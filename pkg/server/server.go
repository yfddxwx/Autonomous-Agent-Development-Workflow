package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"agent/pkg/engine"
	"agent/pkg/git"
	"agent/pkg/llm"
	"agent/pkg/models"
	"agent/pkg/security"
	"agent/pkg/storage"
	"agent/pkg/tools"
)

// Server coordinates HTTP routing, WebSockets, and engine integration.
type Server struct {
	config      *models.Config
	engine      *engine.ReActEngine
	sec         *security.SecurityManager
	registry    *tools.Registry
	gitMgr      *git.Manager
	storage     *storage.Manager
	embedFS     embed.FS
	cancelFunc  context.CancelFunc
	activeWS    *WSConn
	wsMu        sync.Mutex
	cfgMu       sync.RWMutex
	convHistory []models.Message
	activeSess  *models.SessionRecord
}

// NewServer creates a new Server instance.
func NewServer(cfg *models.Config, eng *engine.ReActEngine, sec *security.SecurityManager, reg *tools.Registry, staticFS embed.FS, store *storage.Manager) *Server {
	if store == nil {
		store = storage.NewManager(cfg.WorkspaceDir)
	}

	s := &Server{
		config:      cfg,
		engine:      eng,
		sec:         sec,
		registry:    reg,
		gitMgr:      git.NewManager(sec),
		embedFS:     staticFS,
		storage:     store,
		convHistory: make([]models.Message, 0),
	}

	// Link HITL notification callback to WebSocket broadcast
	sec.GetHITLManager().SetNotifyHandler(func(req *models.ApprovalRequest) {
		s.Broadcast("awaiting_approval", map[string]interface{}{
			"token":   req.Token,
			"tool":    req.Tool,
			"desc":    req.Description,
			"payload": req.Payload,
		})
	})

	return s
}

// Broadcast sends an event to the connected WebSocket client.
func (s *Server) Broadcast(eventType string, data map[string]interface{}) {
	s.wsMu.Lock()
	defer s.wsMu.Unlock()

	if s.activeWS == nil {
		return
	}

	payload := map[string]interface{}{
		"type": eventType,
	}
	for k, v := range data {
		payload[k] = v
	}

	_ = s.activeWS.WriteJSON(payload)
}

// SetupRoutes registers HTTP and WebSocket handlers.
func (s *Server) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/config", s.handleConfig)
	mux.HandleFunc("/api/models", s.handleModels)
	mux.HandleFunc("/api/providers", s.handleProviders)
	mux.HandleFunc("/api/providers/fetch-models", s.handleFetchProviderModels)
	mux.HandleFunc("/api/sessions", s.handleSessions)
	mux.HandleFunc("/api/sessions/detail", s.handleSessionDetail)
	mux.HandleFunc("/api/workspace", s.handleWorkspace)
	mux.HandleFunc("/api/workspace/browse", s.handleBrowseWorkspace)
	mux.HandleFunc("/api/git/status", s.handleGitStatus)
	mux.HandleFunc("/api/git/rollback", s.handleGitRollback)
	mux.HandleFunc("/api/session/snapshots", s.handleSessionSnapshots)
	mux.HandleFunc("/api/session/rollback", s.handleSessionRollback)
	mux.HandleFunc("/api/approval-mode", s.handleApprovalMode)
	mux.HandleFunc("/ws", s.handleWebSocket)

	// Static Web Assets from embed.FS
	mux.Handle("/", http.FileServer(http.FS(s.embedFS)))

	return mux
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.Method == http.MethodGet {
		s.cfgMu.RLock()
		defer s.cfgMu.RUnlock()
		// Always return latest values from security manager
		s.config.WorkspaceDir = s.sec.GetWorkspaceDir()
		s.config.ApprovalMode = s.sec.GetApprovalMode()
		if s.storage != nil {
			appCfg := s.storage.GetConfig()
			s.config.ActiveProfileID = appCfg.ActiveProfileID
			s.config.Profiles = appCfg.Profiles
		}
		_ = json.NewEncoder(w).Encode(s.config)
		return
	}

	if r.Method == http.MethodPost {
		var update models.Config
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.cfgMu.Lock()
		if update.BaseURL != "" {
			s.config.BaseURL = update.BaseURL
		}
		if update.APIKey != "" {
			s.config.APIKey = update.APIKey
		}
		if update.Model != "" {
			s.config.Model = update.Model
		}
		if update.MaxSteps > 0 {
			s.config.MaxSteps = update.MaxSteps
		}
		if update.ApprovalMode != "" {
			s.config.ApprovalMode = update.ApprovalMode
			s.sec.SetApprovalMode(update.ApprovalMode)
			if s.storage != nil {
				s.storage.UpdateApprovalMode(update.ApprovalMode)
			}
		}
		if update.WorkspaceDir != "" {
			if newPath, err := s.sec.SetWorkspaceDir(update.WorkspaceDir); err == nil {
				s.config.WorkspaceDir = newPath
				if s.storage != nil {
					s.storage.UpdateWorkspaceDir(newPath)
				}
			}
		}

		// Update active profile in storage if present
		if s.storage != nil && update.BaseURL != "" && update.Model != "" {
			active := s.storage.GetActiveProfile()
			active.BaseURL = s.config.BaseURL
			active.APIKey = s.config.APIKey
			active.Model = s.config.Model
			active.MaxSteps = s.config.MaxSteps
			_ = s.storage.UpsertProfile(active)
		}

		// Reinitialize engine client
		newClient := llm.NewOpenAIClient(s.config.BaseURL, s.config.APIKey, s.config.Model)
		s.engine.SetClient(newClient)
		s.cfgMu.Unlock()

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if s.storage == nil {
		http.Error(w, "storage not initialized", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		appCfg := s.storage.GetConfig()
		type ProfileWithThinking struct {
			models.ModelProfile
			ThinkingLevels       []string `json:"thinking_levels"`
			DefaultThinkingLevel string   `json:"default_thinking_level"`
		}
		enriched := make([]ProfileWithThinking, len(appCfg.Profiles))
		for i, p := range appCfg.Profiles {
			levels, def := GetModelThinkingCapabilities(p.Model, p.Provider)
			if len(levels) == 0 {
				levels, def = GetModelThinkingCapabilities(p.ID, p.Provider)
			}
			enriched[i] = ProfileWithThinking{
				ModelProfile:         p,
				ThinkingLevels:       levels,
				DefaultThinkingLevel: def,
			}
		}

		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"active_id":            appCfg.ActiveProfileID,
			"profiles":             enriched,
			"thinking_definitions": OfficialThinkingDefinitions,
		})
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			Action    string                `json:"action"`
			ProfileID string                `json:"profile_id"`
			Profile   models.ModelProfile   `json:"profile"`
			Profiles  []models.ModelProfile `json:"profiles"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch req.Action {
		case "switch":
			p, err := s.storage.SetActiveProfile(req.ProfileID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			s.cfgMu.Lock()
			s.config.Provider = p.Provider
			s.config.BaseURL = p.BaseURL
			s.config.APIKey = p.APIKey
			s.config.Model = p.Model
			if p.MaxSteps > 0 {
				s.config.MaxSteps = p.MaxSteps
			}
			newClient := llm.NewOpenAIClient(p.BaseURL, p.APIKey, p.Model)
			s.engine.SetClient(newClient)
			s.cfgMu.Unlock()

			s.Broadcast("model_switched", map[string]interface{}{
				"active_id": p.ID,
				"profile":   p,
			})
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok", "profile": p})
			return

		case "save":
			if err := s.storage.UpsertProfile(req.Profile); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if req.Profile.ID == s.storage.GetConfig().ActiveProfileID {
				s.cfgMu.Lock()
				s.config.Provider = req.Profile.Provider
				s.config.BaseURL = req.Profile.BaseURL
				s.config.APIKey = req.Profile.APIKey
				s.config.Model = req.Profile.Model
				newClient := llm.NewOpenAIClient(req.Profile.BaseURL, req.Profile.APIKey, req.Profile.Model)
				s.engine.SetClient(newClient)
				s.cfgMu.Unlock()
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			return

		case "batch_save":
			if len(req.Profiles) > 0 {
				if err := s.storage.UpsertProfilesBatch(req.Profiles); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			return

		case "delete":
			if err := s.storage.DeleteProfile(req.ProfileID); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			return
		}
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleProviders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_ = json.NewEncoder(w).Encode(OfficialProviders)
}

func (s *Server) handleFetchProviderModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Provider string `json:"provider"`
		BaseURL  string `json:"base_url"`
		APIKey   string `json:"api_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	baseURL := strings.TrimRight(strings.TrimSpace(req.BaseURL), "/")
	if baseURL == "" {
		if meta, ok := GetProviderMeta(req.Provider); ok && meta.BaseURL != "" {
			baseURL = meta.BaseURL
		} else {
			switch strings.ToLower(req.Provider) {
			case "deepseek":
				baseURL = "https://api.deepseek.com"
			default:
				baseURL = "https://api.openai.com/v1"
			}
		}
	}

	type DiscoveredModel struct {
		ID            string `json:"id"`
		Name          string `json:"name"`
		Provider      string `json:"provider"`
		BaseURL       string `json:"base_url"`
		APIKey        string `json:"api_key"`
		Description   string `json:"description,omitempty"`
		Reasoning     bool   `json:"reasoning,omitempty"`
		ContextWindow int    `json:"context_window,omitempty"`
		MaxTokens     int    `json:"max_tokens,omitempty"`
	}

	var discovered []DiscoveredModel

	// Attempt live GET /models discovery
	modelsURL := baseURL + "/models"
	httpReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, modelsURL, nil)
	if err == nil {
		if req.APIKey != "" {
			lowerURL := strings.ToLower(baseURL)
			if strings.Contains(lowerURL, "anthropic") || strings.Contains(lowerURL, "kimi.com") || strings.Contains(lowerURL, "minimax") {
				httpReq.Header.Set("x-api-key", req.APIKey)
				httpReq.Header.Set("anthropic-version", "2023-06-01")
				httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)
			} else if strings.Contains(lowerURL, "azure") {
				httpReq.Header.Set("api-key", req.APIKey)
			} else if strings.Contains(lowerURL, "githubcopilot.com") {
				httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)
				httpReq.Header.Set("Editor-Version", "vscode/1.98.0")
			} else {
				httpReq.Header.Set("Authorization", "Bearer "+req.APIKey)
			}
		}

		client := &http.Client{Timeout: 6 * time.Second}
		resp, err := client.Do(httpReq)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var openAIResp struct {
				Data []struct {
					ID string `json:"id"`
				} `json:"data"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err == nil && len(openAIResp.Data) > 0 {
				for _, m := range openAIResp.Data {
					name := m.ID
					desc := ""
					isReasoning := false
					if strings.Contains(m.ID, "flash") {
						desc = "旗舰级 MoE / 极低延迟与原生多模态"
					} else if strings.Contains(m.ID, "reasoner") || strings.Contains(m.ID, "r1") || strings.Contains(m.ID, "reasoning") {
						desc = "深度长链推理模型"
						isReasoning = true
					} else if strings.Contains(m.ID, "coder") || strings.Contains(m.ID, "code") {
						desc = "专精代码研发模型"
					}
					discovered = append(discovered, DiscoveredModel{
						ID:          m.ID,
						Name:        name,
						Provider:    req.Provider,
						BaseURL:     baseURL,
						APIKey:      req.APIKey,
						Description: desc,
						Reasoning:   isReasoning,
					})
				}
			}
		}
	}

	// Curated latest models fallback from official Pi registry if offline or API key invalid/empty
	if len(discovered) == 0 {
		if meta, ok := GetProviderMeta(req.Provider); ok && len(meta.Models) > 0 {
			for _, m := range meta.Models {
				desc := m.Description
				if desc == "" {
					if m.Reasoning {
						desc = "深度长链推理模型 (Reasoning)"
					} else {
						desc = "官方全能模型"
					}
				}
				discovered = append(discovered, DiscoveredModel{
					ID:            m.ID,
					Name:          m.Name,
					Provider:      req.Provider,
					BaseURL:       baseURL,
					APIKey:        req.APIKey,
					Description:   desc,
					Reasoning:     m.Reasoning,
					ContextWindow: m.ContextWindow,
					MaxTokens:     m.MaxTokens,
				})
			}
		}
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"provider": req.Provider,
		"base_url": baseURL,
		"models":   discovered,
	})
}

func (s *Server) handleSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if s.storage == nil {
		http.Error(w, "storage not initialized", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		list, err := s.storage.ListSessions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		activeID := ""
		if s.activeSess != nil {
			activeID = s.activeSess.ID
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"active_id": activeID,
			"sessions":  list,
		})
		return
	}

	if r.Method == http.MethodPost {
		s.cancelCurrentTask()
		s.convHistory = make([]models.Message, 0)
		s.engine.ClearSnapshots()
		newSess := &models.SessionRecord{
			ID:           fmt.Sprintf("sess_%d", time.Now().UnixNano()),
			Title:        "新研发会话",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			WorkspaceDir: s.sec.GetWorkspaceDir(),
			Messages:     make([]models.Message, 0),
		}
		s.activeSess = newSess
		_ = s.storage.SaveSession(newSess)
		s.Broadcast("session_reset", map[string]interface{}{
			"session_id": newSess.ID,
			"workspace":  newSess.WorkspaceDir,
		})
		_ = json.NewEncoder(w).Encode(newSess)
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleSessionDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing 'id' query param", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		sess, err := s.storage.LoadSession(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		s.cancelCurrentTask()
		s.convHistory = sess.Messages
		s.engine.SetSnapshots(sess.Snapshots)
		s.engine.SetTelemetry(sess.PromptTokens, sess.CompletionTokens)
		s.activeSess = sess

		// Automatically switch workspace directory to this session's folder if recorded
		if sess.WorkspaceDir != "" {
			if newPath, err := s.sec.SetWorkspaceDir(sess.WorkspaceDir); err == nil {
				s.cfgMu.Lock()
				s.config.WorkspaceDir = newPath
				s.cfgMu.Unlock()
				s.Broadcast("workspace_updated", map[string]interface{}{
					"workspace": newPath,
				})
			}
		}

		s.Broadcast("session_loaded", map[string]interface{}{
			"session": sess,
		})

		_ = json.NewEncoder(w).Encode(sess)
		return
	}

	if r.Method == http.MethodDelete {
		_ = s.storage.DeleteSession(id)
		if s.activeSess != nil && s.activeSess.ID == id {
			s.activeSess = nil
			s.convHistory = make([]models.Message, 0)
			s.engine.ClearSnapshots()
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleWorkspace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Workspace string `json:"workspace"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newPath, err := s.sec.SetWorkspaceDir(req.Workspace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.cfgMu.Lock()
	s.config.WorkspaceDir = newPath
	if s.storage != nil {
		s.storage.UpdateWorkspaceDir(newPath)
	}
	if s.activeSess != nil {
		s.activeSess.WorkspaceDir = newPath
		if s.storage != nil {
			_ = s.storage.SaveSession(s.activeSess)
		}
	}
	s.cfgMu.Unlock()

	// Notify all connected clients
	s.Broadcast("workspace_updated", map[string]interface{}{
		"workspace": newPath,
	})

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"workspace": newPath,
	})
}

func (s *Server) handleApprovalMode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ApprovalMode models.ApprovalMode `json:"approval_mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ApprovalMode != models.ApprovalModeManual &&
		req.ApprovalMode != models.ApprovalModeGuarded &&
		req.ApprovalMode != models.ApprovalModeAuto {
		http.Error(w, "invalid approval mode", http.StatusBadRequest)
		return
	}

	s.sec.SetApprovalMode(req.ApprovalMode)

	s.cfgMu.Lock()
	s.config.ApprovalMode = req.ApprovalMode
	if s.storage != nil {
		s.storage.UpdateApprovalMode(req.ApprovalMode)
	}
	s.cfgMu.Unlock()

	s.Broadcast("approval_mode_updated", map[string]interface{}{
		"approval_mode": req.ApprovalMode,
	})

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "ok",
		"approval_mode": req.ApprovalMode,
	})
}

func (s *Server) handleBrowseWorkspace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentWs := s.sec.GetWorkspaceDir()
	escapedWs := strings.ReplaceAll(currentWs, "'", "''")

	// Trigger native Windows FolderBrowserDialog with explicit UTF-8 encoding and current workspace preselected
	script := fmt.Sprintf(`[Console]::OutputEncoding = [System.Text.Encoding]::UTF8; $OutputEncoding = [System.Text.Encoding]::UTF8; Add-Type -AssemblyName System.Windows.Forms; $f = New-Object System.Windows.Forms.FolderBrowserDialog; $f.Description = '选择 Agent 工作空间目录'; $f.ShowNewFolderButton = $true; $f.SelectedPath = '%s'; if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { [Console]::WriteLine($f.SelectedPath) }`, escapedWs)
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script)
	out, err := cmd.Output()
	if err != nil {
		http.Error(w, fmt.Sprintf("无法唤起系统文件夹选择器: %v", err), http.StatusInternalServerError)
		return
	}

	selected := strings.TrimSpace(string(out))
	if selected == "" {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "cancelled"})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"path":   selected,
	})
}

func (s *Server) handleGitStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	ws := s.sec.GetWorkspaceDir()
	isGit := s.gitMgr.IsGitRepo(r.Context(), ws)
	statusText := ""
	if isGit {
		statusText, _ = s.gitMgr.GetStatus(r.Context(), ws)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"is_git": isGit,
		"status": statusText,
	})
}

func (s *Server) handleGitRollback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ws := s.sec.GetWorkspaceDir()
	msg, err := s.gitMgr.Rollback(r.Context(), ws)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.Broadcast("status", map[string]interface{}{"state": "idle", "text": "已回滚 Git 工作区所有改动"})
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": msg})
}

func (s *Server) handleSessionSnapshots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	snaps := s.engine.GetSnapshots()
	_ = json.NewEncoder(w).Encode(snaps)
}

func (s *Server) handleSessionRollback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TurnIndex int `json:"turn_index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.cancelCurrentTask()
	newHistory, err := s.engine.RollbackToTurn(req.TurnIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.convHistory = newHistory

	s.Broadcast("status", map[string]interface{}{"state": "idle", "text": fmt.Sprintf("已回滚至第 %d 轮对话分支", req.TurnIndex+1)})
	s.Broadcast("session_rollback", map[string]interface{}{
		"turn_index": req.TurnIndex,
		"messages":   newHistory,
	})

	s.saveCurrentSession("")

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"messages": newHistory,
	})
}

func (s *Server) saveCurrentSession(firstPrompt string) {
	if s.storage == nil {
		return
	}

	s.cfgMu.Lock()
	defer s.cfgMu.Unlock()

	if s.activeSess == nil {
		s.activeSess = &models.SessionRecord{
			ID:           fmt.Sprintf("sess_%d", time.Now().UnixNano()),
			CreatedAt:    time.Now(),
			WorkspaceDir: s.sec.GetWorkspaceDir(),
		}
	}
	if s.activeSess.WorkspaceDir == "" {
		s.activeSess.WorkspaceDir = s.sec.GetWorkspaceDir()
	}

	if s.activeSess.Title == "" || s.activeSess.Title == "新研发会话" {
		if firstPrompt != "" {
			t := strings.TrimSpace(firstPrompt)
			runes := []rune(t)
			if len(runes) > 20 {
				t = string(runes[:20]) + "..."
			}
			s.activeSess.Title = t
		}
	}
	s.activeSess.Messages = s.convHistory
	s.activeSess.Snapshots = s.engine.GetSnapshots()
	promptTok, complTok, totalTok, costCNY, costUSD := s.engine.GetTelemetry()
	s.activeSess.PromptTokens = promptTok
	s.activeSess.CompletionTokens = complTok
	s.activeSess.TotalTokens = totalTok
	s.activeSess.CostCNY = costCNY
	s.activeSess.CostUSD = costUSD
	_ = s.storage.SaveSession(s.activeSess)
}

type ClientWSMessage struct {
	Type          string              `json:"type"`
	Prompt        string              `json:"prompt,omitempty"`
	ModelProvider string              `json:"model_provider,omitempty"`
	ThinkingLevel string              `json:"thinking_level,omitempty"`
	ApprovalMode  models.ApprovalMode `json:"approval_mode,omitempty"`
	Token         string              `json:"token,omitempty"`
	Approved      bool                `json:"approved,omitempty"`
	Reason        string              `json:"reason,omitempty"`
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := Upgrade(w, r)
	if err != nil {
		log.Printf("WS upgrade error: %v", err)
		return
	}
	defer ws.Close()

	s.wsMu.Lock()
	if s.activeWS != nil {
		_ = s.activeWS.Close()
	}
	s.activeWS = ws
	s.wsMu.Unlock()

	s.Broadcast("status", map[string]interface{}{"state": "idle", "text": "已就绪"})

	for {
		opcode, payload, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if opcode != 0x01 { // text frame
			continue
		}

		var msg ClientWSMessage
		if err := json.Unmarshal(payload, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "user_prompt":
			s.handleUserPrompt(msg.Prompt, msg.ModelProvider, msg.ThinkingLevel)
		case "approval_response":
			s.sec.GetHITLManager().SubmitDecision(msg.Token, msg.Approved, msg.Reason)
		case "ping":
			s.wsMu.Lock()
			if s.activeWS == ws {
				_ = ws.WriteJSON(map[string]string{"type": "pong"})
			}
			s.wsMu.Unlock()
		case "cancel":
			s.cancelCurrentTask()
		case "reset_session":
			s.cancelCurrentTask()
			s.convHistory = make([]models.Message, 0)
			s.engine.ClearSnapshots()
			newSess := &models.SessionRecord{
				ID:           fmt.Sprintf("sess_%d", time.Now().UnixNano()),
				Title:        "新研发会话",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				WorkspaceDir: s.sec.GetWorkspaceDir(),
				Messages:     make([]models.Message, 0),
			}
			s.activeSess = newSess
			if s.storage != nil {
				_ = s.storage.SaveSession(newSess)
			}
			s.Broadcast("status", map[string]interface{}{"state": "idle", "text": "会话已重置"})
			s.Broadcast("session_reset", map[string]interface{}{
				"session_id": newSess.ID,
				"workspace":  newSess.WorkspaceDir,
			})
		}
	}

	s.wsMu.Lock()
	if s.activeWS == ws {
		s.activeWS = nil
	}
	s.wsMu.Unlock()
}

func (s *Server) cancelCurrentTask() {
	if s.cancelFunc != nil {
		s.cancelFunc()
		s.cancelFunc = nil
	}
}

func (s *Server) handleUserPrompt(prompt string, modelProvider string, thinkingLevelOptional ...string) {
	thinkingLevel := ""
	if len(thinkingLevelOptional) > 0 {
		thinkingLevel = thinkingLevelOptional[0]
	}

	s.cancelCurrentTask()

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	// 1. Resolve and dynamically activate the requested model profile from storage
	s.cfgMu.Lock()
	var currentProfile models.ModelProfile
	if s.storage != nil {
		if modelProvider != "" {
			if p, err := s.storage.GetProfile(modelProvider); err == nil && p.ID != "" {
				currentProfile = p
				s.config.Provider = p.Provider
				s.config.BaseURL = p.BaseURL
				s.config.APIKey = p.APIKey
				s.config.Model = p.Model
				if p.MaxSteps > 0 {
					s.config.MaxSteps = p.MaxSteps
				}
				newClient := llm.NewOpenAIClient(p.BaseURL, p.APIKey, p.Model)
				newClient.SetReasoningEffort(thinkingLevel)
				s.engine.SetClient(newClient)
				_, _ = s.storage.SetActiveProfile(p.ID)
			}
		}
		if currentProfile.ID == "" {
			currentProfile = s.storage.GetActiveProfile()
			if currentProfile.ID != "" {
				s.config.Provider = currentProfile.Provider
				s.config.BaseURL = currentProfile.BaseURL
				s.config.APIKey = currentProfile.APIKey
				s.config.Model = currentProfile.Model
				if currentProfile.MaxSteps > 0 {
					s.config.MaxSteps = currentProfile.MaxSteps
				}
				newClient := llm.NewOpenAIClient(currentProfile.BaseURL, currentProfile.APIKey, currentProfile.Model)
				newClient.SetReasoningEffort(thinkingLevel)
				s.engine.SetClient(newClient)
			}
		}
	}
	s.cfgMu.Unlock()

	// 2. Pre-check API Key for remote cloud providers
	isLocal := currentProfile.Provider == "ollama" || currentProfile.Provider == "llama.cpp" || currentProfile.Provider == "lmstudio" ||
		strings.Contains(currentProfile.BaseURL, "localhost") || strings.Contains(currentProfile.BaseURL, "127.0.0.1")
	if strings.TrimSpace(currentProfile.APIKey) == "" && !isLocal && currentProfile.ID != "" {
		s.Broadcast("content", map[string]interface{}{
			"content": fmt.Sprintf("⚠️ **模型尚未配置 API Key**\n\n当前选中的模型 **%s** (`%s`) 尚未配置 API Key，无法发起远程大模型调用。\n\n> 💡 **快速配置**：\n> 1. 点击右上角 **设置图标 (⚙️)**\n> 2. 在左侧官方模型厂商列表中找到该厂商通道\n> 3. 填入该厂商对应的有效 **API Key** 并点击 **保存配置**\n> 4. 保存完成后即可直接在此发起正常对话。", currentProfile.Name, currentProfile.Model),
		})
		s.Broadcast("error", map[string]interface{}{
			"error": fmt.Sprintf("模型 [%s] 缺少 API Key，请先进入右上角设置填写", currentProfile.Name),
		})
		s.Broadcast("status", map[string]interface{}{
			"state": "idle",
			"text":  "未配置 API Key",
		})
		return
	}

	go func() {
		defer cancel()

		onEvent := func(eventType string, data map[string]interface{}) {
			s.Broadcast(eventType, data)
		}

		history, err := s.engine.Run(ctx, prompt, s.convHistory, onEvent)
		if err != nil {
			log.Printf("Engine execution error: %v", err)
			s.Broadcast("content", map[string]interface{}{
				"content": fmt.Sprintf("⚠️ **模型调用失败**: %v\n\n> 💡 **排查建议**：\n> - 请检查右上角设置中该模型的 **Base URL** 与 **API Key** 是否有效配置\n> - 请确认账户或 Token 具有该模型 (`%s`) 的调用权限与可用额度\n> - 若网络环境需要代理，请检查网络连接状态", err, currentProfile.Model),
			})
			s.Broadcast("error", map[string]interface{}{
				"error": fmt.Sprintf("模型生成异常: %v", err),
			})
			s.Broadcast("status", map[string]interface{}{
				"state": "idle",
				"text":  "调用出错",
			})
		} else {
			s.convHistory = history
			s.saveCurrentSession(prompt)
		}
	}()
}

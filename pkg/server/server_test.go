package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"agent/pkg/engine"
	"agent/pkg/llm"
	"agent/pkg/models"
	"agent/pkg/security"
	"agent/pkg/storage"
	"agent/pkg/tools"
	"agent/web"
)

func setupTestServer(t *testing.T) (*Server, http.Handler, string) {
	tempDir := t.TempDir()
	config := &models.Config{
		WorkspaceDir: tempDir,
		ApprovalMode: models.ApprovalModeGuarded,
		Port:         8080,
		Provider:     "deepseek",
		Model:        "deepseek-chat",
		MaxSteps:     10,
	}

	sec := security.NewSecurityManager(tempDir, models.ApprovalModeGuarded)
	reg := tools.NewRegistry()
	tools.RegisterBuiltinTools(reg, sec)
	mockClient := llm.NewOpenAIClient("https://dummy.com", "key", "model")
	reactEngine := engine.NewReActEngine(mockClient, reg, sec, 10)
	store := storage.NewManager(tempDir)

	srv := NewServer(config, reactEngine, sec, reg, web.Assets, store)
	handler := srv.SetupRoutes()
	return srv, handler, tempDir
}

func TestServer_ConfigEndpoints(t *testing.T) {
	_, handler, _ := setupTestServer(t)

	// Test GET /api/config
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var cfg models.Config
	if err := json.NewDecoder(rec.Body).Decode(&cfg); err != nil {
		t.Fatalf("failed to decode config: %v", err)
	}
	if cfg.ApprovalMode != models.ApprovalModeGuarded {
		t.Fatalf("expected guarded approval mode, got %s", cfg.ApprovalMode)
	}

	// Test POST /api/config
	updatePayload := map[string]interface{}{
		"model":         "gpt-4o",
		"approval_mode": "manual",
	}
	body, _ := json.Marshal(updatePayload)
	postReq := httptest.NewRequest(http.MethodPost, "/api/config", bytes.NewReader(body))
	postRec := httptest.NewRecorder()
	handler.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusOK {
		t.Fatalf("expected 200 on config update, got %d", postRec.Code)
	}

	// Verify update took effect
	rec2 := httptest.NewRecorder()
	handler.ServeHTTP(rec2, req)
	var updated models.Config
	_ = json.NewDecoder(rec2.Body).Decode(&updated)
	if updated.Model != "gpt-4o" {
		t.Fatalf("expected updated model gpt-4o, got %s", updated.Model)
	}
	if updated.ApprovalMode != models.ApprovalModeManual {
		t.Fatalf("expected updated approval mode manual, got %s", updated.ApprovalMode)
	}
}

func TestServer_WorkspaceSwitch(t *testing.T) {
	_, handler, baseDir := setupTestServer(t)

	newSubDir := filepath.Join(baseDir, "new_project")
	_ = os.MkdirAll(newSubDir, 0755)

	payload := map[string]string{"workspace": newSubDir}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/workspace", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 on workspace switch, got %d", rec.Code)
	}

	// Test switching to a regular file (not directory)
	filePath := filepath.Join(baseDir, "test_file.txt")
	_ = os.WriteFile(filePath, []byte("content"), 0644)
	badPayload := map[string]string{"workspace": filePath}
	badBody, _ := json.Marshal(badPayload)
	badReq := httptest.NewRequest(http.MethodPost, "/api/workspace", bytes.NewReader(badBody))
	badRec := httptest.NewRecorder()
	handler.ServeHTTP(badRec, badReq)

	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for file workspace, got %d", badRec.Code)
	}
}

func TestServer_ApprovalModeEndpoint(t *testing.T) {
	_, handler, _ := setupTestServer(t)

	modes := []string{"manual", "guarded", "auto"}
	for _, mode := range modes {
		payload := map[string]string{"approval_mode": mode}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/approval-mode", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 setting mode %s, got %d", mode, rec.Code)
		}
	}

	// Invalid mode
	badPayload := map[string]string{"approval_mode": "invalid_mode"}
	badBody, _ := json.Marshal(badPayload)
	badReq := httptest.NewRequest(http.MethodPost, "/api/approval-mode", bytes.NewReader(badBody))
	badRec := httptest.NewRecorder()
	handler.ServeHTTP(badRec, badReq)

	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid mode, got %d", badRec.Code)
	}
}

func TestServer_SessionSnapshots(t *testing.T) {
	_, handler, _ := setupTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/session/snapshots", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from snapshots endpoint, got %d", rec.Code)
	}

	var snaps []models.SessionSnapshot
	if err := json.NewDecoder(rec.Body).Decode(&snaps); err != nil {
		t.Fatalf("failed to decode snapshots: %v", err)
	}
	if len(snaps) != 0 {
		t.Fatalf("expected 0 snapshots initially, got %d", len(snaps))
	}
}

func TestServer_ModelsEndpoints(t *testing.T) {
	_, handler, _ := setupTestServer(t)

	// Test GET /api/models
	req := httptest.NewRequest(http.MethodGet, "/api/models", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from GET /api/models, got %d", rec.Code)
	}

	var resp struct {
		ActiveProfileID string                `json:"active_profile_id"`
		Profiles        []models.ModelProfile `json:"profiles"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode models resp: %v", err)
	}
	if len(resp.Profiles) == 0 {
		t.Fatal("expected default profiles to be present")
	}

	// Test POST /api/models with action=switch
	switchPayload := map[string]string{
		"action":     "switch",
		"profile_id": "gpt-4o",
	}
	body, _ := json.Marshal(switchPayload)
	postReq := httptest.NewRequest(http.MethodPost, "/api/models", bytes.NewReader(body))
	postRec := httptest.NewRecorder()
	handler.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from POST /api/models switch, got %d", postRec.Code)
	}
}

func TestServer_SessionsEndpoints(t *testing.T) {
	_, handler, _ := setupTestServer(t)

	// Test GET /api/sessions
	req := httptest.NewRequest(http.MethodGet, "/api/sessions", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from GET /api/sessions, got %d", rec.Code)
	}

	// Test POST /api/sessions (create new session)
	newPayload := map[string]string{"title": "Test Session"}
	body, _ := json.Marshal(newPayload)
	postReq := httptest.NewRequest(http.MethodPost, "/api/sessions", bytes.NewReader(body))
	postRec := httptest.NewRecorder()
	handler.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from POST /api/sessions, got %d", postRec.Code)
	}

	var sess models.SessionRecord
	if err := json.NewDecoder(postRec.Body).Decode(&sess); err != nil {
		t.Fatalf("failed to decode new session: %v", err)
	}
	if sess.ID == "" {
		t.Fatal("expected non-empty session ID")
	}

	// Test GET /api/sessions/detail?id=...
	detailReq := httptest.NewRequest(http.MethodGet, "/api/sessions/detail?id="+sess.ID, nil)
	detailRec := httptest.NewRecorder()
	handler.ServeHTTP(detailRec, detailReq)

	if detailRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from GET /api/sessions/detail, got %d", detailRec.Code)
	}

	// Test DELETE /api/sessions/detail?id=...
	delReq := httptest.NewRequest(http.MethodDelete, "/api/sessions/detail?id="+sess.ID, nil)
	delRec := httptest.NewRecorder()
	handler.ServeHTTP(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from DELETE /api/sessions/detail, got %d", delRec.Code)
	}
}

func TestServer_SessionWorkspaceSwitchOnLoad(t *testing.T) {
	srv, handler, baseDir := setupTestServer(t)

	// Create a secondary workspace directory
	projectB := filepath.Join(baseDir, "project_b")
	if err := os.MkdirAll(projectB, 0755); err != nil {
		t.Fatalf("failed to create projectB dir: %v", err)
	}

	// Persist a session tied to projectB
	sess := &models.SessionRecord{
		ID:               "sess_project_b",
		Title:            "Project B Session",
		WorkspaceDir:     projectB,
		PromptTokens:     1000,
		CompletionTokens: 500,
		TotalTokens:      1500,
	}
	if err := srv.storage.SaveSession(sess); err != nil {
		t.Fatalf("failed to save session: %v", err)
	}

	// Verify current workspace before load is baseDir
	if srv.sec.GetWorkspaceDir() != baseDir {
		t.Fatalf("expected baseDir %s, got %s", baseDir, srv.sec.GetWorkspaceDir())
	}

	// Load session via GET /api/sessions/detail?id=sess_project_b
	req := httptest.NewRequest(http.MethodGet, "/api/sessions/detail?id=sess_project_b", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from GET /api/sessions/detail, got %d", rec.Code)
	}

	// Verify workspace automatically switched to projectB
	if srv.sec.GetWorkspaceDir() != projectB {
		t.Fatalf("expected security manager workspace to switch to %s, got %s", projectB, srv.sec.GetWorkspaceDir())
	}
	if srv.config.WorkspaceDir != projectB {
		t.Fatalf("expected server config workspace to switch to %s, got %s", projectB, srv.config.WorkspaceDir)
	}

	// Verify telemetry was restored and cost is non-zero
	pTok, cTok, total, costCNY, costUSD := srv.engine.GetTelemetry()
	if pTok != 1000 || cTok != 500 || total != 1500 {
		t.Fatalf("expected tokens 1000/500/1500, got %d/%d/%d", pTok, cTok, total)
	}
	if costCNY <= 0 {
		t.Fatalf("expected positive costCNY, got %f", costCNY)
	}
	if costUSD <= 0 {
		t.Fatalf("expected positive costUSD, got %f", costUSD)
	}
}

func TestServer_FetchProviderModelsAndBatchSave(t *testing.T) {
	_, handler, _ := setupTestServer(t)

	// 1. Test POST /api/providers/fetch-models for DeepSeek
	discPayload := map[string]string{
		"provider": "deepseek",
	}
	body, _ := json.Marshal(discPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/providers/fetch-models", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /api/providers/fetch-models, got %d", rec.Code)
	}

	var discResp struct {
		Provider string `json:"provider"`
		Models   []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&discResp); err != nil {
		t.Fatalf("failed to decode discovered models: %v", err)
	}
	if len(discResp.Models) == 0 {
		t.Fatal("expected discovered models for deepseek")
	}

	// Verify deepseek-flash is present in discovered models
	foundFlash := false
	for _, m := range discResp.Models {
		if m.ID == "deepseek-flash" {
			foundFlash = true
			break
		}
	}
	if !foundFlash {
		t.Fatal("expected deepseek-flash in discovered deepseek models")
	}

	// 2. Test POST /api/models with action=batch_save
	batchPayload := map[string]interface{}{
		"action": "batch_save",
		"profiles": []models.ModelProfile{
			{
				ID:       "deepseek-flash",
				Name:     "DeepSeek V4.1 Flash",
				Provider: "deepseek",
				BaseURL:  "https://api.deepseek.com",
				Model:    "deepseek-flash",
				MaxSteps: 15,
			},
		},
	}
	batchBody, _ := json.Marshal(batchPayload)
	batchReq := httptest.NewRequest(http.MethodPost, "/api/models", bytes.NewReader(batchBody))
	batchRec := httptest.NewRecorder()
	handler.ServeHTTP(batchRec, batchReq)

	if batchRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from batch_save, got %d", batchRec.Code)
	}
}

func TestServer_AllOfficialPiProviders(t *testing.T) {
	_, handler, _ := setupTestServer(t)

	// 1. Test GET /api/providers
	req := httptest.NewRequest(http.MethodGet, "/api/providers", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from GET /api/providers, got %d", rec.Code)
	}

	var provs []ProviderMeta
	if err := json.NewDecoder(rec.Body).Decode(&provs); err != nil {
		t.Fatalf("failed to decode providers: %v", err)
	}

	if len(provs) < 40 {
		t.Fatalf("expected at least 40 providers from Pi registry, got %d", len(provs))
	}

	provMap := make(map[string]ProviderMeta)
	for _, p := range provs {
		provMap[p.ID] = p
	}

	// 2. Verify Key Facts specified by user & Pi ecosystem:
	tests := []struct {
		id           string
		expectedBase string
		mustHaveMod  string
	}{
		{"minimax-cn", "https://api.minimaxi.com/anthropic", "MiniMax-M3"},
		{"minimax", "https://api.minimax.io/anthropic", "MiniMax-M3"},
		{"kimi-coding", "https://api.kimi.com/coding", "kimi-for-coding"},
		{"xiaomi", "https://api.xiaomimimo.com/v1", "mimo-v2.6-flash"},
		{"xiaomi-token-plan-cn", "https://token-plan-cn.xiaomimimo.com/v1", "mimo-v2.6-flash"},
		{"xiaomi-token-plan-ams", "https://token-plan-ams.xiaomimimo.com/v1", "mimo-v2.6-flash"},
		{"xiaomi-token-plan-sgp", "https://token-plan-sgp.xiaomimimo.com/v1", "mimo-v2.6-flash"},
		{"zai", "https://api.z.ai/api/coding/paas/v4", "glm-5.3"},
		{"zai-coding-cn", "https://open.bigmodel.cn/api/coding/paas/v4", "glm-5.3"},
		{"qwen-token-plan-cn", "https://token-plan.cn-beijing.maas.aliyuncs.com/compatible-mode/v1", "deepseek-v4-flash"},
		{"ant-ling", "https://api.ant-ling.com/v1", "Ling-2.6-1T"},
		{"deepseek", "https://api.deepseek.com", "deepseek-flash"},
		{"github-copilot", "https://api.individual.githubcopilot.com", "claude-fable-5"},
	}

	for _, tt := range tests {
		p, exists := provMap[tt.id]
		if !exists {
			t.Fatalf("missing required Pi provider: %s", tt.id)
		}
		if p.BaseURL != tt.expectedBase {
			t.Errorf("provider %s BaseURL mismatch: expected %s, got %s", tt.id, tt.expectedBase, p.BaseURL)
		}

		// Also test POST /api/providers/fetch-models returns correct fallback models
		body, _ := json.Marshal(map[string]string{"provider": tt.id})
		fReq := httptest.NewRequest(http.MethodPost, "/api/providers/fetch-models", bytes.NewReader(body))
		fRec := httptest.NewRecorder()
		handler.ServeHTTP(fRec, fReq)

		if fRec.Code != http.StatusOK {
			t.Fatalf("fetch-models for %s returned code %d", tt.id, fRec.Code)
		}

		var fResp struct {
			Provider string `json:"provider"`
			BaseURL  string `json:"base_url"`
			Models   []struct {
				ID string `json:"id"`
			} `json:"models"`
		}
		_ = json.NewDecoder(fRec.Body).Decode(&fResp)

		if fResp.BaseURL != tt.expectedBase {
			t.Errorf("fetch-models for %s returned BaseURL %s, expected %s", tt.id, fResp.BaseURL, tt.expectedBase)
		}

		found := false
		for _, m := range fResp.Models {
			if m.ID == tt.mustHaveMod {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("provider %s missing expected flagship model %s in fallback discovery", tt.id, tt.mustHaveMod)
		}
	}
}

func TestServer_HandleUserPrompt_DynamicProfileAndMissingKey(t *testing.T) {
	srv, _, _ := setupTestServer(t)

	// 1. Add a remote custom model without API Key
	unconfiguredProfile := models.ModelProfile{
		ID:       "mimo-v2.6-flash-nokey",
		Name:     "mimo-v2.6-flash-nokey",
		Provider: "xiaomi",
		BaseURL:  "https://api.xiaomimimo.com/v1",
		APIKey:   "",
		Model:    "mimo-v2.6-flash",
		MaxSteps: 15,
	}
	if err := srv.storage.UpsertProfile(unconfiguredProfile); err != nil {
		t.Fatalf("failed to upsert profile: %v", err)
	}

	// Track broadcasts
	var broadcastEvents []string
	srv.activeWS = nil // no panic even without ws

	// Test invoking prompt with unconfigured profile
	srv.handleUserPrompt("你好", "mimo-v2.6-flash-nokey")

	// Config should have dynamically switched to the profile
	srv.cfgMu.Lock()
	if srv.config.Model != "mimo-v2.6-flash" {
		t.Fatalf("expected config model mimo-v2.6-flash, got %s", srv.config.Model)
	}
	if srv.config.Provider != "xiaomi" {
		t.Fatalf("expected config provider xiaomi, got %s", srv.config.Provider)
	}
	if srv.config.BaseURL != "https://api.xiaomimimo.com/v1" {
		t.Fatalf("expected config BaseURL https://api.xiaomimimo.com/v1, got %s", srv.config.BaseURL)
	}
	srv.cfgMu.Unlock()

	// 2. Add configured profile with valid settings
	configuredProfile := models.ModelProfile{
		ID:       "mimo-v2.6-flash-configured",
		Name:     "mimo-v2.6-flash-configured",
		Provider: "xiaomi",
		BaseURL:  "https://api.xiaomimimo.com/v1",
		APIKey:   "sk-test-valid-key",
		Model:    "mimo-v2.6-flash",
		MaxSteps: 12,
	}
	if err := srv.storage.UpsertProfile(configuredProfile); err != nil {
		t.Fatalf("failed to upsert profile: %v", err)
	}

	srv.handleUserPrompt("测试", "mimo-v2.6-flash-configured")

	srv.cfgMu.Lock()
	if srv.config.APIKey != "sk-test-valid-key" {
		t.Fatalf("expected configured API key sk-test-valid-key, got %s", srv.config.APIKey)
	}
	if srv.config.MaxSteps != 12 {
		t.Fatalf("expected max steps 12, got %d", srv.config.MaxSteps)
	}
	srv.cfgMu.Unlock()

	_ = broadcastEvents
}

func TestServer_ThinkingDepthAndModelCapabilities(t *testing.T) {
	// 1. Verify Pi-aligned model thinking capabilities
	dsLevels, dsDef := GetModelThinkingCapabilities("deepseek-flash", "deepseek")
	if len(dsLevels) != 3 || dsDef != "high" {
		t.Fatalf("expected deepseek-flash to have 3 levels and default high, got %v, %s", dsLevels, dsDef)
	}

	mimoLevels, mimoDef := GetModelThinkingCapabilities("mimo-v2.6-flash", "xiaomi")
	if len(mimoLevels) != 2 || mimoDef != "high" {
		t.Fatalf("expected mimo-v2.6-flash to have 2 levels [off, high], got %v, %s", mimoLevels, mimoDef)
	}

	geminiLevels, geminiDef := GetModelThinkingCapabilities("gemini-3.5-flash", "google")
	if len(geminiLevels) != 4 || geminiDef != "high" {
		t.Fatalf("expected gemini-3.5-flash to have 4 levels, got %v, %s", geminiLevels, geminiDef)
	}

	gpt4oLevels, _ := GetModelThinkingCapabilities("gpt-4o", "openai")
	if len(gpt4oLevels) != 0 {
		t.Fatalf("expected gpt-4o to have 0 thinking levels, got %v", gpt4oLevels)
	}

	// 2. Verify GET /api/models returns enriched profiles with thinking_definitions
	srv, _, _ := setupTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/models", nil)
	w := httptest.NewRecorder()
	srv.handleModels(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp struct {
		ActiveID            string                      `json:"active_id"`
		Profiles            []map[string]interface{}    `json:"profiles"`
		ThinkingDefinitions map[string]ThinkingLevelDef `json:"thinking_definitions"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode models response: %v", err)
	}

	if len(resp.ThinkingDefinitions) < 7 {
		t.Fatalf("expected 7 thinking definitions, got %d", len(resp.ThinkingDefinitions))
	}
	if _, ok := resp.ThinkingDefinitions["high"]; !ok {
		t.Fatalf("missing 'high' thinking definition")
	}

	// 3. Verify handleUserPrompt passes thinking_level
	srv.handleUserPrompt("测试思考", "deepseek", "low")
	client, ok := srv.engine.GetClient().(*llm.OpenAIClient)
	if !ok {
		t.Fatalf("expected engine client to be OpenAIClient")
	}
	if client.ReasoningEffort != "low" {
		t.Fatalf("expected reasoning effort 'low', got '%s'", client.ReasoningEffort)
	}
}



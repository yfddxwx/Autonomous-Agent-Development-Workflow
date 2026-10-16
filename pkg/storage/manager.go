package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"agent/pkg/models"
)

// Manager handles on-disk persistence for configurations and chat sessions.
type Manager struct {
	mu           sync.RWMutex
	configPath   string
	sessionsDir  string
	config       *models.AppConfig
	activeSessID string
}

// NewManager creates a StorageManager with configured paths.
func NewManager(baseDir string) *Manager {
	agentDir := filepath.Join(baseDir, ".agent")
	_ = os.MkdirAll(agentDir, 0755)

	sessionsDir := filepath.Join(agentDir, "sessions")
	_ = os.MkdirAll(sessionsDir, 0755)

	cfgPath := filepath.Join(baseDir, "config.json")

	mgr := &Manager{
		configPath:  cfgPath,
		sessionsDir: sessionsDir,
	}

	mgr.loadOrCreateConfig(baseDir)
	return mgr
}

func defaultProfiles() []models.ModelProfile {
	return []models.ModelProfile{
		{
			ID:       "deepseek",
			Name:     "DeepSeek V4.1 Flash",
			Provider: "deepseek",
			BaseURL:  "https://api.deepseek.com",
			APIKey:   "",
			Model:    "deepseek-flash",
			MaxSteps: 15,
		},
		{
			ID:       "deepseek-v3",
			Name:     "DeepSeek V3",
			Provider: "deepseek",
			BaseURL:  "https://api.deepseek.com",
			APIKey:   "",
			Model:    "deepseek-chat",
			MaxSteps: 15,
		},
		{
			ID:       "gpt-4o",
			Name:     "OpenAI GPT-4o",
			Provider: "gpt-4o",
			BaseURL:  "https://api.openai.com/v1",
			APIKey:   "",
			Model:    "gpt-4o",
			MaxSteps: 15,
		},
		{
			ID:       "claude-3-7-sonnet",
			Name:     "Claude 3.7 Sonnet",
			Provider: "anthropic",
			BaseURL:  "https://api.anthropic.com/v1",
			APIKey:   "",
			Model:    "claude-3-7-sonnet-latest",
			MaxSteps: 15,
		},
		{
			ID:       "gemini",
			Name:     "Google Gemini 2.0 Flash",
			Provider: "gemini",
			BaseURL:  "https://generativelanguage.googleapis.com/v1beta/openai",
			APIKey:   "",
			Model:    "gemini-2.0-flash",
			MaxSteps: 15,
		},
		{
			ID:       "ollama",
			Name:     "本地 Ollama (Llama 3.3)",
			Provider: "ollama",
			BaseURL:  "http://localhost:11434/v1",
			APIKey:   "",
			Model:    "llama3.3",
			MaxSteps: 15,
		},
	}
}

func (m *Manager) loadOrCreateConfig(workspaceFallback string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.configPath)
	if err == nil {
		var cfg models.AppConfig
		if err := json.Unmarshal(data, &cfg); err == nil && len(cfg.Profiles) > 0 {
			m.config = &cfg
			if m.config.ActiveProfileID == "" {
				m.config.ActiveProfileID = cfg.Profiles[0].ID
			}
			return
		}
	}

	// Create default config
	m.config = &models.AppConfig{
		ActiveProfileID: "deepseek",
		Profiles:        defaultProfiles(),
		WorkspaceDir:    workspaceFallback,
		ApprovalMode:    models.ApprovalModeGuarded,
	}

	// Check environment variables for initial API keys
	for i := range m.config.Profiles {
		p := &m.config.Profiles[i]
		if p.ID == "deepseek" {
			if k := os.Getenv("DEEPSEEK_API_KEY"); k != "" {
				p.APIKey = k
			}
		}
		if p.ID == "gpt-4o" {
			if k := os.Getenv("OPENAI_API_KEY"); k != "" {
				p.APIKey = k
			}
		}
		if p.ID == "gemini" {
			if k := os.Getenv("GEMINI_API_KEY"); k != "" {
				p.APIKey = k
			}
		}
		if k := os.Getenv("AGENT_API_KEY"); k != "" && p.APIKey == "" {
			p.APIKey = k
		}
	}

	_ = m.saveConfigLocked()
}

func (m *Manager) saveConfigLocked() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.configPath, data, 0644)
}

// GetConfig returns the full AppConfig.
func (m *Manager) GetConfig() models.AppConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return *m.config
}

// GetActiveProfile returns the currently selected ModelProfile.
func (m *Manager) GetActiveProfile() models.ModelProfile {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.config.Profiles {
		if p.ID == m.config.ActiveProfileID {
			return p
		}
	}
	if len(m.config.Profiles) > 0 {
		return m.config.Profiles[0]
	}
	return defaultProfiles()[0]
}

// GetProfile retrieves a profile by ID or model name.
func (m *Manager) GetProfile(profileID string) (models.ModelProfile, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, p := range m.config.Profiles {
		if p.ID == profileID || p.Model == profileID {
			return p, nil
		}
	}
	return models.ModelProfile{}, fmt.Errorf("profile with ID '%s' not found", profileID)
}

// SetActiveProfile switches the active model profile ID and saves.
func (m *Manager) SetActiveProfile(profileID string) (models.ModelProfile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.config.Profiles {
		if p.ID == profileID {
			m.config.ActiveProfileID = profileID
			_ = m.saveConfigLocked()
			return p, nil
		}
	}
	return models.ModelProfile{}, fmt.Errorf("profile with ID '%s' not found", profileID)
}

// UpsertProfile adds or updates a model profile and saves.
func (m *Manager) UpsertProfile(profile models.ModelProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if profile.ID == "" {
		profile.ID = fmt.Sprintf("model_%d", time.Now().Unix())
	}
	if profile.Name == "" {
		profile.Name = profile.Model
	}

	found := false
	for i, p := range m.config.Profiles {
		if p.ID == profile.ID {
			m.config.Profiles[i] = profile
			found = true
			break
		}
	}
	if !found {
		m.config.Profiles = append(m.config.Profiles, profile)
	}

	return m.saveConfigLocked()
}

// UpsertProfilesBatch adds or updates multiple model profiles at once.
func (m *Manager) UpsertProfilesBatch(profiles []models.ModelProfile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, profile := range profiles {
		if profile.ID == "" {
			profile.ID = profile.Model
		}
		if profile.Name == "" {
			profile.Name = profile.Model
		}
		if profile.MaxSteps <= 0 {
			profile.MaxSteps = 15
		}

		found := false
		for i, p := range m.config.Profiles {
			if p.ID == profile.ID || p.Model == profile.Model {
				m.config.Profiles[i] = profile
				found = true
				break
			}
		}
		if !found {
			m.config.Profiles = append(m.config.Profiles, profile)
		}
	}

	return m.saveConfigLocked()
}

// DeleteProfile removes a model profile.
func (m *Manager) DeleteProfile(profileID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.config.Profiles) <= 1 {
		return fmt.Errorf("至少需要保留一个模型配置")
	}

	newProfiles := make([]models.ModelProfile, 0, len(m.config.Profiles)-1)
	for _, p := range m.config.Profiles {
		if p.ID != profileID {
			newProfiles = append(newProfiles, p)
		}
	}
	m.config.Profiles = newProfiles

	if m.config.ActiveProfileID == profileID {
		m.config.ActiveProfileID = m.config.Profiles[0].ID
	}

	return m.saveConfigLocked()
}

// UpdateWorkspaceDir updates workspace directory in config.
func (m *Manager) UpdateWorkspaceDir(dir string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.WorkspaceDir = dir
	_ = m.saveConfigLocked()
}

// UpdateApprovalMode updates approval mode in config.
func (m *Manager) UpdateApprovalMode(mode models.ApprovalMode) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.ApprovalMode = mode
	_ = m.saveConfigLocked()
}

// --- Session Persistence ---

// SessionHeader contains summary information about a saved session.
type SessionHeader struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	WorkspaceDir string    `json:"workspace_dir"`
	MessageCount int       `json:"message_count"`
	TotalTokens  int       `json:"total_tokens"`
	CostUSD      float64   `json:"cost_usd"`
	CostCNY      float64   `json:"cost_cny"`
}

// ListSessions returns all saved sessions sorted by updated_at desc.
func (m *Manager) ListSessions() ([]SessionHeader, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entries, err := os.ReadDir(m.sessionsDir)
	if err != nil {
		return nil, err
	}

	headers := make([]SessionHeader, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		path := filepath.Join(m.sessionsDir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var sess models.SessionRecord
		if err := json.Unmarshal(data, &sess); err == nil && sess.ID != "" {
			headers = append(headers, SessionHeader{
				ID:           sess.ID,
				Title:        sess.Title,
				CreatedAt:    sess.CreatedAt,
				UpdatedAt:    sess.UpdatedAt,
				WorkspaceDir: sess.WorkspaceDir,
				MessageCount: len(sess.Messages),
				TotalTokens:  sess.TotalTokens,
				CostUSD:      sess.CostUSD,
				CostCNY:      sess.CostCNY,
			})
		}
	}

	sort.Slice(headers, func(i, j int) bool {
		return headers[i].UpdatedAt.After(headers[j].UpdatedAt)
	})

	return headers, nil
}

// SaveSession persists a session to disk.
func (m *Manager) SaveSession(sess *models.SessionRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if sess.ID == "" {
		sess.ID = fmt.Sprintf("sess_%d", time.Now().UnixNano())
	}
	sess.UpdatedAt = time.Now()
	if sess.CreatedAt.IsZero() {
		sess.CreatedAt = sess.UpdatedAt
	}
	if sess.Title == "" {
		sess.Title = "新研发会话"
	}

	m.activeSessID = sess.ID

	data, err := json.MarshalIndent(sess, "", "  ")
	if err != nil {
		return err
	}

	filePath := filepath.Join(m.sessionsDir, fmt.Sprintf("%s.json", sess.ID))
	return os.WriteFile(filePath, data, 0644)
}

// LoadSession retrieves a session record by ID.
func (m *Manager) LoadSession(id string) (*models.SessionRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	filePath := filepath.Join(m.sessionsDir, fmt.Sprintf("%s.json", id))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("session '%s' not found: %w", id, err)
	}

	var sess models.SessionRecord
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("failed to parse session: %w", err)
	}

	m.activeSessID = id
	return &sess, nil
}

// DeleteSession removes a session by ID.
func (m *Manager) DeleteSession(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	filePath := filepath.Join(m.sessionsDir, fmt.Sprintf("%s.json", id))
	_ = os.Remove(filePath)
	if m.activeSessID == id {
		m.activeSessID = ""
	}
	return nil
}

// GetActiveSessionID returns the currently selected session ID.
func (m *Manager) GetActiveSessionID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activeSessID
}

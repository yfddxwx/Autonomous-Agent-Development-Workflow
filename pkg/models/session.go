package models

import "time"

// SessionSnapshot represents a point-in-time state of the conversation tree.
type SessionSnapshot struct {
	ID        string    `json:"id"`
	ParentID  string    `json:"parent_id,omitempty"`
	Branch    string    `json:"branch,omitempty"`
	TurnIndex int       `json:"turn_index"`
	Prompt    string    `json:"prompt"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
	Messages  []Message `json:"messages,omitempty"`
}

// SessionRecord represents a persisted conversation session.
type SessionRecord struct {
	ID               string            `json:"id"`
	Title            string            `json:"title"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	WorkspaceDir     string            `json:"workspace_dir"`
	Messages         []Message         `json:"messages"`
	Snapshots        []SessionSnapshot `json:"snapshots"`
	PromptTokens     int               `json:"prompt_tokens"`
	CompletionTokens int               `json:"completion_tokens"`
	TotalTokens      int               `json:"total_tokens"`
	CostUSD          float64           `json:"cost_usd"`
	CostCNY          float64           `json:"cost_cny"`
}

// TokenUsage holds raw token counts from the LLM provider.
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// TokenTelemetry tracks token usage and estimated cost.
type TokenTelemetry struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	CostUSD          float64 `json:"cost_usd"`
	CostCNY          float64 `json:"cost_cny"`
	EstimatedCostUSD float64 `json:"estimated_cost_usd"`
}

// SkillDefinition represents a locally discovered skill instruction.
type SkillDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Path        string `json:"path"`
}

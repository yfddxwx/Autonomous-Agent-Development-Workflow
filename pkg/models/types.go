package models

import "time"

// MessageRole represents the role of a participant in a conversation.
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleTool      MessageRole = "tool"
)

// FunctionCall represents a model-requested function invocation.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ToolCall represents a specific tool call with an ID.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// Message represents a chat message.
type Message struct {
	Role       MessageRole `json:"role"`
	Content    string      `json:"content"`
	Thought    string      `json:"thought,omitempty"`
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
	Usage      *TokenUsage `json:"usage,omitempty"`
}

// ToolDefinition represents JSON schema for function calling.
type ToolDefinition struct {
	Type     string                 `json:"type"`
	Function ToolFunctionDefinition `json:"function"`
}

// ToolFunctionDefinition represents function name, description, and schema.
type ToolFunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// StreamChunk represents an incremental piece of output from the model.
type StreamChunk struct {
	Thought   string
	Content   string
	ToolCalls []ToolCall
	Done      bool
}

// ApprovalRequest represents a pending action awaiting user confirmation.
type ApprovalRequest struct {
	Token       string      `json:"token"`
	Tool        string      `json:"tool"`
	Description string      `json:"desc"`
	Payload     interface{} `json:"payload"`
	CreatedAt   time.Time   `json:"created_at"`
}

// ApprovalResponse represents the user's decision.
type ApprovalResponse struct {
	Token    string `json:"token"`
	Approved bool   `json:"approved"`
	Reason   string `json:"reason,omitempty"`
}

// ApprovalMode defines the security approval policy level.
type ApprovalMode string

const (
	ApprovalModeManual  ApprovalMode = "manual"  // 人工审核: 所有写文件与终端命令均需人工放行
	ApprovalModeGuarded ApprovalMode = "guarded" // 有限审核: 安全探测与测试放行，高危命令与关键改写需审批
	ApprovalModeAuto    ApprovalMode = "auto"    // 完全不审核: 全自动执行，适用于高信任场景
)

// ModelProfile represents a configured LLM provider and parameters.
type ModelProfile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
	MaxSteps int    `json:"max_steps"`
}

// AppConfig represents persisted configuration on disk.
type AppConfig struct {
	ActiveProfileID string         `json:"active_profile_id"`
	Profiles        []ModelProfile `json:"profiles"`
	WorkspaceDir    string         `json:"workspace_dir"`
	ApprovalMode    ApprovalMode   `json:"approval_mode"`
}

// TeamRole represents autonomous virtual engineering roles.
type TeamRole string

const (
	TeamRoleSolo      TeamRole = "solo"      // 极简直驱
	TeamRoleArchitect TeamRole = "architect" // 架构规划师
	TeamRoleEngineer  TeamRole = "engineer"  // 全栈工程师
	TeamRoleQA        TeamRole = "qa"        // 测试与安全评审
)

// Config represents runtime configuration.
type Config struct {
	WorkspaceDir    string         `json:"workspace"`
	Port            int            `json:"port"`
	Provider        string         `json:"provider"`
	BaseURL         string         `json:"base_url"`
	APIKey          string         `json:"api_key"`
	Model           string         `json:"model"`
	MaxSteps        int            `json:"max_steps"`
	RequireApproval bool           `json:"require_approval"`
	ApprovalMode    ApprovalMode   `json:"approval_mode"`
	ActiveProfileID string         `json:"active_profile_id,omitempty"`
	Profiles        []ModelProfile `json:"profiles,omitempty"`
}

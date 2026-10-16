package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"agent/pkg/models"
)

// SecurityManager handles path validation, command safety checks, and approval policies.
type SecurityManager struct {
	mu           sync.RWMutex
	WorkspaceDir string
	Mode         models.ApprovalMode
	hitl         *HITLManager
}

// NewSecurityManager creates a new instance.
func NewSecurityManager(workspaceDir string, mode models.ApprovalMode) *SecurityManager {
	if mode == "" {
		mode = models.ApprovalModeGuarded
	}
	absWorkspace, err := filepath.Abs(workspaceDir)
	if err != nil {
		absWorkspace = workspaceDir
	}
	return &SecurityManager{
		WorkspaceDir: absWorkspace,
		Mode:         mode,
		hitl:         NewHITLManager(),
	}
}

// SetWorkspaceDir dynamically updates and validates the active workspace directory.
func (sm *SecurityManager) SetWorkspaceDir(newDir string) (string, error) {
	if strings.TrimSpace(newDir) == "" {
		return "", errors.New("工作空间路径不能为空")
	}

	cleanDir := strings.Trim(strings.TrimSpace(newDir), "\"'")
	cleanDir = filepath.FromSlash(cleanDir)
	absWorkspace, err := filepath.Abs(cleanDir)
	if err != nil {
		return "", fmt.Errorf("无效的路径: %w", err)
	}

	info, err := os.Stat(absWorkspace)
	if err != nil {
		if os.IsNotExist(err) {
			// Try creating directory if not existing
			if mkErr := os.MkdirAll(absWorkspace, 0755); mkErr != nil {
				return "", fmt.Errorf("工作空间目录不存在且无法自动创建: %w", mkErr)
			}
		} else {
			return "", fmt.Errorf("无法访问路径: %w", err)
		}
	} else if !info.IsDir() {
		return "", fmt.Errorf("指定的路径 '%s' 不是目录", absWorkspace)
	}

	sm.mu.Lock()
	sm.WorkspaceDir = absWorkspace
	sm.mu.Unlock()

	return absWorkspace, nil
}

// GetWorkspaceDir safely returns current workspace.
func (sm *SecurityManager) GetWorkspaceDir() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.WorkspaceDir
}

// SetApprovalMode updates the active approval policy.
func (sm *SecurityManager) SetApprovalMode(mode models.ApprovalMode) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.Mode = mode
}

// GetApprovalMode returns the current approval policy.
func (sm *SecurityManager) GetApprovalMode() models.ApprovalMode {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.Mode
}

// GetHITLManager returns the underlying HITL manager.
func (sm *SecurityManager) GetHITLManager() *HITLManager {
	return sm.hitl
}

// ValidatePath ensures the given path does not escape the workspace directory.
func (sm *SecurityManager) ValidatePath(targetPath string) (string, error) {
	if targetPath == "" {
		return "", errors.New("empty path provided")
	}

	sm.mu.RLock()
	currentWorkspace := sm.WorkspaceDir
	sm.mu.RUnlock()

	cleanTarget := strings.TrimSpace(targetPath)
	cleanTarget = strings.Trim(cleanTarget, `"'`)

	// On Windows, if path is not absolute (does not start with drive letter C: or D:),
	// strip leading '/' and '\' so e.g. "/index.html" is treated as relative to workspace,
	// avoiding Go's filepath.Join treating it as volume root.
	if !filepath.IsAbs(cleanTarget) {
		cleanTarget = strings.TrimLeft(cleanTarget, "/\\")
		if cleanTarget == "" {
			cleanTarget = "."
		}
	}

	var fullPath string
	if filepath.IsAbs(cleanTarget) {
		fullPath = filepath.Clean(cleanTarget)
	} else {
		fullPath = filepath.Clean(filepath.Join(currentWorkspace, cleanTarget))
	}

	// Case-insensitive prefix check for Windows compatibility
	workspaceLower := filepath.Clean(strings.ToLower(currentWorkspace))
	fullPathLower := filepath.Clean(strings.ToLower(fullPath))

	sep := string(filepath.Separator)
	isInside := fullPathLower == workspaceLower ||
		strings.HasPrefix(fullPathLower, workspaceLower+sep) ||
		(strings.HasSuffix(workspaceLower, sep) && strings.HasPrefix(fullPathLower, workspaceLower))

	if !isInside {
		return "", fmt.Errorf("security violation: path '%s' escapes workspace boundary '%s'", targetPath, currentWorkspace)
	}

	return fullPath, nil
}

// DangerousCommandRules contains blacklisted command patterns.
var dangerousPatterns = []string{
	"format ",
	"diskpart",
	"mkfs",
	"rm -rf /",
	"rm -rf c:",
	"del /f /s /q c:\\",
	"rmdir /s /q c:\\",
	":(){ :|:& };:",
	"drop database",
	"truncate table",
}

// IsDangerousCommand checks if a shell command contains known catastrophic patterns.
func (sm *SecurityManager) IsDangerousCommand(command string) (bool, string) {
	cmdLower := strings.ToLower(command)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(cmdLower, pattern) {
			return true, fmt.Sprintf("command matches blocked destructive pattern: '%s'", pattern)
		}
	}
	return false, ""
}

// Safe commands list for Guarded mode (read-only, diagnostic, and local test verification)
var safeCommandPrefixes = []string{
	"go version", "go env", "go test", "go vet", "go build", "go list",
	"git status", "git diff", "git log", "git branch",
	"dir", "ls", "type ", "cat ", "echo ", "pwd", "where ", "which ",
	"npm test", "cargo test", "python -m unittest",
}

var guardedRiskyKeywords = []string{
	"rm ", "del ", "rmdir", "remove-item", "drop", "truncate", "shutdown",
	"git push", "git reset --hard", "kill", "stop-process",
	">", ">>", "|", "curl", "wget", "Invoke-WebRequest",
}

// isSafeGuardedCommand checks if a command is benign enough to auto-approve in Guarded mode.
func isSafeGuardedCommand(command string) bool {
	cmdTrim := strings.TrimSpace(strings.ToLower(command))

	// Check for any risky substrings
	for _, risky := range guardedRiskyKeywords {
		if strings.Contains(cmdTrim, risky) {
			return false
		}
	}

	// Check if matches known safe prefixes
	for _, prefix := range safeCommandPrefixes {
		if strings.HasPrefix(cmdTrim, prefix) {
			return true
		}
	}

	return false
}

// CheckApprovalRequirement determines if a tool execution needs human confirmation.
func (sm *SecurityManager) CheckApprovalRequirement(toolName string, args map[string]interface{}) (bool, string) {
	sm.mu.RLock()
	mode := sm.Mode
	sm.mu.RUnlock()

	// 1. 完全不审核模式 (Auto Mode) -> 直接放行
	if mode == models.ApprovalModeAuto {
		return false, ""
	}

	// 2. 人工审核模式 (Manual Mode) -> 任何写操作与 Shell 执行均拦截审核
	if mode == models.ApprovalModeManual {
		switch toolName {
		case "cmd_exec":
			cmd, _ := args["command"].(string)
			return true, fmt.Sprintf("【人工审核】终端命令执行: %s", cmd)
		case "file_write":
			path, _ := args["path"].(string)
			return true, fmt.Sprintf("【人工审核】写入/覆写文件: %s", path)
		case "file_edit":
			path, _ := args["path"].(string)
			return true, fmt.Sprintf("【人工审核】精准差异编辑文件: %s", path)
		case "file_delete":
			path, _ := args["path"].(string)
			return true, fmt.Sprintf("【人工审核】删除文件: %s", path)
		default:
			return false, ""
		}
	}

	// 3. 有限审核模式 (Guarded Mode) -> 智能评估风险
	switch toolName {
	case "cmd_exec":
		cmd, _ := args["command"].(string)
		if isSafeGuardedCommand(cmd) {
			return false, "" // 安全只读/测试命令自动放行
		}
		return true, fmt.Sprintf("【有限审核】潜在高危/非安全只读终端命令: %s", cmd)

	case "file_write":
		path, _ := args["path"].(string)
		validPath, err := sm.ValidatePath(path)
		if err == nil {
			if _, statErr := os.Stat(validPath); statErr == nil {
				// 文件已存在，覆写属于破坏性行为，需审核
				return true, fmt.Sprintf("【有限审核】覆写已有文件: %s", path)
			}
		}
		// 新建文件在工作区内自动放行
		return false, ""

	case "file_edit":
		// 精准代码替换在工作区内通常是日常开发，Guarded 模式默认放行
		return false, ""

	case "file_delete":
		path, _ := args["path"].(string)
		return true, fmt.Sprintf("【有限审核】高危删除操作: %s", path)

	default:
		return false, ""
	}
}

// HITLManager coordinates Human-in-the-Loop approval requests.
type HITLManager struct {
	mu       sync.Mutex
	pending  map[string]chan models.ApprovalResponse
	requests map[string]*models.ApprovalRequest
	onNotify func(req *models.ApprovalRequest)
}

// NewHITLManager creates a new HITLManager.
func NewHITLManager() *HITLManager {
	return &HITLManager{
		pending:  make(map[string]chan models.ApprovalResponse),
		requests: make(map[string]*models.ApprovalRequest),
	}
}

// SetNotifyHandler sets the callback to broadcast approval requests to the UI.
func (h *HITLManager) SetNotifyHandler(handler func(req *models.ApprovalRequest)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onNotify = handler
}

// RequestApproval creates a pending approval token and waits for human input or timeout.
func (h *HITLManager) RequestApproval(ctx context.Context, tool string, desc string, payload interface{}) (bool, string, error) {
	tokenBytes := make([]byte, 16)
	_, _ = rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	respChan := make(chan models.ApprovalResponse, 1)

	req := &models.ApprovalRequest{
		Token:       token,
		Tool:        tool,
		Description: desc,
		Payload:     payload,
		CreatedAt:   time.Now(),
	}

	h.mu.Lock()
	h.pending[token] = respChan
	h.requests[token] = req
	notifier := h.onNotify
	h.mu.Unlock()

	if notifier != nil {
		notifier(req)
	}

	select {
	case <-ctx.Done():
		h.cleanup(token)
		return false, "任务已取消", ctx.Err()
	case <-time.After(10 * time.Minute):
		h.cleanup(token)
		return false, "操作超时（超过 10 分钟未处理，已自动拒绝）", nil
	case resp := <-respChan:
		h.cleanup(token)
		return resp.Approved, resp.Reason, nil
	}
}

// SubmitDecision receives decision from UI.
func (h *HITLManager) SubmitDecision(token string, approved bool, reason string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch, exists := h.pending[token]
	if !exists {
		return false
	}

	ch <- models.ApprovalResponse{
		Token:    token,
		Approved: approved,
		Reason:   reason,
	}
	return true
}

func (h *HITLManager) cleanup(token string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.pending, token)
	delete(h.requests, token)
}

package security

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"agent/pkg/models"
)

func TestValidatePath_BoundaryEnforcement(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sec_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sec := NewSecurityManager(tempDir, models.ApprovalModeManual)

	// Valid subpath
	validPath := filepath.Join(tempDir, "sub", "test.txt")
	resolved, err := sec.ValidatePath(validPath)
	if err != nil {
		t.Errorf("expected valid path, got err: %v", err)
	}
	if resolved != filepath.Clean(validPath) {
		t.Errorf("expected %s, got %s", validPath, resolved)
	}

	// Traversal breakout attempt
	escapePath := filepath.Join(tempDir, "..", "outside.txt")
	_, err = sec.ValidatePath(escapePath)
	if err == nil {
		t.Errorf("expected path escape error for %s, but passed", escapePath)
	}
}

func TestDynamicWorkspaceSwitching(t *testing.T) {
	tempDir1, _ := os.MkdirTemp("", "sec_dir1_*")
	tempDir2, _ := os.MkdirTemp("", "sec_dir2_*")
	defer os.RemoveAll(tempDir1)
	defer os.RemoveAll(tempDir2)

	sec := NewSecurityManager(tempDir1, models.ApprovalModeGuarded)

	// In dir 1
	p1 := filepath.Join(tempDir1, "file1.txt")
	if _, err := sec.ValidatePath(p1); err != nil {
		t.Fatalf("failed to validate path in dir1: %v", err)
	}

	// In dir 2 should fail initially
	p2 := filepath.Join(tempDir2, "file2.txt")
	if _, err := sec.ValidatePath(p2); err == nil {
		t.Fatal("expected dir2 path to fail before switch, but passed")
	}

	// Switch to dir 2
	newWs, err := sec.SetWorkspaceDir(tempDir2)
	if err != nil {
		t.Fatalf("failed to switch workspace: %v", err)
	}
	if newWs != tempDir2 {
		t.Errorf("expected %s, got %s", tempDir2, newWs)
	}

	// Now p2 should pass, p1 should fail
	if _, err := sec.ValidatePath(p2); err != nil {
		t.Errorf("expected dir2 path to pass after switch, got: %v", err)
	}
	if _, err := sec.ValidatePath(p1); err == nil {
		t.Errorf("expected dir1 path to fail after switch, but passed")
	}
}

func TestAntigravityThreeTierApproval(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "sec_tier_*")
	defer os.RemoveAll(tempDir)

	sec := NewSecurityManager(tempDir, models.ApprovalModeAuto)

	// 1. Auto mode: safe & write commands don't need approval
	need, _ := sec.CheckApprovalRequirement("cmd_exec", map[string]interface{}{"command": "git push origin main"})
	if need {
		t.Errorf("expected auto mode to not require approval, got true")
	}

	// 2. Guarded mode: safe commands auto-pass, risky commands need approval
	sec.SetApprovalMode(models.ApprovalModeGuarded)

	// Safe command: go test
	need, _ = sec.CheckApprovalRequirement("cmd_exec", map[string]interface{}{"command": "go test -v ./..."})
	if need {
		t.Errorf("expected safe command 'go test' to auto-pass in guarded mode, but got need=true")
	}

	// Risky command: del or git push
	need, _ = sec.CheckApprovalRequirement("cmd_exec", map[string]interface{}{"command": "git push --force origin main"})
	if !need {
		t.Errorf("expected risky command 'git push' to require approval in guarded mode, but got need=false")
	}

	// File edit inside workspace auto-passes in guarded mode
	need, _ = sec.CheckApprovalRequirement("file_edit", map[string]interface{}{"path": "main.go"})
	if need {
		t.Errorf("expected file_edit to auto-pass in guarded mode, got need=true")
	}

	// 3. Manual mode: even safe commands and edits require approval
	sec.SetApprovalMode(models.ApprovalModeManual)
	need, _ = sec.CheckApprovalRequirement("file_edit", map[string]interface{}{"path": "main.go"})
	if !need {
		t.Errorf("expected manual mode to require approval for file_edit, got false")
	}
	need, _ = sec.CheckApprovalRequirement("cmd_exec", map[string]interface{}{"command": "go version"})
	if !need {
		t.Errorf("expected manual mode to require approval for safe cmd, got false")
	}
}

func TestIsDangerousCommand(t *testing.T) {
	sec := NewSecurityManager(".", models.ApprovalModeGuarded)

	dangerousCases := []string{
		"format c:",
		"rm -rf /",
		"rm -rf c:",
		"del /f /s /q c:\\",
		"diskpart",
	}

	for _, cmd := range dangerousCases {
		isDanger, reason := sec.IsDangerousCommand(cmd)
		if !isDanger {
			t.Errorf("expected command '%s' to be flagged as dangerous, but got false", cmd)
		}
		if reason == "" {
			t.Errorf("expected danger explanation for '%s'", cmd)
		}
	}

	safeCases := []string{
		"echo hello",
		"go version",
		"git status",
		"dir",
		"npm test",
	}

	for _, cmd := range safeCases {
		isDanger, _ := sec.IsDangerousCommand(cmd)
		if isDanger {
			t.Errorf("expected command '%s' to be safe, but was flagged", cmd)
		}
	}
}

func TestHITLManager_Decision(t *testing.T) {
	hitl := NewHITLManager()

	var receivedToken string
	hitl.SetNotifyHandler(func(req *models.ApprovalRequest) {
		receivedToken = req.Token
	})

	// Async request
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	doneCh := make(chan bool)
	go func() {
		// Wait a fraction of a second, then approve
		for i := 0; i < 20; i++ {
			time.Sleep(20 * time.Millisecond)
			hitl.mu.Lock()
			for k := range hitl.pending {
				receivedToken = k
			}
			hitl.mu.Unlock()
			if receivedToken != "" {
				hitl.SubmitDecision(receivedToken, true, "Approved for testing")
				break
			}
		}
	}()

	approved, reason, err := hitl.RequestApproval(ctx, "cmd_exec", "test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !approved {
		t.Errorf("expected approved=true, got false (reason: %s)", reason)
	}
	close(doneCh)
}

func TestChinesePath_ValidationAndNormalization(t *testing.T) {
	tempBase := t.TempDir()
	chineseWs := filepath.Join(tempBase, "中文测试工作空间_研发")
	_ = os.MkdirAll(chineseWs, 0755)

	sec := NewSecurityManager(chineseWs, models.ApprovalModeGuarded)

	// 1. Relative file with Chinese name
	p1, err := sec.ValidatePath("测试文档.txt")
	if err != nil {
		t.Fatalf("expected valid Chinese relative path, got err: %v", err)
	}
	expected1 := filepath.Join(chineseWs, "测试文档.txt")
	if p1 != expected1 {
		t.Errorf("expected %s, got %s", expected1, p1)
	}

	// 2. Path with leading forward slash (common LLM generation pattern "/index.html")
	p2, err := sec.ValidatePath("/index.html")
	if err != nil {
		t.Fatalf("expected valid path for leading slash, got err: %v", err)
	}
	expected2 := filepath.Join(chineseWs, "index.html")
	if p2 != expected2 {
		t.Errorf("expected %s, got %s", expected2, p2)
	}

	// 3. Nested path with forward slashes and Chinese
	p3, err := sec.ValidatePath("/src/components/组件.vue")
	if err != nil {
		t.Fatalf("expected valid path for nested forward slash, got err: %v", err)
	}
	expected3 := filepath.Join(chineseWs, "src", "components", "组件.vue")
	if p3 != expected3 {
		t.Errorf("expected %s, got %s", expected3, p3)
	}

	// 4. Absolute path with Chinese
	absTarget := filepath.Join(chineseWs, "nested", "file.go")
	p4, err := sec.ValidatePath(absTarget)
	if err != nil {
		t.Fatalf("expected valid absolute path, got err: %v", err)
	}
	if p4 != absTarget {
		t.Errorf("expected %s, got %s", absTarget, p4)
	}

	// 5. Escaping path should be rejected
	escapeTarget := filepath.Join(tempBase, "outside_file.txt")
	_, err = sec.ValidatePath(escapeTarget)
	if err == nil {
		t.Errorf("expected error for escaping path %s, but passed", escapeTarget)
	}
}

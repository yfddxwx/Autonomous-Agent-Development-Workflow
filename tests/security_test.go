package tests

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"agent/pkg/models"
	"agent/pkg/security"
	"agent/pkg/tools"
)

func TestSecurity_BoundaryIsolation(t *testing.T) {
	tempWorkspace, err := os.MkdirTemp("", "sec_workspace_*")
	if err != nil {
		t.Fatalf("failed to create temp workspace: %v", err)
	}
	defer os.RemoveAll(tempWorkspace)

	sec := security.NewSecurityManager(tempWorkspace, models.ApprovalModeManual)
	reg := tools.NewRegistry()
	tools.RegisterBuiltinTools(reg, sec)

	ctx := context.Background()

	// 1. Attempt to write outside workspace via relative path
	_, err = reg.Execute(ctx, "file_write", `{"path":"../../system_escape.txt","content":"pwned"}`)
	if err == nil {
		t.Fatal("expected security refusal on path traversal write, but succeeded")
	}
	if !strings.Contains(err.Error(), "security violation") {
		t.Errorf("expected security violation error message, got: %v", err)
	}

	// 2. Attempt to read outside workspace
	_, err = reg.Execute(ctx, "file_read", `{"path":"../some_parent_file.txt"}`)
	if err == nil {
		t.Fatal("expected security refusal on path traversal read, but succeeded")
	}

	// 3. Attempt to run blocked catastrophic command
	_, err = reg.Execute(ctx, "cmd_exec", `{"command":"format c:"}`)
	if err == nil {
		t.Fatal("expected command execution refusal for 'format c:', but succeeded")
	}
	if !strings.Contains(err.Error(), "security refusal") {
		t.Errorf("expected security refusal message, got: %v", err)
	}
}

func TestSecurity_HITLApprovalRejection(t *testing.T) {
	tempWorkspace, err := os.MkdirTemp("", "sec_hitl_test_*")
	if err != nil {
		t.Fatalf("failed to create temp workspace: %v", err)
	}
	defer os.RemoveAll(tempWorkspace)

	sec := security.NewSecurityManager(tempWorkspace, models.ApprovalModeManual)
	hitl := sec.GetHITLManager()

	// 1. Verify approval policy check
	needsApproval, desc := sec.CheckApprovalRequirement("cmd_exec", map[string]interface{}{"command": "git push --force"})
	if !needsApproval {
		t.Errorf("expected cmd_exec to require approval, got false")
	}
	if !strings.Contains(desc, "git push --force") {
		t.Errorf("expected desc to contain command text, got: %s", desc)
	}

	// 2. Test rejection handling
	hitl.SetNotifyHandler(func(req *models.ApprovalRequest) {
		// Asynchronously reject the request
		go func(token string) {
			time.Sleep(30 * time.Millisecond)
			hitl.SubmitDecision(token, false, "拒绝执行此操作")
		}(req.Token)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	approved, reason, err := hitl.RequestApproval(ctx, "cmd_exec", desc, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if approved {
		t.Errorf("expected approved=false, but got true")
	}
	if !strings.Contains(reason, "拒绝执行此操作") {
		t.Errorf("expected reason to contain rejection message, got: %s", reason)
	}
}

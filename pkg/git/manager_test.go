package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"agent/pkg/models"
	"agent/pkg/security"
)

func TestGitManager_NonGitRepo(t *testing.T) {
	tempDir := t.TempDir()
	sec := security.NewSecurityManager(tempDir, models.ApprovalModeGuarded)
	mgr := NewManager(sec)
	ctx := context.Background()

	if mgr.IsGitRepo(ctx, tempDir) {
		t.Fatalf("expected non-git repo for empty temp dir")
	}

	_, err := mgr.GetStatus(ctx, tempDir)
	if err == nil {
		t.Fatalf("expected error from GetStatus on non-git dir")
	}

	_, err = mgr.Rollback(ctx, tempDir)
	if err == nil {
		t.Fatalf("expected error from Rollback on non-git dir")
	}
}

func TestGitManager_Lifecycle(t *testing.T) {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed, skipping git manager test")
	}

	tempDir := t.TempDir()
	sec := security.NewSecurityManager(tempDir, models.ApprovalModeGuarded)
	mgr := NewManager(sec)
	ctx := context.Background()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git init: %v", err)
	}

	// Configure local dummy user for git commit
	_ = exec.Command("git", "-C", tempDir, "config", "user.name", "TestAgent").Run()
	_ = exec.Command("git", "-C", tempDir, "config", "user.email", "test@agent.local").Run()

	if !mgr.IsGitRepo(ctx, tempDir) {
		t.Fatalf("expected IsGitRepo to be true after git init")
	}

	// Create initial file & commit
	initialFile := filepath.Join(tempDir, "hello.txt")
	if err := os.WriteFile(initialFile, []byte("version 1\n"), 0644); err != nil {
		t.Fatalf("failed to write initial file: %v", err)
	}

	_ = exec.Command("git", "-C", tempDir, "add", ".").Run()
	_ = exec.Command("git", "-C", tempDir, "commit", "-m", "initial commit").Run()

	// Status should be clean
	status, err := mgr.GetStatus(ctx, tempDir)
	if err != nil {
		t.Fatalf("GetStatus error: %v", err)
	}
	if strings.TrimSpace(status) != "" {
		t.Fatalf("expected clean status, got: %s", status)
	}

	// Modify tracked file and add untracked file
	_ = os.WriteFile(initialFile, []byte("version 2 modified\n"), 0644)
	untrackedFile := filepath.Join(tempDir, "untracked.txt")
	_ = os.WriteFile(untrackedFile, []byte("i am untracked\n"), 0644)

	// Status should now show changes
	status, err = mgr.GetStatus(ctx, tempDir)
	if err != nil {
		t.Fatalf("GetStatus error after edit: %v", err)
	}
	if !strings.Contains(status, "hello.txt") || !strings.Contains(status, "untracked.txt") {
		t.Fatalf("expected status to show changes, got: %s", status)
	}

	// Execute Rollback
	rollbackMsg, err := mgr.Rollback(ctx, tempDir)
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
	if !strings.Contains(rollbackMsg, "已成功回滚") {
		t.Fatalf("unexpected rollback message: %s", rollbackMsg)
	}

	// Verify tracked file is restored
	restoredContent, err := os.ReadFile(initialFile)
	if err != nil {
		t.Fatalf("failed to read restored file: %v", err)
	}
	if strings.TrimSpace(string(restoredContent)) != "version 1" {
		t.Fatalf("file content not restored, got: %s", string(restoredContent))
	}

	// Verify untracked file is removed
	if _, err := os.Stat(untrackedFile); !os.IsNotExist(err) {
		t.Fatalf("untracked file still exists after rollback")
	}

	// Verify status is clean again
	status, err = mgr.GetStatus(ctx, tempDir)
	if err != nil {
		t.Fatalf("GetStatus error after rollback: %v", err)
	}
	if strings.TrimSpace(status) != "" {
		t.Fatalf("expected clean status after rollback, got: %s", status)
	}
}

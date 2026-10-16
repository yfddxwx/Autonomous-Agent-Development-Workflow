package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"agent/pkg/security"
)

// Manager coordinates Git repository operations and rollback checkpoints.
type Manager struct {
	sec *security.SecurityManager
}

// NewManager creates a new Git manager.
func NewManager(sec *security.SecurityManager) *Manager {
	return &Manager{sec: sec}
}

// runGit executes a git command strictly within the target directory.
func (m *Manager) runGit(ctx context.Context, dir string, args ...string) (string, error) {
	validDir, err := m.sec.ValidatePath(dir)
	if err != nil {
		return "", err
	}

	execCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	gitArgs := append([]string{"-c", "core.quotepath=false"}, args...)
	cmd := exec.CommandContext(execCtx, "git", gitArgs...)
	cmd.Dir = validDir

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git error (%v): %s", err, strings.TrimSpace(errBuf.String()))
	}

	return strings.TrimSpace(outBuf.String()), nil
}

// IsGitRepo checks if the specified directory is inside a git repository.
func (m *Manager) IsGitRepo(ctx context.Context, dir string) bool {
	validDir, err := m.sec.ValidatePath(dir)
	if err != nil {
		return false
	}
	gitPath := filepath.Join(validDir, ".git")
	if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
		return true
	}
	out, err := m.runGit(ctx, dir, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// GetStatus returns the current porcelain status of the repository.
func (m *Manager) GetStatus(ctx context.Context, dir string) (string, error) {
	if !m.IsGitRepo(ctx, dir) {
		return "", errors.New("当前目录不是 Git 仓库")
	}
	return m.runGit(ctx, dir, "status", "--porcelain")
}

// CreateCheckpoint creates a lightweight stash or tag checkpoint.
func (m *Manager) CreateCheckpoint(ctx context.Context, dir string, name string) (string, error) {
	if !m.IsGitRepo(ctx, dir) {
		return "", errors.New("当前目录不是 Git 仓库，无法创建检查点")
	}

	status, err := m.GetStatus(ctx, dir)
	if err != nil {
		return "", err
	}
	if status == "" {
		return "clean", nil // Nothing to checkpoint
	}

	// Create stash without clearing workspace
	commitHash, err := m.runGit(ctx, dir, "stash", "create", name)
	if err != nil {
		return "", err
	}
	if commitHash == "" {
		return "clean", nil
	}
	return commitHash, nil
}

// Rollback restores the git repository by discarding unstaged and untracked changes.
func (m *Manager) Rollback(ctx context.Context, dir string) (string, error) {
	if !m.IsGitRepo(ctx, dir) {
		return "", errors.New("当前目录不是 Git 仓库，无法回滚代码")
	}

	// 1. Revert tracked changes
	_, err := m.runGit(ctx, dir, "checkout", "--", ".")
	if err != nil {
		return "", fmt.Errorf("git checkout 恢复追踪文件失败: %w", err)
	}

	// 2. Clean untracked files
	cleanOut, err := m.runGit(ctx, dir, "clean", "-fd")
	if err != nil {
		return "", fmt.Errorf("git clean 清理新增文件失败: %w", err)
	}

	return fmt.Sprintf("已成功回滚本轮所有文件改动。\n%s", cleanOut), nil
}

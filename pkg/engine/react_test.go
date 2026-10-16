package engine

import (
	"context"
	"strings"
	"testing"

	"agent/pkg/models"
	"agent/pkg/security"
	"agent/pkg/tools"
)

type mockLLMClient struct {
	model string
}

func (m *mockLLMClient) ChatStream(ctx context.Context, messages []models.Message, tools []models.ToolDefinition, onChunk func(models.StreamChunk) error) (*models.Message, error) {
	return &models.Message{
		Role:    models.RoleAssistant,
		Content: "mock reply",
	}, nil
}

func (m *mockLLMClient) GetModelName() string {
	return m.model
}

func (m *mockLLMClient) SetReasoningEffort(effort string) {}

func TestSessionTree_BranchingAndNonDestructiveRollback(t *testing.T) {
	tempDir := t.TempDir()
	sec := security.NewSecurityManager(tempDir, models.ApprovalModeAuto)
	reg := tools.NewRegistry()
	client := &mockLLMClient{model: "deepseek-chat"}

	engine := NewReActEngine(client, reg, sec, 5)

	// Turn 1
	engine.recordSnapshot("创建 README", "已创建 README", []models.Message{{Role: models.RoleUser, Content: "创建 README"}}, nil)
	snaps := engine.GetSnapshots()
	if len(snaps) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(snaps))
	}
	if snaps[0].Branch != "main" {
		t.Fatalf("expected main branch, got %s", snaps[0].Branch)
	}
	if snaps[0].ParentID != "" {
		t.Fatalf("expected empty parent for root turn, got %s", snaps[0].ParentID)
	}

	// Turn 2
	engine.recordSnapshot("实现主服务", "已实现主服务", []models.Message{{Role: models.RoleUser, Content: "实现主服务"}}, nil)
	snaps = engine.GetSnapshots()
	if len(snaps) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(snaps))
	}
	if snaps[1].Branch != "main" {
		t.Fatalf("expected main branch, got %s", snaps[1].Branch)
	}
	if snaps[1].ParentID != snaps[0].ID {
		t.Fatalf("expected parent to be %s, got %s", snaps[0].ID, snaps[1].ParentID)
	}

	// Turn 3
	engine.recordSnapshot("添加测试用例", "测试已全部通过", []models.Message{{Role: models.RoleUser, Content: "添加测试用例"}}, nil)
	snaps = engine.GetSnapshots()
	if len(snaps) != 3 {
		t.Fatalf("expected 3 snapshots, got %d", len(snaps))
	}

	// Time travel back to Turn 1 (index 0)
	rolledMsgs, err := engine.RollbackToTurn(0)
	if err != nil {
		t.Fatalf("RollbackToTurn failed: %v", err)
	}
	if len(rolledMsgs) != 1 {
		t.Fatalf("expected 1 message in rolled turn, got %d", len(rolledMsgs))
	}

	// Non-destructive check: snapshots must NOT be truncated!
	snapsAfterRollback := engine.GetSnapshots()
	if len(snapsAfterRollback) != 3 {
		t.Fatalf("expected snapshots array to remain 3 (non-destructive), got %d", len(snapsAfterRollback))
	}

	// Fork a new branch from Turn 1!
	engine.recordSnapshot("尝试采用微服务架构重构", "微服务方案分支启动", []models.Message{{Role: models.RoleUser, Content: "尝试采用微服务架构重构"}}, nil)
	snapsAfterFork := engine.GetSnapshots()
	if len(snapsAfterFork) != 4 {
		t.Fatalf("expected 4 snapshots after fork, got %d", len(snapsAfterFork))
	}

	forkSnap := snapsAfterFork[3]
	if !strings.HasPrefix(forkSnap.Branch, "fork-") {
		t.Fatalf("expected branch to start with 'fork-', got %s", forkSnap.Branch)
	}
	if forkSnap.ParentID != snaps[0].ID {
		t.Fatalf("expected fork parent to be %s, got %s", snaps[0].ID, forkSnap.ParentID)
	}
}

func TestTelemetry_CalculateCost(t *testing.T) {
	tempDir := t.TempDir()
	sec := security.NewSecurityManager(tempDir, models.ApprovalModeAuto)
	reg := tools.NewRegistry()

	// DeepSeek pricing (2 CNY / 8 CNY per 1M)
	dsEngine := NewReActEngine(&mockLLMClient{model: "deepseek-chat"}, reg, sec, 5)
	cny, usd := dsEngine.calculateCost(1000000, 1000000)
	if cny != 10.0 {
		t.Fatalf("expected 10.0 CNY, got %f", cny)
	}
	if usd <= 0 {
		t.Fatalf("expected positive USD cost, got %f", usd)
	}

	// Ollama free
	ollamaEngine := NewReActEngine(&mockLLMClient{model: "ollama/llama3"}, reg, sec, 5)
	cnyOllama, usdOllama := ollamaEngine.calculateCost(1000000, 1000000)
	if cnyOllama != 0.0 || usdOllama != 0.0 {
		t.Fatalf("expected 0 cost for ollama, got cny=%f, usd=%f", cnyOllama, usdOllama)
	}
}

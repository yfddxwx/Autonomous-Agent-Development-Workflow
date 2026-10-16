package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"agent/pkg/engine"
	"agent/pkg/llm"
	"agent/pkg/models"
	"agent/pkg/security"
	"agent/pkg/tools"
)

func TestReActEngine_AutonomousLifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "agent_integration_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	var requestCount int32

	// Setup Mock LLM server mimicking OpenAI streaming SSE
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("expected flusher")
		}

		if count == 1 {
			// Step 1: Model decides to call file_write
			toolCallChunk := `{"choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"call_123","type":"function","function":{"name":"file_write","arguments":"{\"path\":\"build_info.json\",\"content\":\"{\\\"version\\\":\\\"1.0.0\\\"}\"}"}}]}}]}`
			fmt.Fprintf(w, "data: %s\n\n", toolCallChunk)
			flusher.Flush()

			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
		} else {
			// Step 2: Model reflects on observation and outputs final answer
			contentChunk := `{"choices":[{"index":0,"delta":{"content":"文件 build_info.json 已成功生成并验证通过。"}}]}`
			fmt.Fprintf(w, "data: %s\n\n", contentChunk)
			flusher.Flush()

			fmt.Fprintf(w, "data: [DONE]\n\n")
			flusher.Flush()
		}
	}))
	defer mockServer.Close()

	// Initialize components
	sec := security.NewSecurityManager(tempDir, models.ApprovalModeAuto)
	reg := tools.NewRegistry()
	tools.RegisterBuiltinTools(reg, sec)

	client := llm.NewOpenAIClient(mockServer.URL, "test-key", "mock-model")
	reactEngine := engine.NewReActEngine(client, reg, sec, 5)

	var receivedEvents []string
	onEvent := func(eventType string, data map[string]interface{}) {
		receivedEvents = append(receivedEvents, eventType)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	history, err := reactEngine.Run(ctx, "创建 build_info.json 文件", nil, onEvent)
	if err != nil {
		t.Fatalf("ReActEngine.Run error: %v", err)
	}

	// Verify file was actually written to disk
	createdFilePath := filepath.Join(tempDir, "build_info.json")
	fileBytes, err := os.ReadFile(createdFilePath)
	if err != nil {
		t.Fatalf("expected file to exist at %s: %v", createdFilePath, err)
	}

	var data map[string]string
	if err := json.Unmarshal(fileBytes, &data); err != nil || data["version"] != "1.0.0" {
		t.Errorf("file content mismatch, got: %s", string(fileBytes))
	}

	// Verify events flow
	eventsJoined := strings.Join(receivedEvents, ",")
	if !strings.Contains(eventsJoined, "tool_call") || !strings.Contains(eventsJoined, "tool_result") || !strings.Contains(eventsJoined, "done") {
		t.Errorf("missing expected events in stream: %s", eventsJoined)
	}

	// Verify conversation messages length
	if len(history) < 3 {
		t.Errorf("expected history to contain system, user, tool_call, tool_result, assistant final message. Got %d messages", len(history))
	}
}

package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"agent/pkg/models"
	"agent/pkg/security"
)

func TestFileTools_Lifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tool_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sec := security.NewSecurityManager(tempDir, models.ApprovalModeAuto)

	writer := NewFileWriteTool(sec)
	reader := NewFileReadTool(sec)
	lister := NewFileListTool(sec)

	ctx := context.Background()

	// 1. Write file
	targetFile := filepath.Join("subdir", "sample.txt")
	writeContent := "Line 1: Hello World\nLine 2: Autonomous Agent\nLine 3: Test Passed\n"
	writeOut, err := writer.Execute(ctx, map[string]interface{}{
		"path":      targetFile,
		"content":   writeContent,
		"overwrite": true,
	})
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	if !strings.Contains(writeOut, "Successfully written") {
		t.Errorf("unexpected write output: %s", writeOut)
	}

	// 2. Read file with line slicing
	readOut, err := reader.Execute(ctx, map[string]interface{}{
		"path":       targetFile,
		"start_line": float64(2),
		"end_line":   float64(2),
	})
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if !strings.Contains(readOut, "Autonomous Agent") || strings.Contains(readOut, "Hello World") {
		t.Errorf("slice read unexpected output: %s", readOut)
	}

	// 3. List directory
	listOut, err := lister.Execute(ctx, map[string]interface{}{
		"dir":       ".",
		"recursive": true,
	})
	if err != nil {
		t.Fatalf("failed to list files: %v", err)
	}
	if !strings.Contains(listOut, "subdir") || !strings.Contains(listOut, "sample.txt") {
		t.Errorf("list output missing written files: %s", listOut)
	}
}

func TestFileEditTool_SurgicalReplacement(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "edit_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sec := security.NewSecurityManager(tempDir, models.ApprovalModeAuto)
	writer := NewFileWriteTool(sec)
	editor := NewFileEditTool(sec)

	ctx := context.Background()

	// Create base file
	targetFile := "main.go"
	initContent := "package main\n\nfunc main() {\n\tprintln(\"Hello Old World\")\n}\n"
	_, err = writer.Execute(ctx, map[string]interface{}{
		"path":    targetFile,
		"content": initContent,
	})
	if err != nil {
		t.Fatalf("failed to write initial file: %v", err)
	}

	// 1. Surgical single replacement
	editOut, err := editor.Execute(ctx, map[string]interface{}{
		"path":                targetFile,
		"target_content":      "println(\"Hello Old World\")",
		"replacement_content": "println(\"Hello Pi Agent\")",
	})
	if err != nil {
		t.Fatalf("file_edit failed: %v", err)
	}
	if !strings.Contains(editOut, "Successfully edited") {
		t.Errorf("unexpected output: %s", editOut)
	}

	// Verify content
	updatedBytes, err := os.ReadFile(filepath.Join(tempDir, targetFile))
	if err != nil || !strings.Contains(string(updatedBytes), "Hello Pi Agent") {
		t.Fatalf("content was not updated properly: %s", string(updatedBytes))
	}

	// 2. Error when target content doesn't exist
	_, err = editor.Execute(ctx, map[string]interface{}{
		"path":                targetFile,
		"target_content":      "non_existent_text_block",
		"replacement_content": "replacement",
	})
	if err == nil {
		t.Fatal("expected error for non-existent target content, got nil")
	}
}

func TestCmdExecTool_Execution(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cmd_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sec := security.NewSecurityManager(tempDir, models.ApprovalModeAuto)
	cmdTool := NewCmdExecTool(sec)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Execute simple echo
	out, err := cmdTool.Execute(ctx, map[string]interface{}{
		"command":         "Write-Output 'Agent Test OK'",
		"timeout_seconds": float64(10),
	})
	if err != nil {
		t.Fatalf("cmd_exec failed: %v", err)
	}
	if !strings.Contains(out, "Agent Test OK") {
		t.Errorf("expected output to contain 'Agent Test OK', got: %s", out)
	}
}

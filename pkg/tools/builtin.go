package tools

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"agent/pkg/security"
)

// StreamCallbackKey is the context key for streaming terminal chunks in real time.
type StreamCallbackKey struct{}

// ContextWithStreamCallback injects a stream callback into context.
func ContextWithStreamCallback(ctx context.Context, cb func(chunk string)) context.Context {
	return context.WithValue(ctx, StreamCallbackKey{}, cb)
}

// FileReadTool reads file content.
type FileReadTool struct {
	sec *security.SecurityManager
}

func NewFileReadTool(sec *security.SecurityManager) *FileReadTool {
	return &FileReadTool{sec: sec}
}

func (t *FileReadTool) Name() string { return "file_read" }
func (t *FileReadTool) Description() string {
	return "安全读取指定文件的文本内容，支持行号切片范围读取"
}
func (t *FileReadTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "文件相对或绝对路径",
			},
			"start_line": map[string]interface{}{
				"type":        "integer",
				"description": "起始行号（从 1 开始，可选）",
			},
			"end_line": map[string]interface{}{
				"type":        "integer",
				"description": "结束行号（可选）",
			},
		},
		"required": []string{"path"},
	}
}

func (t *FileReadTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return "", errors.New("missing 'path' argument")
	}

	validPath, err := t.sec.ValidatePath(path)
	if err != nil {
		return "", err
	}

	file, err := os.Open(validPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	startLine := 1
	if sl, ok := args["start_line"].(float64); ok && int(sl) > 0 {
		startLine = int(sl)
	}
	endLine := -1
	if el, ok := args["end_line"].(float64); ok && int(el) > 0 {
		endLine = int(el)
	}

	scanner := bufio.NewScanner(file)
	var sb strings.Builder
	currentLine := 1

	for scanner.Scan() {
		if currentLine >= startLine && (endLine == -1 || currentLine <= endLine) {
			cleanLine := strings.ToValidUTF8(scanner.Text(), "")
			sb.WriteString(fmt.Sprintf("%4d | %s\n", currentLine, cleanLine))
		}
		currentLine++
		if endLine != -1 && currentLine > endLine {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	if sb.Len() == 0 {
		return "(空文件或超出指定行范围)", nil
	}
	return sb.String(), nil
}

// FileWriteTool writes or overwrites file content.
type FileWriteTool struct {
	sec *security.SecurityManager
}

func NewFileWriteTool(sec *security.SecurityManager) *FileWriteTool {
	return &FileWriteTool{sec: sec}
}

func (t *FileWriteTool) Name() string { return "file_write" }
func (t *FileWriteTool) Description() string {
	return "在指定路径创建或覆写文件，自动创建不存在的父目录"
}
func (t *FileWriteTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "文件路径",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "写入的完整文本内容",
			},
			"overwrite": map[string]interface{}{
				"type":        "boolean",
				"description": "如文件存在是否覆写（默认为 true）",
			},
		},
		"required": []string{"path", "content"},
	}
}

func (t *FileWriteTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return "", errors.New("missing 'path' argument")
	}
	content, ok := args["content"].(string)
	if !ok {
		return "", errors.New("missing 'content' argument")
	}

	validPath, err := t.sec.ValidatePath(path)
	if err != nil {
		return "", err
	}

	overwrite := true
	if ow, ok := args["overwrite"].(bool); ok {
		overwrite = ow
	}

	if !overwrite {
		if _, err := os.Stat(validPath); err == nil {
			return "", fmt.Errorf("file '%s' already exists and overwrite is false", path)
		}
	}

	dir := filepath.Dir(validPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}

	if err := os.WriteFile(validPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Successfully written %d bytes to %s", len(content), path), nil
}

// FileListTool lists directory contents.
type FileListTool struct {
	sec *security.SecurityManager
}

func NewFileListTool(sec *security.SecurityManager) *FileListTool {
	return &FileListTool{sec: sec}
}

func (t *FileListTool) Name() string { return "file_list" }
func (t *FileListTool) Description() string {
	return "列出指定目录下的文件与文件夹结构"
}
func (t *FileListTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"dir": map[string]interface{}{
				"type":        "string",
				"description": "目录路径（默认为工作区根目录 '.'）",
			},
			"recursive": map[string]interface{}{
				"type":        "boolean",
				"description": "是否递归扫描子目录",
			},
		},
	}
}

func (t *FileListTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	dir := "."
	if d, ok := args["dir"].(string); ok && d != "" {
		dir = d
	}

	validDir, err := t.sec.ValidatePath(dir)
	if err != nil {
		return "", err
	}

	recursive := false
	if rec, ok := args["recursive"].(bool); ok {
		recursive = rec
	}

	var sb strings.Builder
	if !recursive {
		entries, err := os.ReadDir(validDir)
		if err != nil {
			return "", fmt.Errorf("failed to read dir: %w", err)
		}
		for _, e := range entries {
			typeStr := "FILE"
			if e.IsDir() {
				typeStr = "DIR "
			}
			info, _ := e.Info()
			size := int64(0)
			if info != nil {
				size = info.Size()
			}
			sb.WriteString(fmt.Sprintf("[%s] %-25s %8d B\n", typeStr, e.Name(), size))
		}
	} else {
		err := filepath.WalkDir(validDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			rel, _ := filepath.Rel(validDir, path)
			if rel == "." {
				return nil
			}
			// Skip .git or heavy build dirs
			if d.IsDir() && (d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == "vendor") {
				return filepath.SkipDir
			}
			typeStr := "FILE"
			if d.IsDir() {
				typeStr = "DIR "
			}
			sb.WriteString(fmt.Sprintf("[%s] %s\n", typeStr, filepath.ToSlash(rel)))
			return nil
		})
		if err != nil {
			return "", err
		}
	}

	if sb.Len() == 0 {
		return "(目录为空)", nil
	}
	return sb.String(), nil
}

// CmdExecTool executes system commands with timeout and cancellation.
type CmdExecTool struct {
	sec *security.SecurityManager
}

func NewCmdExecTool(sec *security.SecurityManager) *CmdExecTool {
	return &CmdExecTool{sec: sec}
}

func (t *CmdExecTool) Name() string { return "cmd_exec" }
func (t *CmdExecTool) Description() string {
	return "在本地工作空间内执行 Shell/终端命令并获取输出"
}
func (t *CmdExecTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "要执行的命令行指令",
			},
			"timeout_seconds": map[string]interface{}{
				"type":        "integer",
				"description": "超时时间（秒，默认 60）",
			},
		},
		"required": []string{"command"},
	}
}

func (t *CmdExecTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	command, ok := args["command"].(string)
	if !ok || strings.TrimSpace(command) == "" {
		return "", errors.New("missing or empty 'command' argument")
	}

	// Security check for dangerous destructive patterns
	if isDanger, reason := t.sec.IsDangerousCommand(command); isDanger {
		return "", fmt.Errorf("security refusal: %s", reason)
	}

	timeoutSec := 60
	if to, ok := args["timeout_seconds"].(float64); ok && int(to) > 0 {
		timeoutSec = int(to)
	}

	execCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		psCmd := "[Console]::OutputEncoding = [System.Text.Encoding]::UTF8; [Console]::InputEncoding = [System.Text.Encoding]::UTF8; $OutputEncoding = [System.Text.Encoding]::UTF8; " + command
		cmd = exec.CommandContext(execCtx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", psCmd)
	} else {
		cmd = exec.CommandContext(execCtx, "sh", "-c", command)
	}

	cmd.Dir = t.sec.GetWorkspaceDir()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to open stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to open stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start command: %w", err)
	}

	streamer, _ := ctx.Value(StreamCallbackKey{}).(func(chunk string))

	var outBuf bytes.Buffer
	var bufMu sync.Mutex

	readPipe := func(r io.Reader) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := strings.ToValidUTF8(scanner.Text(), "") + "\n"
			bufMu.Lock()
			outBuf.WriteString(line)
			bufMu.Unlock()
			if streamer != nil {
				streamer(line)
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		readPipe(stdoutPipe)
	}()
	go func() {
		defer wg.Done()
		readPipe(stderrPipe)
	}()

	wg.Wait()
	err = cmd.Wait()
	output := outBuf.String()

	if execCtx.Err() == context.DeadlineExceeded {
		return output, fmt.Errorf("command timed out after %d seconds", timeoutSec)
	}
	if execCtx.Err() == context.Canceled {
		return output, errors.New("command execution canceled by user")
	}

	if err != nil {
		return fmt.Sprintf("Exit Error: %v\nOutput:\n%s", err, output), nil
	}

	if strings.TrimSpace(output) == "" {
		return "(命令执行成功，无终端标准输出)", nil
	}

	// Limit length if too long
	const maxOutput = 100000
	if len(output) > maxOutput {
		output = output[:maxOutput] + "\n...(输出过长已截断)..."
	}

	return output, nil
}

// FileEditTool surgically replaces target text inside a file (Pi-style precise diff editing).
type FileEditTool struct {
	sec *security.SecurityManager
}

func NewFileEditTool(sec *security.SecurityManager) *FileEditTool {
	return &FileEditTool{sec: sec}
}

func (t *FileEditTool) Name() string { return "file_edit" }
func (t *FileEditTool) Description() string {
	return "精准差异替换指定文件中的文本块（Pi 风格的高效编辑工具，无需重写整个文件，大幅降低 Token 消耗并保留未修改代码）"
}
func (t *FileEditTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "待修改的目标文件路径",
			},
			"target_content": map[string]interface{}{
				"type":        "string",
				"description": "文件中待被替换的原始内容块（必须精确匹配，包含空格缩进与换行）",
			},
			"replacement_content": map[string]interface{}{
				"type":        "string",
				"description": "用于替换的新内容块",
			},
			"allow_multiple": map[string]interface{}{
				"type":        "boolean",
				"description": "如文件中存在多处匹配，是否全部替换（默认 false，如匹配到多处则报错提示提供更精准上下文）",
			},
		},
		"required": []string{"path", "target_content", "replacement_content"},
	}
}

func (t *FileEditTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok || path == "" {
		return "", errors.New("missing 'path' argument")
	}
	targetContent, ok := args["target_content"].(string)
	if !ok || targetContent == "" {
		return "", errors.New("missing 'target_content' argument")
	}
	replacementContent, ok := args["replacement_content"].(string)
	if !ok {
		return "", errors.New("missing 'replacement_content' argument")
	}

	validPath, err := t.sec.ValidatePath(path)
	if err != nil {
		return "", err
	}

	fileBytes, err := os.ReadFile(validPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file '%s': %w", path, err)
	}

	originalText := string(fileBytes)
	occurrences := strings.Count(originalText, targetContent)
	if occurrences == 0 {
		return "", fmt.Errorf("target_content not found in '%s'. Please verify whitespace, indentation, or use file_read to inspect current content", path)
	}

	allowMultiple := false
	if am, ok := args["allow_multiple"].(bool); ok {
		allowMultiple = am
	}

	if occurrences > 1 && !allowMultiple {
		return "", fmt.Errorf("target_content matches %d locations in '%s'. To prevent unintended replacements, provide a larger surrounding block or set allow_multiple=true", occurrences, path)
	}

	var newText string
	if allowMultiple {
		newText = strings.ReplaceAll(originalText, targetContent, replacementContent)
	} else {
		newText = strings.Replace(originalText, targetContent, replacementContent, 1)
	}

	if err := os.WriteFile(validPath, []byte(newText), 0644); err != nil {
		return "", fmt.Errorf("failed to save updated file '%s': %w", path, err)
	}

	return fmt.Sprintf("Successfully edited '%s' (%d location(s) replaced)", path, occurrences), nil
}

// RegisterBuiltinTools registers all standard tools into registry.
func RegisterBuiltinTools(r *Registry, sec *security.SecurityManager) {
	r.Register(NewFileReadTool(sec))
	r.Register(NewFileWriteTool(sec))
	r.Register(NewFileEditTool(sec))
	r.Register(NewFileListTool(sec))
	r.Register(NewCmdExecTool(sec))
}

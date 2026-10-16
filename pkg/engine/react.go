package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"agent/pkg/llm"
	"agent/pkg/models"
	"agent/pkg/security"
	"agent/pkg/skills"
	"agent/pkg/tools"
)

// EventCallback defines the signature for real-time lifecycle event broadcasting.
type EventCallback func(eventType string, data map[string]interface{})

// ReActEngine coordinates autonomous planning, execution, and reflection.
type ReActEngine struct {
	client       llm.Client
	registry     *tools.Registry
	sec          *security.SecurityManager
	maxSteps     int
	skillsLoader *skills.Loader
	snapshots    []models.SessionSnapshot
	activeSnapID string
	branchCount  int
	totalPrompt  int
	totalCompl   int
}

// NewReActEngine creates a new ReAct engine.
func NewReActEngine(client llm.Client, registry *tools.Registry, sec *security.SecurityManager, maxSteps int) *ReActEngine {
	if maxSteps <= 0 {
		maxSteps = 15
	}
	loader := skills.NewLoader()
	loader.EnsureDefaultSkills(sec.GetWorkspaceDir())

	return &ReActEngine{
		client:       client,
		registry:     registry,
		sec:          sec,
		maxSteps:     maxSteps,
		skillsLoader: loader,
		snapshots:    make([]models.SessionSnapshot, 0),
	}
}

// SetClient updates the active LLM client (for runtime model switching).
func (e *ReActEngine) SetClient(client llm.Client) {
	e.client = client
}

// GetClient returns the current active LLM client.
func (e *ReActEngine) GetClient() llm.Client {
	return e.client
}

// GetSnapshots returns the current session snapshots list.
func (e *ReActEngine) GetSnapshots() []models.SessionSnapshot {
	return e.snapshots
}

// SetSnapshots sets snapshots from a restored session.
func (e *ReActEngine) SetSnapshots(snapshots []models.SessionSnapshot) {
	e.snapshots = snapshots
	if len(snapshots) > 0 {
		e.activeSnapID = snapshots[len(snapshots)-1].ID
	}
}

// ClearSnapshots clears the snapshot history.
func (e *ReActEngine) ClearSnapshots() {
	e.snapshots = make([]models.SessionSnapshot, 0)
	e.activeSnapID = ""
	e.branchCount = 0
	e.totalPrompt = 0
	e.totalCompl = 0
}

// SetTelemetry restores token usage telemetry.
func (e *ReActEngine) SetTelemetry(promptTok, complTok int) {
	e.totalPrompt = promptTok
	e.totalCompl = complTok
}

func (e *ReActEngine) calculateCost(promptTok, complTok int) (float64, float64) {
	modelName := strings.ToLower(e.client.GetModelName())
	var promptRateCNY, complRateCNY float64

	switch {
	case strings.Contains(modelName, "gpt-4"):
		promptRateCNY = 18.0
		complRateCNY = 72.0
	case strings.Contains(modelName, "claude"):
		promptRateCNY = 21.0
		complRateCNY = 105.0
	case strings.Contains(modelName, "gemini"):
		promptRateCNY = 9.0
		complRateCNY = 36.0
	case strings.Contains(modelName, "ollama") || strings.Contains(modelName, "llama"):
		promptRateCNY = 0.0
		complRateCNY = 0.0
	case strings.Contains(modelName, "flash"): // DeepSeek V4.1 Flash (ultra cost-effective)
		promptRateCNY = 1.0
		complRateCNY = 4.0
	default: // DeepSeek V3 / R1 / Custom
		promptRateCNY = 2.0
		complRateCNY = 8.0
	}

	costCNY := (float64(promptTok)*promptRateCNY + float64(complTok)*complRateCNY) / 1000000.0
	costUSD := costCNY / 7.25
	return costCNY, costUSD
}

// GetTelemetry returns current token counters, CNY and USD cost.
func (e *ReActEngine) GetTelemetry() (int, int, int, float64, float64) {
	total := e.totalPrompt + e.totalCompl
	costCNY, costUSD := e.calculateCost(e.totalPrompt, e.totalCompl)
	return e.totalPrompt, e.totalCompl, total, costCNY, costUSD
}

// RollbackToTurn activates a target snapshot node for non-destructive time travel without pruning branches.
func (e *ReActEngine) RollbackToTurn(turnIndex int) ([]models.Message, error) {
	if turnIndex < 0 || turnIndex >= len(e.snapshots) {
		return nil, fmt.Errorf("invalid turn index: %d", turnIndex)
	}
	targetSnap := e.snapshots[turnIndex]
	e.activeSnapID = targetSnap.ID
	return targetSnap.Messages, nil
}

// BuildSystemPrompt constructs the initial agent instructions.
func (e *ReActEngine) BuildSystemPrompt() string {
	wsDir := e.sec.GetWorkspaceDir()
	basePrompt := fmt.Sprintf(`你是一个运行在本地机器上的资深全栈研发自主 Agent，具备多 Agent 团队协同与自主编排能力。
你的工作空间根目录为: %s

【核心行为准则】:
1. 你具备完全自主的本地工程研发与调试能力，可以利用提供的工具直接读取、改写文件并运行 Shell 终端指令。
2. 遇到任务时，先进行严谨的架构思考与步骤拆解，然后再调用工具执行操作。
3. 【Pi 风格精准编辑】：修改已有代码时，优先使用 file_edit 进行精准块级差异替换，避免使用 file_write 重写整个大文件，以极大节省 Token 并保障代码安全性。仅在新建文件时使用 file_write。
4. 编码完成后，务必使用终端工具运行编译、运行单元测试或检查输出，确保交付物质量无误。
5. 遇到报错或测试失败时，切忌盲目猜测，应主动读取出错日志并反思原因，精准修正代码。
6. 遵守最小权限与工作区边界约束，严禁执行越界破坏性操作。
7. 输出使用中文，保持专业、严谨且条理清晰。

【深度思考协议 (Deep Thinking Protocol)】:
在执行任何工具调用或输出回复前，你必须首先在 <think>...</think> 闭合标签内展开真实、严密、深度的内心推理：
1. 深入剖析用户的核心研发诉求与边界条件；
2. 结合当前工作区目录状态与前序步骤执行结果，推演技术架构；
3. 规划下一步工具调用的精准参数，并推演可能遇到的编译报错与异常修复策略。
无论任务简单还是复杂，严禁跳过思考直接调用工具。你的内心推理必须完整包裹在 <think>...</think> 闭合标签内。

【虚拟工程团队协同与任务复杂度自适应机制】:
你有权根据用户任务的复杂度自主选择工作模式：
- 模式一【单兵直驱执行 (Direct Solo)】：针对简单明确的任务（如修改少量代码、列出目录、查询状态等），无需分工，直接调用工具快速完成。
- 模式二【虚拟工程团队协同 (Multi-Agent Team Orchestration)】：针对复杂的综合性研发工程（从零构建完整应用、多模块架构重构、全套功能及测试闭环开发），由你自主承担并切换虚拟工程团队的各个专业角色，按照严谨的工程标准推进：
  * [Role: 架构规划师 (Architect)]: 梳理需求、设计系统架构与模块契约；
  * [Role: 全栈研发工程师 (Engineer)]: 编写核心代码，使用 file_edit 微创精准差异编辑；
  * [Role: 质量与安全评审员 (QA & Reviewer)]: 运行终端命令执行编译测试与安全性校验，确保 100%% 验收通过。

当你决定激活团队协同模式并在不同角色间推进时，请在思考或回复开头明确标示角色（如：👥 [团队角色: 架构规划师] 或 👥 [团队角色: 全栈研发工程师] 或 👥 [团队角色: 质量与安全评审员]），以便系统实时向用户可视化展示团队分工进度。`, wsDir)

	// Append dynamic workspace skills if available
	e.skillsLoader.EnsureDefaultSkills(wsDir)
	skillsList, err := e.skillsLoader.LoadSkills(wsDir)
	if err == nil && len(skillsList) > 0 {
		basePrompt += e.skillsLoader.FormatSkillsPrompt(skillsList)
	} else {
		basePrompt += fmt.Sprintf(`

【本地专属技能即插即用解析器 (Agent Skills 标准支持)】:
你所在的运行环境已完全挂载并激活了 Agent Skills 解析引擎，支持读取工作区根目录下的标准技能定义：
1. 技能支持路径：%s/.agent/skills/ 与 %s/skills/
2. 技能组织规范：支持标准 Agent Skills 目录结构（包含 SKILL.md 或 skill.md 的独立文件夹，支持 YAML Frontmatter 元数据如 name、description），以及扁平的 .md 规范文档。
3. 当前工作区暂未检测到其他自定义技能。用户随时可将 SKILL.md 技能文件夹放入上述目录，引擎将在下一轮会话即插即用自动挂载注入上下文。`, wsDir, wsDir)
	}

	return basePrompt
}

// compactMessages ensures history does not grow unboundedly (Pi-style context management)
func compactMessages(msgs []models.Message, maxCount int) []models.Message {
	if len(msgs) <= maxCount {
		return msgs
	}
	result := make([]models.Message, len(msgs))
	copy(result, msgs)

	cutoff := len(result) - 12
	for i := 1; i < cutoff; i++ {
		if result[i].Role == models.RoleTool && len(result[i].Content) > 200 {
			result[i].Content = result[i].Content[:100] + "\n...(历史工具结果已由上下文引擎自动压缩)..."
		}
	}
	return result
}

// Run executes the autonomous ReAct cycle for a user prompt.
func (e *ReActEngine) Run(ctx context.Context, userPrompt string, conversationHistory []models.Message, onEvent EventCallback) ([]models.Message, error) {
	toolDefs := e.registry.GetDefinitions()

	messages := make([]models.Message, 0, len(conversationHistory)+10)
	if len(conversationHistory) == 0 {
		messages = append(messages, models.Message{
			Role:    models.RoleSystem,
			Content: e.BuildSystemPrompt(),
		})
	} else {
		// Update system prompt with latest workspace
		compacted := compactMessages(conversationHistory, 20)
		if len(compacted) > 0 && compacted[0].Role == models.RoleSystem {
			compacted[0].Content = e.BuildSystemPrompt()
		}
		messages = append(messages, compacted...)
	}

	messages = append(messages, models.Message{
		Role:    models.RoleUser,
		Content: userPrompt,
	})

	for step := 1; step <= e.maxSteps; step++ {
		select {
		case <-ctx.Done():
			if onEvent != nil {
				onEvent("status", map[string]interface{}{"state": "idle", "text": "任务已取消"})
			}
			return messages, ctx.Err()
		default:
		}

		if onEvent != nil {
			onEvent("status", map[string]interface{}{
				"state": "thinking",
				"text":  fmt.Sprintf("正在思考 (第 %d/%d 轮)...", step, e.maxSteps),
			})
		}

		// Stream completion
		assistantMsg, err := e.client.ChatStream(ctx, messages, toolDefs, func(chunk models.StreamChunk) error {
			if onEvent != nil {
				if chunk.Thought != "" {
					onEvent("thought", map[string]interface{}{"content": chunk.Thought})
					if strings.Contains(chunk.Thought, "架构") || strings.Contains(chunk.Thought, "Architect") {
						onEvent("team_role", map[string]interface{}{"role": "architect", "name": "架构规划师", "icon": "fa-compass-drafting"})
					} else if strings.Contains(chunk.Thought, "开发") || strings.Contains(chunk.Thought, "Engineer") || strings.Contains(chunk.Thought, "编码") {
						onEvent("team_role", map[string]interface{}{"role": "engineer", "name": "全栈开发工程师", "icon": "fa-code"})
					} else if strings.Contains(chunk.Thought, "测试") || strings.Contains(chunk.Thought, "安全") || strings.Contains(chunk.Thought, "QA") {
						onEvent("team_role", map[string]interface{}{"role": "qa", "name": "质量与安全评审员", "icon": "fa-vial-circle-check"})
					}
				}
				if chunk.Content != "" {
					onEvent("content", map[string]interface{}{"content": chunk.Content})
				}
			}
			return nil
		})

		if err != nil {
			if onEvent != nil {
				onEvent("error", map[string]interface{}{"error": fmt.Sprintf("模型生成失败: %v", err)})
			}
			return messages, err
		}

		// Calculate Token Telemetry
		if assistantMsg.Usage != nil && assistantMsg.Usage.TotalTokens > 0 {
			e.totalPrompt += assistantMsg.Usage.PromptTokens
			e.totalCompl += assistantMsg.Usage.CompletionTokens
		} else {
			stepPromptChars := 0
			for _, m := range messages {
				stepPromptChars += len(m.Content) + len(m.Thought)
			}
			stepComplChars := len(assistantMsg.Content) + len(assistantMsg.Thought)

			pTok := stepPromptChars / 3
			if pTok < 10 && stepPromptChars > 0 {
				pTok = 10
			}
			cTok := stepComplChars / 3
			if cTok < 5 && stepComplChars > 0 {
				cTok = 5
			}
			e.totalPrompt += pTok
			e.totalCompl += cTok
		}

		totalToks := e.totalPrompt + e.totalCompl
		costCNY, costUSD := e.calculateCost(e.totalPrompt, e.totalCompl)

		if onEvent != nil {
			onEvent("telemetry", map[string]interface{}{
				"prompt_tokens":      e.totalPrompt,
				"completion_tokens":  e.totalCompl,
				"total_tokens":       totalToks,
				"cost_cny":           costCNY,
				"cost_usd":           costUSD,
				"estimated_cost_usd": costUSD,
			})
		}

		messages = append(messages, *assistantMsg)

		// Check if there are tool calls to execute
		if len(assistantMsg.ToolCalls) == 0 {
			// Finished - create snapshot
			e.recordSnapshot(userPrompt, assistantMsg.Content, messages, onEvent)
			if onEvent != nil {
				onEvent("done", map[string]interface{}{"step": step})
			}
			return messages, nil
		}

		// Execute each tool call
		for _, tc := range assistantMsg.ToolCalls {
			select {
			case <-ctx.Done():
				return messages, ctx.Err()
			default:
			}

			if onEvent != nil {
				onEvent("status", map[string]interface{}{
					"state": "executing",
					"text":  fmt.Sprintf("正在调用工具: %s", tc.Function.Name),
				})
				onEvent("tool_call", map[string]interface{}{
					"call_id": tc.ID,
					"tool":    tc.Function.Name,
					"args":    tc.Function.Arguments,
				})
			}

			var parsedArgs map[string]interface{}
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &parsedArgs)

			// Security & HITL Approval Check
			needsApproval, desc := e.sec.CheckApprovalRequirement(tc.Function.Name, parsedArgs)
			var toolOutput string
			var toolErr error

			toolCtx := tools.ContextWithStreamCallback(ctx, func(chunk string) {
				if onEvent != nil {
					onEvent("cmd_stream", map[string]interface{}{
						"call_id": tc.ID,
						"chunk":   chunk,
					})
				}
			})

			if needsApproval {
				if onEvent != nil {
					onEvent("status", map[string]interface{}{
						"state": "awaiting_approval",
						"text":  "需要人工审核确认",
					})
				}
				approved, reason, hitlErr := e.sec.GetHITLManager().RequestApproval(ctx, tc.Function.Name, desc, parsedArgs)
				if hitlErr != nil {
					toolErr = hitlErr
					toolOutput = fmt.Sprintf("审批中断或取消: %v", hitlErr)
				} else if !approved {
					toolOutput = fmt.Sprintf("【操作被用户在界面上明确拒绝】: %s。请反思该操作为何被拒绝，并调整执行方案或向用户说明。", reason)
				} else {
					// User approved, proceed with execution
					toolOutput, toolErr = e.registry.Execute(toolCtx, tc.Function.Name, tc.Function.Arguments)
				}
			} else {
				toolOutput, toolErr = e.registry.Execute(toolCtx, tc.Function.Name, tc.Function.Arguments)
			}

			errStr := ""
			if toolErr != nil {
				errStr = toolErr.Error()
				if toolOutput == "" {
					toolOutput = fmt.Sprintf("执行错误: %v", toolErr)
				}
			}

			if onEvent != nil {
				onEvent("tool_result", map[string]interface{}{
					"call_id": tc.ID,
					"output":  toolOutput,
					"error":   errStr,
				})
			}

			// Feed tool observation back to context
			messages = append(messages, models.Message{
				Role:       models.RoleTool,
				Content:    toolOutput,
				ToolCallID: tc.ID,
			})
		}
	}

	e.recordSnapshot(userPrompt, "已达最大循环轮数", messages, onEvent)
	if onEvent != nil {
		onEvent("done", map[string]interface{}{
			"warning": fmt.Sprintf("已达最大单次自主步骤上限 (%d 轮)", e.maxSteps),
		})
	}

	return messages, nil
}

func (e *ReActEngine) recordSnapshot(prompt, rawSummary string, msgs []models.Message, onEvent EventCallback) {
	summary := rawSummary
	if len(summary) > 60 {
		summary = summary[:60] + "..."
	}
	if summary == "" {
		summary = "任务步骤执行完成"
	}

	snapID := fmt.Sprintf("snap_%d", time.Now().UnixNano())
	turnIdx := len(e.snapshots)
	parentID := e.activeSnapID
	branch := "main"

	if parentID != "" {
		hasChildren := false
		var parentSnap *models.SessionSnapshot
		for i := range e.snapshots {
			if e.snapshots[i].ID == parentID {
				parentSnap = &e.snapshots[i]
			}
			if e.snapshots[i].ParentID == parentID {
				hasChildren = true
			}
		}

		if hasChildren {
			e.branchCount++
			branch = fmt.Sprintf("fork-%d", e.branchCount)
		} else if parentSnap != nil && parentSnap.Branch != "" {
			branch = parentSnap.Branch
		}
	}

	snapshot := models.SessionSnapshot{
		ID:        snapID,
		ParentID:  parentID,
		Branch:    branch,
		TurnIndex: turnIdx,
		Prompt:    prompt,
		Summary:   summary,
		CreatedAt: time.Now(),
		Messages:  msgs,
	}
	e.snapshots = append(e.snapshots, snapshot)
	e.activeSnapID = snapID

	if onEvent != nil {
		onEvent("snapshot_created", map[string]interface{}{
			"id":         snapID,
			"parent_id":  parentID,
			"branch":     branch,
			"turn_index": turnIdx,
			"prompt":     prompt,
			"summary":    summary,
		})
	}
}

package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"agent/pkg/models"
)

// Client is the interface for communicating with language models.
type Client interface {
	ChatStream(ctx context.Context, messages []models.Message, tools []models.ToolDefinition, onChunk func(chunk models.StreamChunk) error) (*models.Message, error)
	GetModelName() string
	SetReasoningEffort(effort string)
}

// OpenAIClient implements Client for any OpenAI-compatible provider (OpenAI, DeepSeek, Ollama, Gemini-OpenAI, etc.).
type OpenAIClient struct {
	BaseURL         string
	APIKey          string
	ModelName       string
	ReasoningEffort string
	HTTPClient      *http.Client
}

// NewOpenAIClient creates a new OpenAI compatible client.
func NewOpenAIClient(baseURL, apiKey, modelName string) *OpenAIClient {
	cleanURL := strings.TrimRight(baseURL, "/")
	if cleanURL == "" {
		cleanURL = "https://api.deepseek.com/v1"
	}
	if modelName == "" {
		modelName = "deepseek-chat"
	}
	return &OpenAIClient{
		BaseURL:   cleanURL,
		APIKey:    apiKey,
		ModelName: modelName,
		HTTPClient: &http.Client{
			Timeout: 180 * time.Second,
		},
	}
}

func (c *OpenAIClient) GetModelName() string {
	return c.ModelName
}

func (c *OpenAIClient) SetReasoningEffort(effort string) {
	c.ReasoningEffort = effort
}

// StreamOptions requests token usage metadata in streaming mode.
type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

// ChatRequest represents the payload sent to /chat/completions.
type ChatRequest struct {
	Model           string                  `json:"model"`
	Messages        []interface{}           `json:"messages"`
	Tools           []models.ToolDefinition `json:"tools,omitempty"`
	Stream          bool                    `json:"stream"`
	StreamOptions   *StreamOptions          `json:"stream_options,omitempty"`
	ReasoningEffort string                  `json:"reasoning_effort,omitempty"`
}

// OpenAIDelta represents an incremental piece in streaming response.
type OpenAIDelta struct {
	Role             string            `json:"role,omitempty"`
	Content          string            `json:"content,omitempty"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCalls        []OpenAIToolDelta `json:"tool_calls,omitempty"`
}

type OpenAIToolDelta struct {
	Index    int                 `json:"index"`
	ID       string              `json:"id,omitempty"`
	Type     string              `json:"type,omitempty"`
	Function OpenAIFunctionDelta `json:"function,omitempty"`
}

type OpenAIFunctionDelta struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type OpenAIStreamChunk struct {
	Choices []struct {
		Index        int         `json:"index"`
		Delta        OpenAIDelta `json:"delta"`
		FinishReason *string     `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

// ChatStream executes streaming chat completion.
func (c *OpenAIClient) ChatStream(ctx context.Context, messages []models.Message, tools []models.ToolDefinition, onChunk func(chunk models.StreamChunk) error) (*models.Message, error) {
	endpoint := fmt.Sprintf("%s/chat/completions", c.BaseURL)

	// Format messages
	rawMessages := make([]interface{}, 0, len(messages))
	for _, m := range messages {
		msgMap := map[string]interface{}{
			"role": string(m.Role),
		}
		if m.Content != "" {
			msgMap["content"] = m.Content
		}
		// Preserve reasoning_content on assistant messages (required by DeepSeek / reasoning models)
		if m.Role == models.RoleAssistant && m.Thought != "" {
			msgMap["reasoning_content"] = m.Thought
		}
		if len(m.ToolCalls) > 0 {
			msgMap["tool_calls"] = m.ToolCalls
		}
		if m.ToolCallID != "" {
			msgMap["tool_call_id"] = m.ToolCallID
		}
		rawMessages = append(rawMessages, msgMap)
	}

	reqBody := ChatRequest{
		Model:    c.ModelName,
		Messages: rawMessages,
		Tools:    tools,
		Stream:   true,
		StreamOptions: &StreamOptions{
			IncludeUsage: true,
		},
	}

	effort := strings.ToLower(strings.TrimSpace(c.ReasoningEffort))
	if effort == "off" || effort == "none" {
		reqBody.ReasoningEffort = ""
	} else if effort != "" {
		switch effort {
		case "minimal":
			reqBody.ReasoningEffort = "low"
		case "xhigh", "max":
			reqBody.ReasoningEffort = "high"
		default:
			reqBody.ReasoningEffort = effort
		}
	} else {
		// Default fallback if reasoning model and effort not explicitly provided
		lowerModel := strings.ToLower(c.ModelName)
		lowerURL := strings.ToLower(c.BaseURL)
		if strings.Contains(lowerURL, "deepseek") || strings.Contains(lowerModel, "deepseek") || strings.Contains(lowerModel, "flash") || strings.Contains(lowerModel, "reasoner") || strings.Contains(lowerModel, "r1") {
			reqBody.ReasoningEffort = "high"
		}
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

		lowerURL := strings.ToLower(c.BaseURL)
		if strings.Contains(lowerURL, "anthropic") || strings.Contains(lowerURL, "minimax") || strings.Contains(lowerURL, "kimi.com") {
			httpReq.Header.Set("x-api-key", c.APIKey)
			httpReq.Header.Set("anthropic-version", "2023-06-01")
		}
		if strings.Contains(lowerURL, "azure") {
			httpReq.Header.Set("api-key", c.APIKey)
		}
		if strings.Contains(lowerURL, "githubcopilot.com") {
			httpReq.Header.Set("Editor-Version", "vscode/1.98.0")
			httpReq.Header.Set("Editor-Plugin-Version", "copilot/1.250.0")
		}
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	scanner := bufio.NewScanner(resp.Body)
	// Allow large tokens buffer
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var assembledContent strings.Builder
	var assembledThought strings.Builder
	toolCallMap := make(map[int]*models.ToolCall)
	var reportedUsage *models.TokenUsage
	inThink := false

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var chunk OpenAIStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			if chunk.Usage != nil {
				reportedUsage = &models.TokenUsage{
					PromptTokens:     chunk.Usage.PromptTokens,
					CompletionTokens: chunk.Usage.CompletionTokens,
					TotalTokens:      chunk.Usage.TotalTokens,
				}
			}

			if len(chunk.Choices) == 0 {
				continue
			}

			delta := chunk.Choices[0].Delta
			var streamChunk models.StreamChunk

			if delta.ReasoningContent != "" {
				assembledThought.WriteString(delta.ReasoningContent)
				streamChunk.Thought += delta.ReasoningContent
			}

			if delta.Content != "" {
				raw := delta.Content
				if strings.Contains(raw, "<think>") {
					parts := strings.Split(raw, "<think>")
					if parts[0] != "" {
						assembledContent.WriteString(parts[0])
						streamChunk.Content += parts[0]
					}
					inThink = true
					if len(parts) > 1 {
						raw = parts[1]
					} else {
						raw = ""
					}
				}

				if inThink {
					if strings.Contains(raw, "</think>") {
						parts := strings.Split(raw, "</think>")
						assembledThought.WriteString(parts[0])
						streamChunk.Thought += parts[0]
						inThink = false
						if len(parts) > 1 && parts[1] != "" {
							assembledContent.WriteString(parts[1])
							streamChunk.Content += parts[1]
						}
					} else {
						assembledThought.WriteString(raw)
						streamChunk.Thought += raw
					}
				} else if raw != "" {
					assembledContent.WriteString(raw)
					streamChunk.Content += raw
				}
			}

			for _, tcDelta := range delta.ToolCalls {
				tc, exists := toolCallMap[tcDelta.Index]
				if !exists {
					tc = &models.ToolCall{
						Type: "function",
					}
					toolCallMap[tcDelta.Index] = tc
				}
				if tcDelta.ID != "" {
					tc.ID = tcDelta.ID
				}
				if tcDelta.Function.Name != "" {
					tc.Function.Name += tcDelta.Function.Name
				}
				if tcDelta.Function.Arguments != "" {
					tc.Function.Arguments += tcDelta.Function.Arguments
				}
			}

			if onChunk != nil && (streamChunk.Content != "" || streamChunk.Thought != "") {
				if err := onChunk(streamChunk); err != nil {
					return nil, err
				}
			}
		}
	}

	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("stream read error: %w", err)
	}

	// Reassemble final tool calls in order
	finalToolCalls := make([]models.ToolCall, 0, len(toolCallMap))
	for i := 0; i < len(toolCallMap); i++ {
		if tc, ok := toolCallMap[i]; ok {
			finalToolCalls = append(finalToolCalls, *tc)
		}
	}

	resultMsg := &models.Message{
		Role:      models.RoleAssistant,
		Content:   assembledContent.String(),
		Thought:   assembledThought.String(),
		ToolCalls: finalToolCalls,
		Usage:     reportedUsage,
	}

	return resultMsg, nil
}

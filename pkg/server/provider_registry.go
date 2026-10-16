package server

import "strings"

// Code generated from official Pi Agent registry (@earendil-works/pi-ai). DO NOT EDIT MANUALLY.

type RegistryModel struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Reasoning     bool   `json:"reasoning,omitempty"`
	ContextWindow int    `json:"context_window,omitempty"`
	MaxTokens     int    `json:"max_tokens,omitempty"`
	Description   string `json:"description,omitempty"`
}

type ProviderMeta struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Type     string          `json:"type"`     // "account" or "api_key"
	Category string          `json:"category"` // "china", "global", "chips", "gateway", "local"
	BaseURL  string          `json:"base_url"`
	AuthEnv  string          `json:"auth_env"`
	AuthType string          `json:"auth_type"` // "bearer", "x-api-key", "token", "none"
	Models   []RegistryModel `json:"models"`
}

var OfficialProviders = []ProviderMeta{
	{
		ID:       "amazon-bedrock",
		Name:     "AWS credentials or bearer token",
		Type:     "api_key",
		Category: "gateway",
		BaseURL:  "https://bedrock-runtime.us-east-1.amazonaws.com",
		AuthEnv:  "",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "amazon.nova-2-lite-v1:0", Name: "Nova 2 Lite", Reasoning: true, ContextWindow: 1000000, MaxTokens: 65535, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "amazon.nova-lite-v1:0", Name: "Nova Lite", Reasoning: false, ContextWindow: 300000, MaxTokens: 10000, Description: ""},
			{ID: "amazon.nova-micro-v1:0", Name: "Nova Micro", Reasoning: false, ContextWindow: 128000, MaxTokens: 10000, Description: ""},
			{ID: "amazon.nova-pro-v1:0", Name: "Nova Pro", Reasoning: false, ContextWindow: 300000, MaxTokens: 10000, Description: ""},
			{ID: "anthropic.claude-fable-5", Name: "Claude Fable 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic.claude-fable-5-1", Name: "Claude Fable 5.1", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic.claude-haiku-4-5-20251001-v1:0", Name: "Claude Haiku 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic.claude-haiku-5-5", Name: "Claude Haiku 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic.claude-opus-4-1-20250805-v1:0", Name: "Claude Opus 4.1", Reasoning: true, ContextWindow: 200000, MaxTokens: 32000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic.claude-opus-4-5-20251101-v1:0", Name: "Claude Opus 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic.claude-opus-4-6-v1", Name: "Claude Opus 4.6", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic.claude-opus-4-7", Name: "Claude Opus 4.7", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic.claude-opus-4-8", Name: "Claude Opus 4.8", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic.claude-opus-5-5", Name: "Claude Opus 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic.claude-sonnet-4-5-20250929-v1:0", Name: "Claude Sonnet 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "ant-ling",
		Name:     "Ant Ling",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://api.ant-ling.com/v1",
		AuthEnv:  "[\"ANT_LING_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "Ling-2.6-1T", Name: "Ling 2.6 1T", Reasoning: false, ContextWindow: 262144, MaxTokens: 65536, Description: ""},
			{ID: "Ling-2.6-flash", Name: "Ling 2.6 Flash", Reasoning: false, ContextWindow: 262144, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "Ring-2.6-1T", Name: "Ring 2.6 1T", Reasoning: true, ContextWindow: 262144, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "anthropic",
		Name:     "Anthropic API key",
		Type:     "account",
		Category: "global",
		BaseURL:  "https://api.anthropic.com",
		AuthEnv:  "",
		AuthType: "x-api-key",
		Models: []RegistryModel{
			{ID: "claude-fable-5", Name: "Claude Fable 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-fable-5-1", Name: "Claude Fable 5.1", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-4-5", Name: "Claude Haiku 4.5 (latest)", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-4-5-20251001", Name: "Claude Haiku 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-5-5", Name: "Claude Haiku 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-5", Name: "Claude Opus 4.5 (latest)", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-5-20251101", Name: "Claude Opus 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-6", Name: "Claude Opus 4.6", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-7", Name: "Claude Opus 4.7", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-8", Name: "Claude Opus 4.8", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-5", Name: "Claude Opus 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-5-5", Name: "Claude Opus 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-4-5", Name: "Claude Sonnet 4.5 (latest)", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-4-5-20250929", Name: "Claude Sonnet 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-4-6", Name: "Claude Sonnet 4.6", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "azure",
		Name:     "Azure",
		Type:     "api_key",
		Category: "gateway",
		BaseURL:  "",
		AuthEnv:  "[\"AZURE_OPENAI_API_KEY\"]",
		AuthType: "azure-api-key",
		Models: []RegistryModel{
			{ID: "gpt-4", Name: "GPT-4", Reasoning: false, ContextWindow: 8192, MaxTokens: 8192, Description: ""},
			{ID: "gpt-4-turbo", Name: "GPT-4 Turbo", Reasoning: false, ContextWindow: 128000, MaxTokens: 4096, Description: ""},
			{ID: "gpt-4.1", Name: "GPT-4.1", Reasoning: false, ContextWindow: 1047576, MaxTokens: 32768, Description: ""},
			{ID: "gpt-4.1-mini", Name: "GPT-4.1 mini", Reasoning: false, ContextWindow: 1047576, MaxTokens: 32768, Description: ""},
			{ID: "gpt-4.1-nano", Name: "GPT-4.1 nano", Reasoning: false, ContextWindow: 1047576, MaxTokens: 32768, Description: ""},
			{ID: "gpt-4o", Name: "GPT-4o", Reasoning: false, ContextWindow: 128000, MaxTokens: 16384, Description: ""},
			{ID: "gpt-4o-2024-05-13", Name: "GPT-4o (2024-05-13)", Reasoning: false, ContextWindow: 128000, MaxTokens: 4096, Description: ""},
			{ID: "gpt-4o-2024-08-06", Name: "GPT-4o (2024-08-06)", Reasoning: false, ContextWindow: 128000, MaxTokens: 16384, Description: ""},
			{ID: "gpt-4o-2024-11-20", Name: "GPT-4o (2024-11-20)", Reasoning: false, ContextWindow: 128000, MaxTokens: 16384, Description: ""},
			{ID: "gpt-4o-mini", Name: "GPT-4o mini", Reasoning: false, ContextWindow: 128000, MaxTokens: 16384, Description: ""},
			{ID: "gpt-5", Name: "GPT-5", Reasoning: true, ContextWindow: 400000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-5-chat-latest", Name: "GPT-5 Chat Latest", Reasoning: false, ContextWindow: 128000, MaxTokens: 16384, Description: ""},
			{ID: "gpt-5-mini", Name: "GPT-5 Mini", Reasoning: true, ContextWindow: 400000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-5-nano", Name: "GPT-5 Nano", Reasoning: true, ContextWindow: 400000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-5-pro", Name: "GPT-5 Pro", Reasoning: true, ContextWindow: 400000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "baseten",
		Name:     "Baseten",
		Type:     "api_key",
		Category: "chips",
		BaseURL:  "https://inference.baseten.co/v1",
		AuthEnv:  "[\"BASETEN_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "deepseek-ai/DeepSeek-V4-Flash-0731", Name: "DeepSeek V4 Flash 0731", Reasoning: true, ContextWindow: 1048576, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4-Pro", Name: "DeepSeek V4 Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4-Pro-0813", Name: "DeepSeek V4 Pro 0813", Reasoning: true, ContextWindow: 1048576, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4.1-Flash", Name: "DeepSeek V4.1 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4.1-Flash-Fast", Name: "DeepSeek V4.1 Flash Fast", Reasoning: true, ContextWindow: 1048576, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "moonshotai/Kimi-K2.5", Name: "Kimi K2.5", Reasoning: true, ContextWindow: 262000, MaxTokens: 262000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "moonshotai/Kimi-K2.6", Name: "Kimi K2.6", Reasoning: true, ContextWindow: 262000, MaxTokens: 262000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "moonshotai/Kimi-K2.7-Code", Name: "Kimi K2.7 Code", Reasoning: true, ContextWindow: 262000, MaxTokens: 262000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "moonshotai/Kimi-K3", Name: "Kimi K3", Reasoning: true, ContextWindow: 1048576, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "nvidia/Nemotron-120B-A12B", Name: "Nemotron Super", Reasoning: true, ContextWindow: 202800, MaxTokens: 202800, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "nvidia/NVIDIA-Nemotron-3-Ultra-550B-A55B", Name: "Nemotron Ultra", Reasoning: true, ContextWindow: 202800, MaxTokens: 202800, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "openai/gpt-oss-120b", Name: "OpenAI GPT 120B", Reasoning: true, ContextWindow: 128072, MaxTokens: 128072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "thinkingmachines/inkling", Name: "Inkling", Reasoning: true, ContextWindow: 1048576, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "thinkingmachines/inkling-small", Name: "Inkling Small", Reasoning: true, ContextWindow: 1048576, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "zai-org/GLM-4.7", Name: "GLM 4.7", Reasoning: true, ContextWindow: 200000, MaxTokens: 200000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "cerebras",
		Name:     "Cerebras",
		Type:     "api_key",
		Category: "chips",
		BaseURL:  "https://api.cerebras.ai/v1",
		AuthEnv:  "[\"CEREBRAS_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "gpt-oss-120b", Name: "GPT OSS 120B", Reasoning: true, ContextWindow: 131072, MaxTokens: 40960, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "qwen-3.8-27b", Name: "Qwen3.8 27B", Reasoning: true, ContextWindow: 131072, MaxTokens: 40960, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "cloudflare-ai-gateway",
		Name:     "Cloudflare AI Gateway",
		Type:     "api_key",
		Category: "gateway",
		BaseURL:  "https://gateway.ai.cloudflare.com/v1/{CLOUDFLARE_ACCOUNT_ID}/{CLOUDFLARE_GATEWAY_ID}/anthropic",
		AuthEnv:  "",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "claude-fable-5", Name: "Claude Fable 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-fable-5-1", Name: "Claude Fable 5.1", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-4-5", Name: "Claude Haiku 4.5 (latest)", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-5-5", Name: "Claude Haiku 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-5", Name: "Claude Opus 4.5 (latest)", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-6", Name: "Claude Opus 4.6", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-7", Name: "Claude Opus 4.7", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-8", Name: "Claude Opus 4.8", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-5", Name: "Claude Opus 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-5-5", Name: "Claude Opus 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-4-5", Name: "Claude Sonnet 4.5 (latest)", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-4-6", Name: "Claude Sonnet 4.6", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-5", Name: "Claude Sonnet 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-5-5", Name: "Claude Sonnet 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "workers-ai/@cf/deepseek-ai/deepseek-v4-flash-0731", Name: "DeepSeek V4 Flash 0731", Reasoning: true, ContextWindow: 1048576, MaxTokens: 1048576, Description: "极速 MoE / 高效响应"},
		},
	},
	{
		ID:       "cloudflare-workers-ai",
		Name:     "Cloudflare Workers AI",
		Type:     "api_key",
		Category: "gateway",
		BaseURL:  "https://api.cloudflare.com/client/v4/accounts/{CLOUDFLARE_ACCOUNT_ID}/ai",
		AuthEnv:  "",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "@cf/cloudflare/clef", Name: "Clef", Reasoning: false, ContextWindow: 65536, MaxTokens: 0, Description: ""},
			{ID: "@cf/cloudflare/clef-flash", Name: "Clef Flash", Reasoning: false, ContextWindow: 65536, MaxTokens: 0, Description: "极速 MoE / 高效响应"},
			{ID: "typesafe/jev", Name: "Jev", Reasoning: false, ContextWindow: 32000, MaxTokens: 0, Description: ""},
			{ID: "@cf/deepseek-ai/deepseek-v4-flash-0731", Name: "DeepSeek V4 Flash 0731", Reasoning: true, ContextWindow: 1048576, MaxTokens: 1048576, Description: "极速 MoE / 高效响应"},
			{ID: "@cf/deepseek-ai/deepseek-v4-pro-0813", Name: "DeepSeek V4 Pro 0813", Reasoning: true, ContextWindow: 1048576, MaxTokens: 1048576, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "@cf/google/gemma-4-26b-a4b-it", Name: "Gemma 4 26B A4B IT", Reasoning: true, ContextWindow: 256000, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "@cf/ibm-granite/granite-4.0-h-micro", Name: "Granite 4.0 H Micro", Reasoning: false, ContextWindow: 131000, MaxTokens: 131000, Description: ""},
			{ID: "@cf/meta/llama-3.3-70b-instruct-fp8-fast", Name: "Llama 3.3 70B Instruct fp8 Fast", Reasoning: false, ContextWindow: 24000, MaxTokens: 24000, Description: ""},
			{ID: "@cf/meta/llama-4-scout-17b-16e-instruct", Name: "Llama 4 Scout 17B 16E Instruct", Reasoning: false, ContextWindow: 131000, MaxTokens: 16384, Description: ""},
			{ID: "@cf/mistralai/mistral-small-3.1-24b-instruct", Name: "Mistral Small 3.1 24B Instruct", Reasoning: false, ContextWindow: 128000, MaxTokens: 128000, Description: ""},
			{ID: "@cf/moonshotai/kimi-k2.6", Name: "Kimi K2.6", Reasoning: true, ContextWindow: 262144, MaxTokens: 256000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "@cf/moonshotai/kimi-k2.7-code", Name: "Kimi K2.7 Code", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "专业代码与工程开发模型"},
			{ID: "@cf/nvidia/nemotron-3-120b-a12b", Name: "Nemotron 3 Super 120B", Reasoning: true, ContextWindow: 256000, MaxTokens: 256000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "@cf/openai/gpt-oss-120b", Name: "GPT OSS 120B", Reasoning: true, ContextWindow: 128000, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "@cf/openai/gpt-oss-20b", Name: "GPT OSS 20B", Reasoning: true, ContextWindow: 128000, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "deepseek",
		Name:     "DeepSeek",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://api.deepseek.com",
		AuthEnv:  "[\"DEEPSEEK_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "deepseek-flash", Name: "DeepSeek V4.1 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "deepseek-v4-pro", Name: "DeepSeek V4 Pro", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "fireworks",
		Name:     "Fireworks",
		Type:     "api_key",
		Category: "chips",
		BaseURL:  "https://api.fireworks.ai/inference",
		AuthEnv:  "[\"FIREWORKS_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "accounts/fireworks/models/deepseek-v4p1-flash", Name: "DeepSeek V4.1 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "accounts/fireworks/models/ember-1", Name: "Ember-1", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/models/gpt-oss-120b", Name: "GPT OSS 120B", Reasoning: true, ContextWindow: 131072, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/models/inkling", Name: "Inkling", Reasoning: true, ContextWindow: 1048576, MaxTokens: 1048576, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/models/minimax-m3", Name: "MiniMax-M3", Reasoning: true, ContextWindow: 512000, MaxTokens: 512000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/models/nemotron-3-ultra-nvfp4", Name: "Nemotron 3 Ultra 550B A55B", Reasoning: true, ContextWindow: 262144, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/models/nemotron-lightning-3p5-30b-a3b", Name: "Nemotron 3.5 Lightning 30B A3B", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/models/qwen3p8-2p4t-a95b", Name: "Qwen3.8 2.4T A95B", Reasoning: true, ContextWindow: 262144, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/models/qwen3p8-max", Name: "Qwen3.8 Max", Reasoning: true, ContextWindow: 262144, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/routers/deepseek-flash-latest", Name: "DeepSeek Flash Latest", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "accounts/fireworks/routers/kimi-fast-latest", Name: "Kimi Fast Latest", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/routers/kimi-latest", Name: "Kimi Latest", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/routers/minimax-latest", Name: "MiniMax Latest", Reasoning: true, ContextWindow: 512000, MaxTokens: 512000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/routers/qwen-max-latest", Name: "Qwen Max Latest (Qwen3.8 Max)", Reasoning: true, ContextWindow: 262144, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "accounts/fireworks/models/glm-5p3", Name: "GLM 5.3", Reasoning: true, ContextWindow: 1048573, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "github-copilot",
		Name:     "GitHub Copilot",
		Type:     "account",
		Category: "global",
		BaseURL:  "https://api.individual.githubcopilot.com",
		AuthEnv:  "[\"COPILOT_GITHUB_TOKEN\"]",
		AuthType: "copilot-token",
		Models: []RegistryModel{
			{ID: "claude-fable-5", Name: "Claude Fable 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-fable-5.1", Name: "Claude Fable 5.1", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-4.5", Name: "Claude Haiku 4.5 (latest)", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-5.5", Name: "Claude Haiku 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4.7", Name: "Claude Opus 4.7", Reasoning: true, ContextWindow: 1000000, MaxTokens: 32000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4.8", Name: "Claude Opus 4.8", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-5", Name: "Claude Opus 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-5.5", Name: "Claude Opus 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-4.6", Name: "Claude Sonnet 4.6", Reasoning: true, ContextWindow: 1000000, MaxTokens: 32000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-5", Name: "Claude Sonnet 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-5.5", Name: "Claude Sonnet 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gemini-3.5-flash", Name: "Gemini 3.5 Flash", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.6-flash", Name: "Gemini 3.6 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.7-flash", Name: "Gemini 3.7 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.8-flash", Name: "Gemini 3.8 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "极速 MoE / 高效响应"},
		},
	},
	{
		ID:       "google",
		Name:     "Google",
		Type:     "api_key",
		Category: "global",
		BaseURL:  "https://generativelanguage.googleapis.com/v1beta",
		AuthEnv:  "[\"GEMINI_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "deep-research-max-preview-04-2026", Name: "Deep Research Max Preview (Apr-21-2026)", Reasoning: true, ContextWindow: 131072, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deep-research-preview-04-2026", Name: "Deep Research Preview (Apr-21-2026)", Reasoning: true, ContextWindow: 131072, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gemini-2.5-computer-use-preview-10-2025", Name: "Gemini 2.5 Computer Use Preview 10-2025", Reasoning: true, ContextWindow: 128000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-2.5-flash-lite", Name: "Gemini 2.5 Flash-Lite", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gemini-3-flash-preview", Name: "Gemini 3 Flash Preview", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.1-flash-lite", Name: "Gemini 3.1 Flash Lite", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.1-flash-lite-image", Name: "Nano Banana 2 Lite", Reasoning: true, ContextWindow: 65536, MaxTokens: 4096, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.1-flash-lite-preview", Name: "Gemini 3.1 Flash Lite Preview", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.1-flash-live-preview", Name: "Gemini 3.1 Flash Live Preview", Reasoning: true, ContextWindow: 131072, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.1-pro-preview", Name: "Gemini 3.1 Pro Preview", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gemini-3.1-pro-preview-customtools", Name: "Gemini 3.1 Pro Preview Custom Tools", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gemini-3.5-flash", Name: "Gemini 3.5 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.5-flash-lite", Name: "Gemini 3.5 Flash Lite", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
		},
	},
	{
		ID:       "google-vertex",
		Name:     "Google Cloud credentials",
		Type:     "api_key",
		Category: "global",
		BaseURL:  "https://{location}-aiplatform.googleapis.com",
		AuthEnv:  "",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "gemini-2.5-flash", Name: "Gemini 2.5 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-2.5-flash-lite", Name: "Gemini 2.5 Flash-Lite", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-2.5-pro", Name: "Gemini 2.5 Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gemini-3-flash-preview", Name: "Gemini 3 Flash Preview", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.1-flash-lite", Name: "Gemini 3.1 Flash Lite", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.1-flash-lite-image", Name: "Nano Banana 2 Lite", Reasoning: true, ContextWindow: 65536, MaxTokens: 4096, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.1-pro-preview", Name: "Gemini 3.1 Pro Preview", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gemini-3.1-pro-preview-customtools", Name: "Gemini 3.1 Pro Preview Custom Tools", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gemini-3.5-flash", Name: "Gemini 3.5 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.5-flash-lite", Name: "Gemini 3.5 Flash Lite", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.6-flash", Name: "Gemini 3.6 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.7-flash", Name: "Gemini 3.7 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "gemini-3.8-flash", Name: "Gemini 3.8 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
		},
	},
	{
		ID:       "groq",
		Name:     "Groq",
		Type:     "api_key",
		Category: "chips",
		BaseURL:  "https://api.groq.com/openai/v1",
		AuthEnv:  "[\"GROQ_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "llama-3.1-8b-instant", Name: "Llama 3.1 8B", Reasoning: false, ContextWindow: 131072, MaxTokens: 131072, Description: ""},
			{ID: "llama-3.3-70b-versatile", Name: "Llama 3.3 70B", Reasoning: false, ContextWindow: 131072, MaxTokens: 32768, Description: ""},
			{ID: "openai/gpt-oss-120b", Name: "GPT OSS 120B", Reasoning: true, ContextWindow: 131072, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "openai/gpt-oss-20b", Name: "GPT OSS 20B", Reasoning: true, ContextWindow: 131072, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "openai/gpt-oss-safeguard-20b", Name: "Safety GPT OSS 20B", Reasoning: true, ContextWindow: 131072, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "qwen/qwen3.6-27b", Name: "Qwen3.6 27B", Reasoning: true, ContextWindow: 131072, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "qwen/qwen3.8-27b", Name: "Qwen3.8 27B", Reasoning: true, ContextWindow: 131042, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "huggingface",
		Name:     "Hugging Face",
		Type:     "api_key",
		Category: "global",
		BaseURL:  "https://router.huggingface.co/v1",
		AuthEnv:  "[\"HF_TOKEN\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "deepseek-ai/DeepSeek-R1", Name: "DeepSeek-R1", Reasoning: true, ContextWindow: 64000, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-R1-0528", Name: "DeepSeek-R1-0528", Reasoning: true, ContextWindow: 163840, MaxTokens: 163840, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V3", Name: "DeepSeek-V3", Reasoning: false, ContextWindow: 64000, MaxTokens: 8192, Description: ""},
			{ID: "deepseek-ai/DeepSeek-V3-0324", Name: "DeepSeek V3 0324", Reasoning: false, ContextWindow: 163840, MaxTokens: 163840, Description: ""},
			{ID: "deepseek-ai/DeepSeek-V3.1", Name: "DeepSeek-V3.1", Reasoning: true, ContextWindow: 131072, MaxTokens: 8192, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V3.2", Name: "DeepSeek-V3.2", Reasoning: true, ContextWindow: 163840, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4-Flash", Name: "DeepSeek V4 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4-Flash-0731", Name: "DeepSeek V4 Flash 0731", Reasoning: true, ContextWindow: 1048576, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4-Flash-Vision-Exp", Name: "DeepSeek V4 Flash Vision Exp", Reasoning: true, ContextWindow: 1048576, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4-Pro", Name: "DeepSeek V4 Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 393216, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4-Pro-0813", Name: "DeepSeek V4 Pro 0813", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4.1-Flash", Name: "DeepSeek V4.1 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "google/gemma-3-12b-it", Name: "Gemma 3 12B IT", Reasoning: false, ContextWindow: 131072, MaxTokens: 131072, Description: ""},
			{ID: "google/gemma-3-27b-it", Name: "Gemma 3 27B IT", Reasoning: false, ContextWindow: 131072, MaxTokens: 131072, Description: ""},
			{ID: "google/gemma-3-4b-it", Name: "Gemma 3 4B IT", Reasoning: false, ContextWindow: 131072, MaxTokens: 131072, Description: ""},
		},
	},
	{
		ID:       "kimi-coding",
		Name:     "Kimi For Coding",
		Type:     "account",
		Category: "global",
		BaseURL:  "https://api.kimi.com/coding",
		AuthEnv:  "[\"KIMI_API_KEY\"]",
		AuthType: "x-api-key",
		Models: []RegistryModel{
			{ID: "k3", Name: "Kimi K3", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "k3-256k", Name: "Kimi K3-256K", Reasoning: true, ContextWindow: 262144, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-for-coding", Name: "kimi-for-coding", Reasoning: true, ContextWindow: 1048576, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-for-coding-highspeed", Name: "Kimi For Coding HighSpeed", Reasoning: true, ContextWindow: 262144, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "llama.cpp",
		Name:     "llama.cpp",
		Type:     "api_key",
		Category: "local",
		BaseURL:  "http://localhost:8080/v1",
		AuthEnv:  "None (Local server)",
		AuthType: "none",
		Models: []RegistryModel{
			{ID: "local-model", Name: "llama.cpp 本地自载入模型", Reasoning: true, ContextWindow: 128000, MaxTokens: 0, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "meta",
		Name:     "Meta",
		Type:     "account",
		Category: "global",
		BaseURL:  "https://api.meta.ai/v1",
		AuthEnv:  "[\"META_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "muse-spark-1.1", Name: "Muse Spark 1.1", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "muse-spark-1.2", Name: "Muse Spark 1.2", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "muse-spark-1.2-contributor", Name: "Muse Spark 1.2 Contributor", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "muse-spark-1.3", Name: "Muse Spark 1.3", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "muse-spark-1.3-contributor", Name: "Muse Spark 1.3 Contributor", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "muse-wonder-1.0", Name: "Muse Wonder", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "minimax",
		Name:     "MiniMax",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://api.minimax.io/anthropic",
		AuthEnv:  "[\"MINIMAX_API_KEY\"]",
		AuthType: "x-api-key",
		Models: []RegistryModel{
			{ID: "MiniMax-M2.7", Name: "MiniMax-M2.7", Reasoning: true, ContextWindow: 204800, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "MiniMax-M2.7-highspeed", Name: "MiniMax-M2.7-highspeed", Reasoning: true, ContextWindow: 204800, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "MiniMax-M3", Name: "MiniMax-M3", Reasoning: true, ContextWindow: 1000000, MaxTokens: 512000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "minimax-cn",
		Name:     "MiniMax CN",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://api.minimaxi.com/anthropic",
		AuthEnv:  "[\"MINIMAX_CN_API_KEY\"]",
		AuthType: "x-api-key",
		Models: []RegistryModel{
			{ID: "MiniMax-M2.7", Name: "MiniMax-M2.7", Reasoning: true, ContextWindow: 204800, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "MiniMax-M2.7-highspeed", Name: "MiniMax-M2.7-highspeed", Reasoning: true, ContextWindow: 204800, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "MiniMax-M3", Name: "MiniMax-M3", Reasoning: true, ContextWindow: 1000000, MaxTokens: 512000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "mistral",
		Name:     "Mistral",
		Type:     "api_key",
		Category: "global",
		BaseURL:  "https://api.mistral.ai",
		AuthEnv:  "[\"MISTRAL_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "codestral-2508", Name: "Codestral 25.08", Reasoning: false, ContextWindow: 256000, MaxTokens: 8192, Description: "专业代码与工程开发模型"},
			{ID: "codestral-latest", Name: "Codestral (latest)", Reasoning: false, ContextWindow: 256000, MaxTokens: 4096, Description: "专业代码与工程开发模型"},
			{ID: "devstral-2512", Name: "Devstral 2", Reasoning: false, ContextWindow: 262144, MaxTokens: 262144, Description: ""},
			{ID: "devstral-latest", Name: "Devstral 2", Reasoning: false, ContextWindow: 262144, MaxTokens: 262144, Description: ""},
			{ID: "devstral-medium-2507", Name: "Devstral Medium", Reasoning: false, ContextWindow: 128000, MaxTokens: 128000, Description: ""},
			{ID: "devstral-medium-latest", Name: "Devstral 2 (latest)", Reasoning: false, ContextWindow: 262144, MaxTokens: 262144, Description: ""},
			{ID: "devstral-small-2505", Name: "Devstral Small 2505", Reasoning: false, ContextWindow: 128000, MaxTokens: 128000, Description: ""},
			{ID: "devstral-small-2507", Name: "Devstral Small", Reasoning: false, ContextWindow: 128000, MaxTokens: 128000, Description: ""},
			{ID: "glm-5-2", Name: "GLM-5.2", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "labs-devstral-small-2512", Name: "Devstral Small 2", Reasoning: false, ContextWindow: 256000, MaxTokens: 256000, Description: ""},
			{ID: "labs-leanstral-1-5-1", Name: "Leanstral 1.5", Reasoning: true, ContextWindow: 262144, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "magistral-medium-latest", Name: "Magistral Medium (latest)", Reasoning: true, ContextWindow: 262144, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "ministral-14b-2512", Name: "Ministral 3 14B", Reasoning: false, ContextWindow: 262144, MaxTokens: 262144, Description: ""},
			{ID: "ministral-3b-2512", Name: "Ministral 3 3B", Reasoning: false, ContextWindow: 131072, MaxTokens: 262144, Description: ""},
			{ID: "ministral-3b-latest", Name: "Ministral 3B (latest)", Reasoning: false, ContextWindow: 131072, MaxTokens: 128000, Description: ""},
		},
	},
	{
		ID:       "moonshotai",
		Name:     "Moonshot AI",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://api.moonshot.ai/v1",
		AuthEnv:  "[\"MOONSHOT_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "kimi-k2.6", Name: "Kimi K2.6", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-k2.7-code", Name: "Kimi K2.7 Code", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "专业代码与工程开发模型"},
			{ID: "kimi-k2.7-code-highspeed", Name: "Kimi K2.7 Code HighSpeed", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "专业代码与工程开发模型"},
			{ID: "kimi-k3", Name: "Kimi K3", Reasoning: true, ContextWindow: 1048576, MaxTokens: 1048576, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "moonshotai-cn",
		Name:     "Moonshot AI CN",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://api.moonshot.cn/v1",
		AuthEnv:  "[\"MOONSHOT_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "kimi-k2.6", Name: "Kimi K2.6", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-k2.7-code", Name: "Kimi K2.7 Code", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "专业代码与工程开发模型"},
			{ID: "kimi-k2.7-code-highspeed", Name: "Kimi K2.7 Code HighSpeed", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "专业代码与工程开发模型"},
			{ID: "kimi-k3", Name: "Kimi K3", Reasoning: true, ContextWindow: 1048576, MaxTokens: 1048576, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "nvidia",
		Name:     "NVIDIA",
		Type:     "api_key",
		Category: "chips",
		BaseURL:  "https://integrate.api.nvidia.com/v1",
		AuthEnv:  "[\"NVIDIA_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "deepseek-ai/deepseek-v4.1-flash", Name: "DeepSeek V4.1 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "google/diffusiongemma-26b-a4b-it", Name: "DiffusionGemma 26B A4B IT", Reasoning: true, ContextWindow: 250000, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "google/gemma-3-12b-it", Name: "Gemma 3 12B IT", Reasoning: false, ContextWindow: 131072, MaxTokens: 16384, Description: ""},
			{ID: "google/gemma-3-4b-it", Name: "Gemma 3 4B IT", Reasoning: false, ContextWindow: 131072, MaxTokens: 16384, Description: ""},
			{ID: "meta/llama-3.2-11b-vision-instruct", Name: "Llama 3.2 11b Vision Instruct", Reasoning: false, ContextWindow: 128000, MaxTokens: 4096, Description: ""},
			{ID: "meta/llama-3.2-90b-vision-instruct", Name: "Llama-3.2-90B-Vision-Instruct", Reasoning: false, ContextWindow: 128000, MaxTokens: 8192, Description: ""},
			{ID: "meta/muse-glimmer-30b", Name: "Muse Glimmer 30B", Reasoning: true, ContextWindow: 131072, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mistralai/mistral-7b-instruct-v0.3", Name: "Mistral-7B-Instruct-v0.3", Reasoning: false, ContextWindow: 65536, MaxTokens: 65536, Description: ""},
			{ID: "moonshotai/kimi-k2.6", Name: "Kimi K2.6", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "moonshotai/kimi-k3", Name: "Kimi K3", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "nvidia/cosmos-reason2-8b", Name: "Cosmos Reason2 8B", Reasoning: true, ContextWindow: 131072, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "nvidia/llama-3.1-nemotron-70b-instruct", Name: "Llama 3.1 Nemotron 70B Instruct", Reasoning: false, ContextWindow: 128000, MaxTokens: 8192, Description: ""},
			{ID: "nvidia/llama-3.1-nemotron-ultra-253b-v1", Name: "Llama 3.1 Nemotron Ultra 253B", Reasoning: true, ContextWindow: 128000, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "nvidia/nemotron-3-nano-omni-30b-a3b-reasoning", Name: "Nemotron 3 Nano Omni", Reasoning: true, ContextWindow: 256000, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "nvidia/nemotron-3-super-120b-a12b", Name: "Nemotron 3 Super", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "openai",
		Name:     "OpenAI",
		Type:     "account",
		Category: "global",
		BaseURL:  "https://api.openai.com/v1",
		AuthEnv:  "[\"OPENAI_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "gpt-6-luna", Name: "GPT-6 Luna", Reasoning: false, ContextWindow: 922000, MaxTokens: 0, Description: ""},
			{ID: "gpt-4", Name: "GPT-4", Reasoning: false, ContextWindow: 8192, MaxTokens: 8192, Description: ""},
			{ID: "gpt-4-turbo", Name: "GPT-4 Turbo", Reasoning: false, ContextWindow: 128000, MaxTokens: 4096, Description: ""},
			{ID: "gpt-4.1", Name: "GPT-4.1", Reasoning: false, ContextWindow: 1047576, MaxTokens: 32768, Description: ""},
			{ID: "gpt-4.1-mini", Name: "GPT-4.1 mini", Reasoning: false, ContextWindow: 1047576, MaxTokens: 32768, Description: ""},
			{ID: "gpt-4.1-nano", Name: "GPT-4.1 nano", Reasoning: false, ContextWindow: 1047576, MaxTokens: 32768, Description: ""},
			{ID: "gpt-4o", Name: "GPT-4o", Reasoning: false, ContextWindow: 128000, MaxTokens: 16384, Description: ""},
			{ID: "gpt-4o-2024-05-13", Name: "GPT-4o (2024-05-13)", Reasoning: false, ContextWindow: 128000, MaxTokens: 4096, Description: ""},
			{ID: "gpt-4o-2024-08-06", Name: "GPT-4o (2024-08-06)", Reasoning: false, ContextWindow: 128000, MaxTokens: 16384, Description: ""},
			{ID: "gpt-4o-2024-11-20", Name: "GPT-4o (2024-11-20)", Reasoning: false, ContextWindow: 128000, MaxTokens: 16384, Description: ""},
			{ID: "gpt-4o-mini", Name: "GPT-4o mini", Reasoning: false, ContextWindow: 128000, MaxTokens: 16384, Description: ""},
			{ID: "gpt-5", Name: "GPT-5", Reasoning: true, ContextWindow: 400000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-5-chat-latest", Name: "GPT-5 Chat Latest", Reasoning: false, ContextWindow: 128000, MaxTokens: 16384, Description: ""},
			{ID: "gpt-5-mini", Name: "GPT-5 Mini", Reasoning: true, ContextWindow: 400000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-5-nano", Name: "GPT-5 Nano", Reasoning: true, ContextWindow: 400000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "openai-codex",
		Name:     "OpenAI Codex (legacy)",
		Type:     "account",
		Category: "global",
		BaseURL:  "https://chatgpt.com/backend-api",
		AuthEnv:  "",
		AuthType: "session-token",
		Models: []RegistryModel{
			{ID: "gpt-5.3-codex-spark", Name: "GPT-5.3 Codex Spark", Reasoning: true, ContextWindow: 128000, MaxTokens: 128000, Description: "专业代码与工程开发模型"},
			{ID: "gpt-5.5", Name: "GPT-5.5", Reasoning: true, ContextWindow: 272000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-5.6-luna", Name: "GPT-5.6 Luna", Reasoning: true, ContextWindow: 272000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-5.6-sol", Name: "GPT-5.6 Sol", Reasoning: true, ContextWindow: 272000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-5.6-terra", Name: "GPT-5.6 Terra", Reasoning: true, ContextWindow: 272000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-6-astra", Name: "GPT-6 Astra", Reasoning: true, ContextWindow: 272000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-6-luna", Name: "GPT-6 Luna", Reasoning: true, ContextWindow: 272000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-6-sol", Name: "GPT-6 Sol", Reasoning: true, ContextWindow: 272000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "gpt-6.1-sol", Name: "GPT-6.1 Sol", Reasoning: true, ContextWindow: 272000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "opencode",
		Name:     "OpenCode Zen",
		Type:     "api_key",
		Category: "global",
		BaseURL:  "https://opencode.ai/zen",
		AuthEnv:  "[\"OPENCODE_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "claude-fable-5", Name: "Claude Fable 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-fable-5-1", Name: "Claude Fable 5.1", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-4-5", Name: "Claude Haiku 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-5-5", Name: "Claude Haiku 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-5", Name: "Claude Opus 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-6", Name: "Claude Opus 4.6", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-7", Name: "Claude Opus 4.7", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-8", Name: "Claude Opus 4.8", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-5", Name: "Claude Opus 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-5-5", Name: "Claude Opus 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-4", Name: "Claude Sonnet 4", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-4-5", Name: "Claude Sonnet 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-4-6", Name: "Claude Sonnet 4.6", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-5", Name: "Claude Sonnet 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-5-5", Name: "Claude Sonnet 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "opencode-go",
		Name:     "OpenCode Go",
		Type:     "api_key",
		Category: "global",
		BaseURL:  "https://opencode.ai/zen/go",
		AuthEnv:  "[\"OPENCODE_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "claude-haiku-5-5", Name: "Claude Haiku 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "minimax-m3", Name: "MiniMax-M3", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "qwen3.7-plus", Name: "Qwen3.7 Plus", Reasoning: true, ContextWindow: 1000000, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "qwen3.8-flash", Name: "Qwen3.8 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "极速 MoE / 高效响应"},
			{ID: "qwen3.8-max", Name: "Qwen3.8 Max", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-v4-flash", Name: "DeepSeek V4 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "deepseek-v4-flash-vision-exp", Name: "DeepSeek V4 Flash Vision Exp", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "deepseek-v4-pro", Name: "DeepSeek V4 Pro (New)", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-v4.1-flash", Name: "DeepSeek V4.1 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "glm-5.2", Name: "GLM-5.2", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.3", Name: "GLM-5.3", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.3-flash", Name: "GLM-5.3-Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "极速 MoE / 高效响应"},
			{ID: "hy3", Name: "Hy3", Reasoning: true, ContextWindow: 256000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "hy4-preview", Name: "Hy4 preview", Reasoning: true, ContextWindow: 1024000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-k2.7-code", Name: "Kimi K2.7 Code", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "专业代码与工程开发模型"},
		},
	},
	{
		ID:       "openrouter",
		Name:     "OpenRouter",
		Type:     "account",
		Category: "gateway",
		BaseURL:  "https://openrouter.ai/api",
		AuthEnv:  "[\"OPENROUTER_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "anthropic/claude-fable-5", Name: "Anthropic: Claude Fable 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-fable-5.1", Name: "Anthropic: Claude Fable 5.1", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-haiku-4.5", Name: "Anthropic: Claude Haiku 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-haiku-5.5", Name: "Anthropic: Claude Haiku 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-opus-4.1", Name: "Anthropic: Claude Opus 4.1", Reasoning: true, ContextWindow: 200000, MaxTokens: 32000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-opus-4.5", Name: "Anthropic: Claude Opus 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-opus-4.6", Name: "Anthropic: Claude Opus 4.6", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-opus-4.7", Name: "Anthropic: Claude Opus 4.7", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-opus-4.8", Name: "Anthropic: Claude Opus 4.8", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-opus-5", Name: "Anthropic: Claude Opus 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-opus-5.5", Name: "Anthropic: Claude Opus 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-sonnet-4", Name: "Anthropic: Claude Sonnet 4", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-sonnet-4.5", Name: "Anthropic: Claude Sonnet 4.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-sonnet-4.6", Name: "Anthropic: Claude Sonnet 4.6", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "anthropic/claude-sonnet-5", Name: "Anthropic: Claude Sonnet 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "qwen-token-plan",
		Name:     "Qwen Token Plan",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
		AuthEnv:  "[\"QWEN_TOKEN_PLAN_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "deepseek-v3.2", Name: "DeepSeek V3.2", Reasoning: true, ContextWindow: 131072, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-v4-flash", Name: "DeepSeek V4 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "deepseek-v4-flash-0731", Name: "DeepSeek V4 Flash 0731", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "deepseek-v4-pro", Name: "DeepSeek V4 Pro", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-v4-pro-0813", Name: "DeepSeek V4 Pro 0813", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-v4.1-flash", Name: "DeepSeek V4.1 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "glm-5", Name: "GLM-5", Reasoning: true, ContextWindow: 202752, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.1", Name: "GLM-5.1", Reasoning: true, ContextWindow: 202752, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.2", Name: "GLM-5.2", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.3", Name: "GLM-5.3", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-k2.5", Name: "Kimi K2.5", Reasoning: true, ContextWindow: 262144, MaxTokens: 98304, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-k2.6", Name: "Kimi K2.6", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-k2.7-code", Name: "Kimi K2.7 Code", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "专业代码与工程开发模型"},
			{ID: "MiniMax-M2.5", Name: "MiniMax-M2.5", Reasoning: true, ContextWindow: 196608, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "qwen3.6-flash", Name: "Qwen3.6 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
		},
	},
	{
		ID:       "qwen-token-plan-cn",
		Name:     "Qwen Token Plan CN",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://token-plan.cn-beijing.maas.aliyuncs.com/compatible-mode/v1",
		AuthEnv:  "[\"QWEN_TOKEN_PLAN_CN_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "deepseek-v3.2", Name: "DeepSeek V3.2", Reasoning: true, ContextWindow: 131072, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-v4-flash", Name: "DeepSeek V4 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "deepseek-v4-flash-0731", Name: "DeepSeek V4 Flash 0731", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "deepseek-v4-pro", Name: "DeepSeek V4 Pro", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-v4-pro-0813", Name: "DeepSeek V4 Pro 0813", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-v4.1-flash", Name: "DeepSeek V4.1 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "glm-5", Name: "GLM-5", Reasoning: true, ContextWindow: 202752, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.1", Name: "GLM-5.1", Reasoning: true, ContextWindow: 202752, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.2", Name: "GLM-5.2", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.3", Name: "GLM-5.3", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-k2.5", Name: "Kimi K2.5", Reasoning: true, ContextWindow: 262144, MaxTokens: 98304, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-k2.6", Name: "Kimi K2.6", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "kimi-k2.7-code", Name: "Kimi K2.7 Code", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "专业代码与工程开发模型"},
			{ID: "MiniMax-M2.5", Name: "MiniMax-M2.5", Reasoning: true, ContextWindow: 196608, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "qwen3.6-flash", Name: "Qwen3.6 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
		},
	},
	{
		ID:       "qwen-token-plan-individual",
		Name:     "Qwen Token Plan Individual",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://token-plan.ap-southeast-1.maas.aliyuncs.com/compatible-mode/v1",
		AuthEnv:  "[\"QWEN_TOKEN_PLAN_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "deepseek-v4-flash-0731", Name: "DeepSeek V4 Flash 0731", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "极速 MoE / 高效响应"},
			{ID: "deepseek-v4-pro", Name: "DeepSeek V4 Pro", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-v4-pro-0813", Name: "DeepSeek V4 Pro 0813", Reasoning: true, ContextWindow: 1000000, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.2", Name: "GLM-5.2", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "qwen3.6-flash", Name: "Qwen3.6 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 65536, Description: "极速 MoE / 高效响应"},
			{ID: "qwen3.7-max", Name: "Qwen3.7 Max", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "qwen3.7-plus", Name: "Qwen3.7 Plus", Reasoning: true, ContextWindow: 1000000, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "qwen3.8-flash", Name: "Qwen3.8 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "极速 MoE / 高效响应"},
			{ID: "qwen3.8-max", Name: "Qwen3.8 Max", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "radius",
		Name:     "radius",
		Type:     "account",
		Category: "gateway",
		BaseURL:  "https://radius.pi.dev/v1",
		AuthEnv:  "[\"RADIUS_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "balanced", Name: "Balanced", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "cheap", Name: "Cheap", Reasoning: true, ContextWindow: 1000000, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-fable-5", Name: "Claude Fable 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-fable-5-1", Name: "Claude Fable 5.1", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-4-5", Name: "Claude Haiku 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-haiku-5-5", Name: "Claude Haiku 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-5", Name: "Claude Opus 4.5", Reasoning: true, ContextWindow: 200000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-4-8", Name: "Claude Opus 4.8", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-5", Name: "Claude Opus 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-opus-5-5", Name: "Claude Opus 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-4-5", Name: "Claude Sonnet 4.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-5", Name: "Claude Sonnet 5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "claude-sonnet-5-5", Name: "Claude Sonnet 5.5", Reasoning: true, ContextWindow: 1000000, MaxTokens: 128000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-v4.1-flash", Name: "DeepSeek V4.1 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 32768, Description: "极速 MoE / 高效响应"},
			{ID: "glm-5.3", Name: "GLM-5.3", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "together",
		Name:     "Together",
		Type:     "api_key",
		Category: "chips",
		BaseURL:  "https://api.together.ai/v1",
		AuthEnv:  "[\"TOGETHER_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "deepseek-ai/DeepSeek-V4-Flash-0731", Name: "DeepSeek V4 Flash 0731", Reasoning: true, ContextWindow: 1048576, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4-Pro-0813", Name: "DeepSeek V4 Pro 0813", Reasoning: true, ContextWindow: 1048576, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "deepseek-ai/DeepSeek-V4.1-Flash", Name: "DeepSeek V4.1 Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 384000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "meta-llama/Llama-3.3-70B-Instruct-Turbo", Name: "Llama 3.3 70B", Reasoning: false, ContextWindow: 131072, MaxTokens: 131072, Description: ""},
			{ID: "MiniMaxAI/MiniMax-M2.7", Name: "MiniMax-M2.7", Reasoning: true, ContextWindow: 196608, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "MiniMaxAI/MiniMax-M3", Name: "MiniMax-M3", Reasoning: true, ContextWindow: 524288, MaxTokens: 250000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "moonshotai/Kimi-K3", Name: "Kimi K3", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "nvidia/nemotron-3-ultra-550b-a55b", Name: "Nemotron 3 Ultra 550B A55B", Reasoning: true, ContextWindow: 512300, MaxTokens: 512300, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "openai/gpt-oss-120b", Name: "GPT OSS 120B", Reasoning: true, ContextWindow: 131072, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "Qwen/Qwen2.5-7B-Instruct-Turbo", Name: "Qwen 2.5 7B Instruct Turbo", Reasoning: false, ContextWindow: 32768, MaxTokens: 32768, Description: ""},
			{ID: "Qwen/Qwen3.5-9B", Name: "Qwen3.5 9B", Reasoning: true, ContextWindow: 262144, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "Qwen/Qwen3.6-Plus", Name: "Qwen3.6 Plus", Reasoning: true, ContextWindow: 1000000, MaxTokens: 500000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "Qwen/Qwen3.7-Max", Name: "Qwen3.7 Max", Reasoning: false, ContextWindow: 1000000, MaxTokens: 500000, Description: ""},
			{ID: "thinkingmachines/Inkling", Name: "Inkling", Reasoning: true, ContextWindow: 524288, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "zai-org/GLM-5.2", Name: "GLM-5.2", Reasoning: true, ContextWindow: 1048575, MaxTokens: 164000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "typesafe",
		Name:     "TypeSafe",
		Type:     "api_key",
		Category: "global",
		BaseURL:  "https://api.typesafe.ai/v1/",
		AuthEnv:  "[\"TYPESAFE_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "jev-latest", Name: "Jev", Reasoning: false, ContextWindow: 64000, MaxTokens: 0, Description: ""},
		},
	},
	{
		ID:       "vercel-ai-gateway",
		Name:     "Vercel AI Gateway",
		Type:     "api_key",
		Category: "gateway",
		BaseURL:  "https://ai-gateway.vercel.sh",
		AuthEnv:  "[\"AI_GATEWAY_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "alibaba/qwen-3-14b", Name: "Qwen3-14B", Reasoning: true, ContextWindow: 40960, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "alibaba/qwen-3-235b", Name: "Qwen3 235B A22B", Reasoning: false, ContextWindow: 262144, MaxTokens: 16384, Description: ""},
			{ID: "alibaba/qwen-3-30b", Name: "Qwen3-30B-A3B", Reasoning: true, ContextWindow: 40960, MaxTokens: 16384, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "alibaba/qwen-3-32b", Name: "Qwen 3 32B", Reasoning: true, ContextWindow: 128000, MaxTokens: 8192, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "alibaba/qwen3-coder", Name: "Qwen3 Coder 480B A35B Instruct", Reasoning: false, ContextWindow: 262144, MaxTokens: 32768, Description: "专业代码与工程开发模型"},
			{ID: "alibaba/qwen3-coder-30b-a3b", Name: "Qwen 3 Coder 30B A3B Instruct", Reasoning: false, ContextWindow: 262144, MaxTokens: 8192, Description: "专业代码与工程开发模型"},
			{ID: "alibaba/qwen3-coder-next", Name: "Qwen3 Coder Next", Reasoning: false, ContextWindow: 256000, MaxTokens: 256000, Description: "专业代码与工程开发模型"},
			{ID: "alibaba/qwen3-max-thinking", Name: "Qwen 3 Max Thinking", Reasoning: true, ContextWindow: 256000, MaxTokens: 65536, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "alibaba/qwen3-next-80b-a3b-instruct", Name: "Qwen3 Next 80B A3B Instruct", Reasoning: false, ContextWindow: 262114, MaxTokens: 262114, Description: ""},
			{ID: "alibaba/qwen3-next-80b-a3b-thinking", Name: "Qwen3 Next 80B A3B Thinking", Reasoning: true, ContextWindow: 262144, MaxTokens: 262144, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "alibaba/qwen3-vl-instruct", Name: "Qwen3 VL 235B A22B Instruct", Reasoning: false, ContextWindow: 262144, MaxTokens: 262144, Description: ""},
			{ID: "alibaba/qwen3.5-flash", Name: "Qwen 3.5 Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "极速 MoE / 高效响应"},
			{ID: "alibaba/qwen3.5-plus", Name: "Qwen 3.5 Plus", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "alibaba/qwen3.6-27b", Name: "Qwen 3.6 27B", Reasoning: true, ContextWindow: 256000, MaxTokens: 256000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "alibaba/qwen3.6-plus", Name: "Qwen 3.6 Plus", Reasoning: true, ContextWindow: 1000000, MaxTokens: 64000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "xai",
		Name:     "xAI",
		Type:     "account",
		Category: "global",
		BaseURL:  "https://api.x.ai/v1",
		AuthEnv:  "[\"XAI_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "grok-4.3", Name: "Grok 4.3", Reasoning: true, ContextWindow: 1000000, MaxTokens: 30000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "grok-4.5", Name: "Grok 4.5", Reasoning: true, ContextWindow: 500000, MaxTokens: 500000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "grok-4.6", Name: "Grok 4.6", Reasoning: true, ContextWindow: 500000, MaxTokens: 500000, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "grok-4.7", Name: "Grok 4.7", Reasoning: true, ContextWindow: 500000, MaxTokens: 500000, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "xiaomi",
		Name:     "Xiaomi",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://api.xiaomimimo.com/v1",
		AuthEnv:  "[\"XIAOMI_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "mimo-v2.5", Name: "MiMo-V2.5", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mimo-v2.5-pro", Name: "MiMo-V2.5-Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mimo-v2.5-pro-ultraspeed", Name: "MiMo-V2.5-Pro-UltraSpeed", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mimo-v2.6-flash", Name: "MiMo-V2.6-Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "极速 MoE / 高效响应"},
			{ID: "mimo-v2.6-pro", Name: "MiMo-V2.6-Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mimo-v2.6-pro-ultraspeed", Name: "MiMo-V2.6-Pro-UltraSpeed", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "xiaomi-token-plan-ams",
		Name:     "Xiaomi Token Plan AMS",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://token-plan-ams.xiaomimimo.com/v1",
		AuthEnv:  "[\"XIAOMI_TOKEN_PLAN_AMS_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "mimo-v2.5", Name: "MiMo-V2.5", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mimo-v2.5-pro", Name: "MiMo-V2.5-Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mimo-v2.6-flash", Name: "MiMo-V2.6-Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "极速 MoE / 高效响应"},
			{ID: "mimo-v2.6-pro", Name: "MiMo-V2.6-Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "xiaomi-token-plan-cn",
		Name:     "Xiaomi Token Plan CN",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://token-plan-cn.xiaomimimo.com/v1",
		AuthEnv:  "[\"XIAOMI_TOKEN_PLAN_CN_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "mimo-v2.5", Name: "MiMo-V2.5", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mimo-v2.5-pro", Name: "MiMo-V2.5-Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mimo-v2.6-flash", Name: "MiMo-V2.6-Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "极速 MoE / 高效响应"},
			{ID: "mimo-v2.6-pro", Name: "MiMo-V2.6-Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "xiaomi-token-plan-sgp",
		Name:     "Xiaomi Token Plan SGP",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://token-plan-sgp.xiaomimimo.com/v1",
		AuthEnv:  "[\"XIAOMI_TOKEN_PLAN_SGP_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "mimo-v2.5", Name: "MiMo-V2.5", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mimo-v2.5-pro", Name: "MiMo-V2.5-Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "mimo-v2.6-flash", Name: "MiMo-V2.6-Flash", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "极速 MoE / 高效响应"},
			{ID: "mimo-v2.6-pro", Name: "MiMo-V2.6-Pro", Reasoning: true, ContextWindow: 1048576, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "zai",
		Name:     "Z.AI",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://api.z.ai/api/coding/paas/v4",
		AuthEnv:  "[\"ZAI_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "glm-4.7", Name: "GLM-4.7", Reasoning: true, ContextWindow: 204800, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5-turbo", Name: "GLM-5-Turbo", Reasoning: true, ContextWindow: 200000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.2", Name: "GLM-5.2", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.2-highspeed", Name: "GLM-5.2 Highspeed", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.3", Name: "GLM-5.3", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.3-flash", Name: "GLM-5.3-Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "极速 MoE / 高效响应"},
			{ID: "glm-5.3-highspeed", Name: "GLM-5.3 Highspeed", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
	{
		ID:       "zai-coding-cn",
		Name:     "Z.AI Coding CN",
		Type:     "api_key",
		Category: "china",
		BaseURL:  "https://open.bigmodel.cn/api/coding/paas/v4",
		AuthEnv:  "[\"ZAI_CODING_CN_API_KEY\"]",
		AuthType: "bearer",
		Models: []RegistryModel{
			{ID: "glm-4.6v", Name: "GLM-4.6V", Reasoning: true, ContextWindow: 128000, MaxTokens: 32768, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.3", Name: "GLM-5.3", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
			{ID: "glm-5.3-flash", Name: "GLM-5.3-Flash", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "极速 MoE / 高效响应"},
			{ID: "glm-5.3-highspeed", Name: "GLM-5.3 Highspeed", Reasoning: true, ContextWindow: 1000000, MaxTokens: 131072, Description: "深度长链推理模型 (Reasoning)"},
		},
	},
}

// GetProviderMeta finds provider metadata by provider ID or alias.
func GetProviderMeta(providerID string) (ProviderMeta, bool) {
	for _, p := range OfficialProviders {
		if p.ID == providerID {
			return p, true
		}
	}
	return ProviderMeta{}, false
}

// ThinkingLevelDef represents the metadata and description of a thinking depth level in Pi Agent.
type ThinkingLevelDef struct {
	Level  string `json:"level"`
	Label  string `json:"label"`
	Tokens string `json:"tokens"`
	Desc   string `json:"desc"`
}

// OfficialThinkingDefinitions maps the standard Pi Agent thinking levels to human-readable metadata.
var OfficialThinkingDefinitions = map[string]ThinkingLevelDef{
	"off":     {Level: "off", Label: "关闭思考", Tokens: "0", Desc: "不进行推理 (No reasoning)"},
	"minimal": {Level: "minimal", Label: "极简思考", Tokens: "~1k", Desc: "极短推理 (~1k tokens)"},
	"low":     {Level: "low", Label: "浅度思考", Tokens: "~2k", Desc: "轻度推理 (~2k tokens)"},
	"medium":  {Level: "medium", Label: "中度思考", Tokens: "~8k", Desc: "中度推理 (~8k tokens)"},
	"high":    {Level: "high", Label: "深度思考", Tokens: "~16k", Desc: "深度推理 (~16k tokens)"},
	"xhigh":   {Level: "xhigh", Label: "超深思考", Tokens: "~32k", Desc: "超高推理 (~32k tokens)"},
	"max":     {Level: "max", Label: "极限思考", Tokens: "Max", Desc: "最大算力推理 (Max reasoning)"},
}

// GetModelThinkingCapabilities returns the thinking levels supported by a model based on official Pi Agent rules.
func GetModelThinkingCapabilities(modelID, providerID string) (levels []string, defaultLevel string) {
	m := strings.ToLower(modelID)
	p := strings.ToLower(providerID)

	// 1. DeepSeek
	if strings.Contains(m, "deepseek-flash") {
		return []string{"low", "high", "max"}, "high"
	}
	if strings.Contains(m, "deepseek-v4-pro") || strings.Contains(m, "deepseek-reasoner") || strings.Contains(m, "deepseek-r1") || strings.Contains(m, "-r1") {
		return []string{"high", "max"}, "high"
	}
	if strings.Contains(m, "deepseek-chat") || strings.Contains(m, "deepseek-v3") {
		return nil, ""
	}

	// 2. Xiaomi MiMo
	if strings.Contains(m, "mimo-v2.5") || strings.Contains(m, "mimo-v2.6") {
		if strings.Contains(m, "asr") || strings.Contains(m, "tts") {
			return nil, ""
		}
		return []string{"off", "high"}, "high"
	}

	// 3. Google Gemini / Gemma
	if strings.Contains(m, "gemini-3.1-pro") || strings.Contains(m, "gemini-3.7") || strings.Contains(m, "gemini-3.8") || strings.Contains(m, "gemini-4") {
		return []string{"low", "medium", "high"}, "high"
	}
	if strings.Contains(m, "gemini-3") || strings.Contains(m, "gemini-3.5") || strings.Contains(m, "gemini-flash-latest") || strings.Contains(m, "gemini-flash-lite-latest") {
		return []string{"minimal", "low", "medium", "high"}, "high"
	}
	if strings.Contains(m, "gemma-4") || strings.Contains(m, "gemini-3.1-flash-lite-image") {
		return []string{"minimal", "high"}, "high"
	}
	if strings.Contains(m, "gemini-2.5") {
		return []string{"off", "high"}, "high"
	}

	// 4. Anthropic Claude
	if strings.Contains(m, "claude-sonnet-5-5") || strings.Contains(m, "claude-opus-5-5") {
		return []string{"low", "medium", "high", "xhigh", "max"}, "high"
	}
	if strings.Contains(m, "claude-opus-4-7") || strings.Contains(m, "claude-opus-4-8") || strings.Contains(m, "claude-sonnet-5") || strings.Contains(m, "claude-fable") {
		return []string{"xhigh", "max"}, "max"
	}
	if strings.Contains(m, "claude-3-7") || strings.Contains(m, "claude-sonnet-4-5") || strings.Contains(m, "claude-opus-4-6") {
		return []string{"off", "low", "medium", "high", "max"}, "high"
	}

	// 5. OpenAI
	if strings.Contains(m, "gpt-5.1") || strings.Contains(m, "gpt-5.2") || strings.Contains(m, "gpt-5.3") || strings.Contains(m, "gpt-5.4") || strings.Contains(m, "gpt-5.5") || strings.Contains(m, "gpt-5.6") || strings.Contains(m, "gpt-6") {
		return []string{"off", "low", "medium", "high", "xhigh", "max"}, "high"
	}
	if strings.HasPrefix(m, "gpt-5") {
		return []string{"off", "minimal", "low", "medium", "high"}, "high"
	}
	if strings.HasPrefix(m, "o1") || strings.HasPrefix(m, "o3") || strings.HasPrefix(m, "o4") {
		return []string{"low", "medium", "high"}, "medium"
	}

	// 6. xAI
	if strings.Contains(m, "grok-4") {
		return []string{"off", "low", "medium", "high"}, "high"
	}

	// 7. Z.AI / GLM
	if strings.Contains(m, "glm-5.3-highspeed") || strings.Contains(m, "glm-5.2-highspeed") {
		return []string{"low", "high", "max"}, "high"
	}
	if strings.Contains(m, "glm-5.3") {
		return []string{"off", "low", "high"}, "high"
	}

	// 8. Qwen / Aliyun
	if strings.Contains(m, "deepseek-v4-flash") {
		return []string{"low", "high", "max"}, "high"
	}
	if strings.Contains(m, "deepseek-v4-pro") {
		return []string{"high", "max"}, "high"
	}

	// 9. Kimi Coding
	if strings.Contains(m, "kimi-for-coding") || strings.Contains(m, "k3") {
		return []string{"off", "high"}, "high"
	}

	// 10. General heuristic fallback (e.g. for Ollama / Custom endpoints)
	if strings.Contains(m, "reasoner") || strings.Contains(m, "reasoning") || strings.Contains(m, "thinking") || strings.Contains(m, "thought") || strings.Contains(p, "reasoning") {
		return []string{"off", "low", "medium", "high"}, "high"
	}

	return nil, ""
}

<p align="center">
  <h1 align="center">Autonomous Agent Development Workflow</h1>
  <p align="center">
    <strong>A World-Class, Single-Binary, Autonomous AI Software Engineering Agent in Pure Go</strong>
  </p>
  <p align="center">
    Zero External Runtime Dependencies • ReAct Autonomous Loop • 43+ Model Providers • 7-Level Reasoning Depth • 3-Tier HITL Guardrails • Live Terminal Streaming • Git Rollback
  </p>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/Build-Passing-brightgreen.svg" alt="Build Status"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/Tests-100%25%20Passing-success.svg" alt="Test Status"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/pulls"><img src="https://img.shields.io/badge/PRs-Welcome-brightgreen.svg" alt="PRs Welcome"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow"><img src="https://img.shields.io/badge/Providers-43%2B%20Supported-purple.svg" alt="Providers"></a>
</p>

---

<p align="center">
  <strong>Language Switcher / 语言切换</strong><br>
  <a href="README.md">English</a> •
  <a href="README.zh-CN.md">简体中文</a> •
  <a href="README.ja.md">日本語</a> •
  <a href="README.ko.md">한국어</a> •
  <a href="README.de.md">Deutsch</a> •
  <a href="README.fr.md">Français</a> •
  <a href="README.ru.md">Русский</a> •
  <a href="README.pt.md">Português</a> •
  <a href="README.es.md">Español</a>
</p>

---

## 📖 Table of Contents

- [Overview](#-overview)
- [Key Features](#-key-features)
- [System Architecture](#-system-architecture)
- [Pi-Aligned Thinking Depth Selector](#-pi-aligned-thinking-depth-selector)
- [Human-in-the-Loop (HITL) Security](#-human-in-the-loop-hitl-security)
- [Multi-Provider Registry (43+ Platforms)](#-multi-provider-registry-43-platforms)
- [Quick Start](#-quick-start)
  - [Prerequisites](#prerequisites)
  - [Run from Source](#1-run-from-source)
  - [Build Single Binary](#2-build-single-binary)
- [Configuration Guide](#-configuration-guide)
- [Interactive Web UI](#-interactive-web-ui)
- [Testing & Quality Assurance](#-testing--quality-assurance)
- [Contributing](#-contributing)
- [Author & License](#-author--license)

---

## 🌟 Overview

**Autonomous Agent Development Workflow** is an ultra-lightweight, high-performance, single-binary AI software development agent built from the ground up in Go. It empowers software developers and engineering teams with an autonomous pairing partner capable of planning tasks, reading and writing code, executing terminal commands, self-healing compilation errors, and exploring alternative design branches.

Unlike Python-heavy frameworks requiring gigabytes of virtual environments and node dependencies, this agent compiles into a **single ~9MB binary** with its entire Google Material 3 / Gemini-inspired web interface embedded directly via Go's `//go:embed`.

---

## ✨ Key Features

1. **Zero-Dependency Single Binary (`//go:embed`)**
   - The entire frontend (HTML, Tailwind CSS, offline Markdown engine, syntax icons) is compiled into one self-contained executable.
   - Run anywhere on Windows, Linux, or macOS with zero installation overhead.

2. **Autonomous ReAct Loop & Self-Healing**
   - Implements `Think -> Act -> Observe -> Reflect`.
   - When test suites or builds fail, the agent intercepts the stdout/stderr, diagnoses the root cause, and autonomously refactors code until green.

3. **Pi-Aligned 7-Level Thinking Depth Control**
   - Real-time reasoning effort slider/popover in the prompt bar: `off`, `minimal`, `low`, `medium`, `high`, `xhigh`, and `max`.
   - Dynamically adapts available levels per model (e.g. DeepSeek V4.1 Flash, Xiaomi MiMo, Claude 3.7 Sonnet, OpenAI o-series, Gemini 3).

4. **3-Tier HITL Safety Gatekeeper (Antigravity Guardrails)**
   - **Manual**: Strict review required for every file modification and command.
   - **Guarded**: Safe diagnostics and non-destructive reads run automatically; file overwrites, deletions, network requests, and dangerous commands require 1-click human authorization.
   - **Auto**: Full autonomy for air-gapped or sandboxed test runs.

5. **Live Terminal Streaming (Real-Time Output)**
   - Long-running commands stream output line-by-line via RFC 6455 WebSockets into an embedded dark code terminal.

6. **Git-Aware Automatic Rollback**
   - Monitors repository status in real-time. Provides single-click atomic rollback to undo changes and clean untracked files safely.

7. **Non-Destructive Session Branch Tree (Time-Travel)**
   - Automatically takes snapshots after each task turn. Jump between branches without losing historical reasoning paths.

8. **Token Telemetry & Cost Meter**
   - Live tracking of prompt and completion tokens with real-time currency calculation.

---

## 🏗️ System Architecture

```text
┌────────────────────────────────────────────────────────────────────────┐
│               Embedded Google Material 3 Web Interface                │
│    (Prompt Bar, Thinking Selector, Session Tree, Live Terminal)        │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ RFC 6455 WebSockets / REST
┌───────────────────────────────────▼────────────────────────────────────┐
│                    DevAgent Go Core (agent.exe)                        │
│                                                                        │
│   ┌───────────────────────────┐    ┌───────────────────────────────┐   │
│   │    ReAct Engine Core      │◄───┤    Multi-Provider Hub         │   │
│   │  (Think-Act-Observe-Reflect)   │    (DeepSeek, Xiaomi, Gemini, │   │
│   └─────────────┬─────────────┘    │     Claude, OpenAI, Ollama)   │   │
│                 │                  └───────────────────────────────┘   │
│                 ▼                                                      │
│   ┌───────────────────────────┐    ┌───────────────────────────────┐   │
│   │  HITL Security Gatekeeper │    │   Session Storage & Branch    │   │
│   │  (Path Jail, Regex Filter)│    │   (Snapshots, Time-Travel)    │   │
│   └─────────────┬─────────────┘    └───────────────────────────────┘   │
│                 │                                                      │
│                 ▼                                                      │
│   ┌────────────────────────────────────────────────────────────┐       │
│   │                  Built-in Tool Ecosystem                   │       │
│   │  file_read • file_write • file_edit • list_dir • cmd_exec  │       │
│   └────────────────────────────────────────────────────────────┘       │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                         ┌──────────▼──────────┐
                         │   Host Filesystem   │
                         │    & Shell/Git      │
                         └─────────────────────┘
```

---

## 🧠 Pi-Aligned Thinking Depth Selector

The prompt bar features an integrated Thinking Depth controller calibrated to the exact official model specifications of the Pi Agent framework:

| Level | UI Name | Token Budget | Typical Use Case |
| :--- | :--- | :--- | :--- |
| `off` | **关闭思考 (Off)** | 0 Tokens | Direct responses, formatting, conversational Q&A |
| `minimal` | **极简思考 (Minimal)** | ~1k Tokens | Quick summarization, single-line lookups |
| `low` | **浅度思考 (Low)** | ~2k Tokens | Minor refactorings, quick script generation |
| `medium` | **中度思考 (Medium)** | ~8k Tokens | Standard feature implementation and debugging |
| `high` | **深度思考 (High)** | ~16k Tokens | Complex algorithmic design, full unit-test loops |
| `xhigh` | **超深思考 (Extra-High)**| ~32k Tokens | Architecture-level synthesis, cross-module refactors |
| `max` | **极限思考 (Max)** | Max Capacity | Exhaustive multi-step proof and deep bug root-causing |

*Models automatically expose only their valid levels. Models without reasoning support (such as GPT-4o) gracefully render a disabled indicator.*

---

## 🛡️ Human-in-the-Loop (HITL) Security

- **Path Confinement**: Eliminates directory traversal (`../`) attacks, anchoring operations inside the designated workspace.
- **Dangerous Command Interception**: Patterns matching destructive operations (`rm -rf /`, formatting drives, credential access) are blocked or elevated to explicit approval.
- **Visual Diff Inspection**: File modifications highlight surgical diff replacements in red/green before application.

---

## 🌐 Multi-Provider Registry (43+ Platforms)

Built-in support with verified Base URLs and recommended model matrices:

- **Flagship Cloud Providers**: DeepSeek, Xiaomi (MiMo), Google Gemini, Anthropic Claude, OpenAI, xAI (Grok), Mistral, Z.AI (GLM), Moonshot/Kimi.
- **Dedicated Coding Gateways**: Kimi for Coding, Z.AI Coding CN, Qwen Token Plan, GitHub Copilot.
- **High-Speed Inference Chips**: Groq, Cerebras, Fireworks, Together AI, Baseten.
- **Cloud & Edge Gateways**: Cloudflare AI Gateway, Vercel AI Gateway, Azure OpenAI, AWS Bedrock.
- **Local Private Models**: Ollama, llama.cpp, LocalAI, vLLM.

---

## 🚀 Quick Start

### Prerequisites
- [Go 1.22+](https://golang.org/dl/) (only required if compiling from source)
- Git (optional, for version control integration)

### 1. Run from Source
```bash
# Clone the repository
git clone https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow.git
cd Autonomous-Agent-Development-Workflow

# Launch DevAgent server (defaults to port 8080)
go run cmd/agent/main.go -port 8080 -dir .
```

Open your browser at:
```
http://127.0.0.1:8080
```

### 2. Build Single Binary
Generate a stripped, production-grade executable:
```bash
# On Linux/macOS
go build -ldflags "-s -w" -o agent ./cmd/agent

# On Windows
go build -ldflags "-s -w" -o agent.exe ./cmd/agent
```

---

## ⚙️ Configuration Guide

Copy `config.example.json` to `config.json` (or configure dynamically through the Settings UI):

```bash
cp config.example.json config.json
```

Sample profile definition:
```json
{
  "active_profile_id": "deepseek-flash",
  "profiles": [
    {
      "id": "deepseek-flash",
      "name": "DeepSeek V4.1 Flash",
      "provider": "deepseek",
      "base_url": "https://api.deepseek.com",
      "api_key": "YOUR_API_KEY_HERE",
      "model": "deepseek-flash",
      "max_steps": 15
    }
  ],
  "workspace_dir": ".",
  "approval_mode": "guarded"
}
```

*Note: Environment variables such as `DEEPSEEK_API_KEY`, `OPENAI_API_KEY`, `GEMINI_API_KEY`, and `ANTHROPIC_API_KEY` are automatically detected on startup.*

---

## 🖥️ Interactive Web UI

- **Google Gemini Theme**: Styled with Google Sans typography, JetBrains Mono code fonts, and dark surface aesthetics (`#131314`).
- **Collapsible Thinking Accordions**: Deep thought logs expand during reasoning and smoothly collapse upon completion.
- **Responsive Navigation**: Top bar adapts to any screen width and laptop scaling without element clipping.
- **Native OS Folder Picker**: Click the workspace button to launch the native operating system directory dialog.

---

## 🧪 Testing & Quality Assurance

Every pull request and release is validated across four quality gates:

```bash
# Run all unit, integration, and security tests
go test -v ./...
```

Included verification test suites:
- `pkg/security`: Path containment, dangerous command interception, and approval dispatch.
- `pkg/tools`: File surgical block replacement, directory iteration, and command execution streaming.
- `pkg/server`: Dynamic model profile switching, thinking level validation, and HTTP API routing.
- `pkg/git`: Rollback safety and repository status tracking.
- `tests`: End-to-end multi-turn ReAct autonomous agent simulation.

---

## 🤝 Contributing

Contributions are welcomed and greatly appreciated! Please review [CONTRIBUTING.md](CONTRIBUTING.md) and the [Code of Conduct](CODE_OF_CONDUCT.md) before opening a pull request.

---

## 👤 Author & License

- **Author**: [yifangdong](https://github.com/yfddxwx)
- **Repository**: [Autonomous-Agent-Development-Workflow](https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow)
- **License**: [MIT License](LICENSE) © 2026 yifangdong

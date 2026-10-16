<p align="center">
  <h1 align="center">自主研发 Agent 工程工作流 (Autonomous Agent Development Workflow)</h1>
  <p align="center">
    <strong>基于 Go 语言构建的世界级极轻量、单二进制自包含自动化 AI 研发助理</strong>
  </p>
  <p align="center">
    零外部运行时依赖 • ReAct 自主规划与反思自愈 • 43+ 大模型厂商矩阵 • 7 级思考深度调控 • 三级人机协同安全门禁 • 终端实时流式回显 • Git 秒级撤销回滚
  </p>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/构建-通过-brightgreen.svg" alt="Build Status"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/测试-100%25%20通过-success.svg" alt="Test Status"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/pulls"><img src="https://img.shields.io/badge/PRs-欢迎提交-brightgreen.svg" alt="PRs Welcome"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow"><img src="https://img.shields.io/badge/支持厂商-43%2B-purple.svg" alt="Providers"></a>
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

## 📖 目录

- [项目概述](#-项目概述)
- [核心特性](#-核心特性)
- [系统架构图](#-系统架构图)
- [对标 Pi 的 7 级思考深度选择器](#-对标-pi-的-7-级思考深度选择器)
- [人机协同 (HITL) 安全门禁](#-人机协同-hitl-安全门禁)
- [43+ 大模型厂商纳管矩阵](#-43-大模型厂商纳管矩阵)
- [快速上手指南](#-快速上手指南)
  - [环境准备](#环境准备)
  - [源码运行](#1-源码直接运行)
  - [单二进制编译打包](#2-编译自包含可执行程序)
- [配置文件说明](#-配置文件说明)
- [现代化交互界面体验](#-现代化交互界面体验)
- [自动化工程测试套件](#-自动化工程测试套件)
- [贡献指南](#-贡献指南)
- [作者与开源许可证](#-作者与开源许可证)

---

## 🌟 项目概述

**Autonomous Agent Development Workflow** 是一个完全使用 **Go 语言** 从零构建的高性能、极轻量、单二进制本地自主研发 Agent。它为开发者和研发团队提供了一位常驻本地的高智商结对编程助理，能够独立进行任务拆解与推演、代码全生命周期读写、系统终端命令执行、编译报错自动修复以及历史分支时间旅行。

区别于传统基于 Python 堆叠数个 G 虚拟环境与 Node 复杂依赖的庞大框架，本项目通过 Go 原生 `//go:embed` 技术将全套现代化 Web 前端一体化嵌入，编译后仅为一个 **~9MB 的自包含单文件**，即下即用，秒级启动。

---

## ✨ 核心特性

1. **零外部运行时依赖单二进制 (`//go:embed`)**
   - HTML、Tailwind CSS、离线 Markdown 渲染器与矢量图标全部静态嵌入二进制，无需安装 Python、Node.js、NPM 等任何宿主环境。
2. **ReAct 自主循环与自愈反思 (Self-Healing Loop)**
   - 贯彻 `Think -> Act -> Observe -> Reflect` 标准循环。遇到编译报错、测试红灯或命令异常时，自动拦截输出日志并自我反思迭代，直至实现任务目标。
3. **对标 Pi Agent 的 7 级思考深度选择器 (Thinking Depth Selector)**
   - 输入框原生内置思考深度控制胶囊与弹出菜单，支持 `关闭思考 (off)`、`极简思考 (minimal)`、`浅度思考 (low)`、`中度思考 (medium)`、`深度思考 (high)`、`超深思考 (xhigh)` 与 `极限思考 (max)`。
   - 智能识别各模型特征（如 DeepSeek V4.1 Flash、小米 MiMo、Claude 3.7、OpenAI o-系列、Gemini 3 等）并动态展示专属档位。
4. **Antigravity 三级人机协同审核 (HITL Gatekeeper)**
   - **完全不审核 (Auto)**：适合沙箱隔离或自动化批量任务。
   - **有限审核 (Guarded)**：只读诊断与安全测试自动放行；文件删除、已有文件覆写、破坏性指令及外部网络请求强制唤起人工审查弹窗。
   - **人工审核 (Manual)**：每一项工具调用均需开发者明确批准。
5. **终端命令实时流式回显 (Live Terminal Streaming)**
   - 终端命令通过 RFC 6455 WebSocket 逐行流式推送到前端深色终端卡片，告别长时间黑盒等待。
6. **Git 感知与秒级安全回滚 (Git-Aware Rollback)**
   - 实时监听本地仓库代码变动，支持一键丢弃修改并安全清理未追踪文件。
7. **会话分支树历史检查点与时间旅行 (Session Tree)**
   - 任务每轮执行完毕后自动建立快照分支，支持在历史轮次间自由跳转并衍生全新分支，探索历史永不丢失。
8. **Token 消耗与成本实时仪表盘**
   - 实时折算 Prompt / Completion Token 并按行业标准费率精确统计人民币与美元成本。

---

## 🏗️ 系统架构图

```text
┌────────────────────────────────────────────────────────────────────────┐
│               内嵌 Google Material 3 / Gemini 风格前端 UI              │
│    (任务输入框、思考深度滑块、分支树下拉选择器、实时代码流式终端)      │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ RFC 6455 WebSocket / RESTful API
┌───────────────────────────────────▼────────────────────────────────────┐
│                    DevAgent Go 核心引擎 (agent.exe)                    │
│                                                                        │
│   ┌───────────────────────────┐    ┌───────────────────────────────┐   │
│   │    ReAct 自主循环核心     │◄───┤    多模型统一适配中枢         │   │
│   │  (Think-Act-Observe-Reflect)   │    (DeepSeek、Xiaomi、Gemini、│   │
│   └─────────────┬─────────────┘    │     Claude、OpenAI、Ollama等) │   │
│                 │                  └───────────────────────────────┘   │
│                 ▼                                                      │
│   ┌───────────────────────────┐    ┌───────────────────────────────┐   │
│   │  HITL 人机安全协同门禁    │    │    会话分支持久化与时间旅行   │   │
│   │ (路径防穿越沙箱、高危过滤)│    │   (快照分支树、会话持久化)    │   │
│   └─────────────┬─────────────┘    └───────────────────────────────┘   │
│                 │                                                      │
│                 ▼                                                      │
│   ┌────────────────────────────────────────────────────────────┐       │
│   │                     本地核心工具生态                       │       │
│   │  file_read • file_write • file_edit • list_dir • cmd_exec  │       │
│   └────────────────────────────────────────────────────────────┘       │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                         ┌──────────▼──────────┐
                         │   宿主操作系统环境  │
                         │    文件系统与终端   │
                         └─────────────────────┘
```

---

## 🧠 对标 Pi 的 7 级思考深度选择器

在输入框操作栏右侧集成了深度思考调控胶囊，严格遵循 Pi 官方底层规约：

| 档位代号 | 界面名称 | Token 预算 | 典型应用场景 |
| :--- | :--- | :--- | :--- |
| `off` | **关闭思考** | 0 Tokens | 格式整理、直接问答、简单配置编写 |
| `minimal` | **极简思考** | ~1k Tokens | 单行代码修改、快速信息检索 |
| `low` | **浅度思考** | ~2k Tokens | 常规脚本片段编写、轻量级小重构 |
| `medium` | **中度思考** | ~8k Tokens | 标准业务功能实现、跨函数逻辑排查 |
| `high` | **深度思考** | ~16k Tokens | 复杂算法推演、单元测试闭环自愈 (默认推荐) |
| `xhigh` | **超深思考** | ~32k Tokens | 复杂架构推演、全仓库规模设计分析 |
| `max` | **极限思考** | 顶格极限算力 | 极限深度长链逻辑推导与攻坚排障 |

*若当前选中的模型不支持深度思考（如 GPT-4o），按钮将自动平滑置灰并禁用，避免无效传参。*

---

## 🛡️ 人机协同 (HITL) 安全门禁

- **工作空间路径严密隔离**：防止任何 `../` 越权路径逃逸，所有操作严格锚定在用户指定的工作目录内。
- **高危破坏性指令拦截**：自动捕获格式化、递归删除根目录等危险操作，强制挂起并交由开发者审查。
- **差异化代码替换审计**：支持精准文本块增量替换（`file_edit`），审查弹窗以红绿对比色高亮显示待修改内容。

---

## 🌐 43+ 大模型厂商纳管矩阵

预设权威官方 Base URL 与推荐主力模型：

- **国内前沿模型**：DeepSeek、小米 MiMo、月之暗面 Kimi、智谱 Z.AI、阿里通义 Qwen、MiniMax、蚂蚁 Ant Ling。
- **国际主流模型**：Anthropic Claude、OpenAI、Google Gemini、xAI (Grok)、Mistral、Meta Llama。
- **专用极速通道**：Kimi For Coding、Z.AI Coding CN、Qwen Token Plan、GitHub Copilot。
- **极速推理算力**：Groq、Cerebras、Fireworks、Together AI、Baseten。
- **边缘网关与本地私有**：Cloudflare AI Gateway、Vercel AI Gateway、Azure、AWS Bedrock、Ollama、llama.cpp、vLLM。

---

## 🚀 快速上手指南

### 环境准备
- [Go 1.22+](https://golang.org/dl/)（仅从源码编译时需要）
- Git（可选，用于版本回滚特性联动）

### 1. 源码直接运行
```bash
# 克隆代码仓库
git clone https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow.git
cd Autonomous-Agent-Development-Workflow

# 启动 Agent 本地服务 (默认监听 8080 端口)
go run cmd/agent/main.go -port 8080 -dir .
```

浏览器访问：
```
http://127.0.0.1:8080
```

### 2. 编译自包含可执行程序
生成独立免安装二进制可执行文件：
```bash
# Linux / macOS 环境
go build -ldflags "-s -w" -o agent ./cmd/agent

# Windows 环境
go build -ldflags "-s -w" -o agent.exe ./cmd/agent
```

---

## ⚙️ 配置文件说明

可将 `config.example.json` 复制为 `config.json`，或直接在 Web 设置界面图形化添加：

```bash
cp config.example.json config.json
```

配置样例：
```json
{
  "active_profile_id": "deepseek-flash",
  "profiles": [
    {
      "id": "deepseek-flash",
      "name": "DeepSeek V4.1 Flash",
      "provider": "deepseek",
      "base_url": "https://api.deepseek.com",
      "api_key": "在此填入您的API_KEY",
      "model": "deepseek-flash",
      "max_steps": 15
    }
  ],
  "workspace_dir": ".",
  "approval_mode": "guarded"
}
```

*系统同时支持从系统环境变量（如 `DEEPSEEK_API_KEY`、`OPENAI_API_KEY`、`GEMINI_API_KEY`）自动装载默认密钥。*

---

## 🖥️ 现代化交互界面体验

- **Google Gemini 品牌设计规范**：采用 Google Sans 字体、JetBrains Mono 终端字体及 `#131314` 深灰表面色调。
- **折叠式思考动效**：大模型思考过程在生成中展开显示，回答完成后平滑收起。
- **自适应视口排版**：顶栏与操作控件经受住各类屏幕缩放考验，无遮挡与越界。
- **唤起原生文件资源管理器**：点击工作空间按钮可直接唤起系统图形文件选择器选择任意文件夹。

---

## 🧪 自动化工程测试套件

本项目执行严苛的工程四重门禁标准：

```bash
# 执行全量单元测试与集成测试
go test -v ./...
```

测试覆盖范围：
- `pkg/security`：路径越权阻断、高危命令拦截、人工审核挂起与放行机制。
- `pkg/tools`：精确差异文本替换、多行切片查看、命令实时流式捕获。
- `pkg/server`：动态模型配置加载、思考深度解析校验、REST 路由健全性。
- `pkg/git`：一键秒级快照回滚与仓库状态检测。
- `tests`：模拟大模型多轮调用与自动化反思自愈全链路集成。

---

## 🤝 贡献指南

欢迎各位开发者提交 Issue 或 Pull Request！在提交前请详细参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 与 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)。

---

## 👤 作者与开源许可证

- **作者**：[yifangdong](https://github.com/yfddxwx)
- **代码仓库**：[Autonomous-Agent-Development-Workflow](https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow)
- **开源许可证**：[MIT License](LICENSE) © 2026 yifangdong

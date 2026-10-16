package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"agent/pkg/engine"
	"agent/pkg/llm"
	"agent/pkg/models"
	"agent/pkg/security"
	"agent/pkg/server"
	"agent/pkg/storage"
	"agent/pkg/tools"
	"agent/web"
)

func main() {
	port := flag.Int("port", 8080, "Web UI 监听端口")
	workspace := flag.String("dir", ".", "Agent 操作的工作空间目录")
	provider := flag.String("provider", "deepseek", "默认模型协议 (deepseek / gpt-4o / gemini / ollama)")
	baseURL := flag.String("base-url", "", "模型 API Base URL")
	apiKey := flag.String("api-key", "", "模型 API Key（默认优先读取环境变量 AGENT_API_KEY 或 DEEPSEEK_API_KEY）")
	model := flag.String("model", "", "默认调用的模型名称")
	requireApproval := flag.Bool("require-approval", true, "兼容参数: 是否开启交互审核")
	approvalModeFlag := flag.String("approval-mode", "guarded", "Antigravity 审核模式: manual(人工审核) / guarded(有限审核) / auto(完全不审核)")
	maxSteps := flag.Int("max-steps", 15, "Agent 单任务最大自主循环轮数")
	flag.Parse()

	// Track which CLI flags were explicitly provided
	flagsProvided := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		flagsProvided[f.Name] = true
	})

	// Resolve absolute workspace
	cleanDir := strings.Trim(strings.TrimSpace(*workspace), "\"'")
	if len(flag.Args()) > 0 {
		fullAttempt := cleanDir + " " + strings.Join(flag.Args(), " ")
		if _, err := os.Stat(fullAttempt); err == nil {
			cleanDir = fullAttempt
		}
	}
	absWorkspace, err := filepath.Abs(cleanDir)
	if err != nil {
		absWorkspace = cleanDir
	}

	// Initialize Storage Manager for persistent configuration and session history
	storageMgr := storage.NewManager(absWorkspace)
	appCfg := storageMgr.GetConfig()
	activeProf := storageMgr.GetActiveProfile()

	usedProvider := activeProf.Provider
	if flagsProvided["provider"] {
		usedProvider = *provider
	}

	usedBaseURL := activeProf.BaseURL
	if flagsProvided["base-url"] {
		usedBaseURL = *baseURL
	}

	usedAPIKey := activeProf.APIKey
	if flagsProvided["api-key"] {
		usedAPIKey = *apiKey
	}
	if usedAPIKey == "" {
		usedAPIKey = os.Getenv("AGENT_API_KEY")
	}
	if usedAPIKey == "" {
		usedAPIKey = os.Getenv("DEEPSEEK_API_KEY")
	}
	if usedAPIKey == "" {
		usedAPIKey = os.Getenv("OPENAI_API_KEY")
	}
	if usedAPIKey == "" {
		usedAPIKey = os.Getenv("GEMINI_API_KEY")
	}

	usedModel := activeProf.Model
	if flagsProvided["model"] {
		usedModel = *model
	}

	usedApprovalMode := appCfg.ApprovalMode
	if flagsProvided["approval-mode"] {
		usedApprovalMode = models.ApprovalMode(*approvalModeFlag)
	}
	if flagsProvided["require-approval"] && !*requireApproval {
		usedApprovalMode = models.ApprovalModeAuto
	}
	if usedApprovalMode != models.ApprovalModeManual && usedApprovalMode != models.ApprovalModeGuarded && usedApprovalMode != models.ApprovalModeAuto {
		usedApprovalMode = models.ApprovalModeGuarded
	}

	usedMaxSteps := activeProf.MaxSteps
	if usedMaxSteps <= 0 {
		usedMaxSteps = 15
	}
	if flagsProvided["max-steps"] {
		usedMaxSteps = *maxSteps
	}

	if usedBaseURL == "" {
		switch usedProvider {
		case "deepseek":
			usedBaseURL = "https://api.deepseek.com/v1"
			if usedModel == "" {
				usedModel = "deepseek-chat"
			}
		case "gpt-4o":
			usedBaseURL = "https://api.openai.com/v1"
			if usedModel == "" {
				usedModel = "gpt-4o"
			}
		case "gemini":
			usedBaseURL = "https://generativelanguage.googleapis.com/v1beta/openai"
			if usedModel == "" {
				usedModel = "gemini-1.5-pro"
			}
		case "ollama":
			usedBaseURL = "http://localhost:11434/v1"
			if usedModel == "" {
				usedModel = "llama3.1"
			}
		default:
			usedBaseURL = "https://api.deepseek.com/v1"
			if usedModel == "" {
				usedModel = "deepseek-chat"
			}
		}
	}

	config := &models.Config{
		WorkspaceDir:    absWorkspace,
		Port:            *port,
		Provider:        usedProvider,
		BaseURL:         usedBaseURL,
		APIKey:          usedAPIKey,
		Model:           usedModel,
		MaxSteps:        usedMaxSteps,
		RequireApproval: usedApprovalMode != models.ApprovalModeAuto,
		ApprovalMode:    usedApprovalMode,
	}

	// Initialize Security Manager with 3-tier approval mode
	sec := security.NewSecurityManager(absWorkspace, usedApprovalMode)

	// Initialize Tools
	reg := tools.NewRegistry()
	tools.RegisterBuiltinTools(reg, sec)

	// Initialize LLM Client
	llmClient := llm.NewOpenAIClient(usedBaseURL, usedAPIKey, usedModel)

	// Initialize ReAct Engine
	reactEngine := engine.NewReActEngine(llmClient, reg, sec, usedMaxSteps)

	// Initialize Server
	srv := server.NewServer(config, reactEngine, sec, reg, web.Assets, storageMgr)
	handler := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: handler,
	}

	// Graceful shutdown handling
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("=====================================================")
		log.Printf(" DevAgent 已启动成功！")
		log.Printf(" Web 访问地址: http://127.0.0.1:%d", *port)
		log.Printf(" 工作空间目录: %s", absWorkspace)
		log.Printf(" 默认模型提供商: %s (%s)", usedProvider, usedModel)
		log.Printf(" 审核模式 (Approval Mode): %s", usedApprovalMode)
		log.Printf("=====================================================")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	<-stopChan
	log.Println("正在优雅关闭 DevAgent 服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
	log.Println("DevAgent 已安全退出。")
}

<p align="center">
  <h1 align="center">Autonomous Agent Development Workflow</h1>
  <p align="center">
    <strong>순수 Go 언어로 구축된 세계적 수준의 단일 바이너리 자율형 AI 소프트웨어 엔지니어링 에이전트</strong>
  </p>
  <p align="center">
    외부 런타임 종속성 제로 • ReAct 자율 계획 및 자기 복구 루프 • 43개 이상의 LLM 모델 프로바이더 • 7단계 추론 깊이 제어 • 3단계 HITL 안전 검증 • 터미널 실시간 스트리밍 • Git 롤백
  </p>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/빌드-통과-brightgreen.svg" alt="Build Status"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/테스트-100%25%20통과-success.svg" alt="Test Status"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/pulls"><img src="https://img.shields.io/badge/PRs-환영-brightgreen.svg" alt="PRs Welcome"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow"><img src="https://img.shields.io/badge/프로바이더-43%2B-purple.svg" alt="Providers"></a>
</p>

---

<p align="center">
  <strong>Language Switcher / 언어 선택</strong><br>
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

## 📖 목차

- [프로젝트 개요](#-프로젝트-개요)
- [핵심 기능](#-핵심-기능)
- [시스템 아키텍처](#-시스템-아키텍처)
- [Pi 표준 7단계 추론 깊이 선택기](#-pi-표준-7단계-추론-깊이-선택기)
- [인간 참여형 (HITL) 보안 가드레일](#-인간-참여형-hitl-보안-가드레일)
- [43개 이상의 모델 프로바이더 지원](#-43개-이상의-모델-프로바이더-지원)
- [빠른 시작 가이드](#-빠른-시작-가이드)
  - [사전 요구 사항](#사전-요구-사항)
  - [소스 코드에서 실행](#1-소스-코드에서-실행)
  - [단일 바이너리 빌드](#2-단일-바이너리-빌드)
- [설정 가이드](#-설정-가이드)
- [모던 웹 인터페이스](#-모던-웹-인터페이스)
- [테스트 및 품질 검증](#-테스트-및-품질-검증)
- [기여하기](#-기여하기)
- [작성자 및 라이선스](#-작성자-및-라이선스)

---

## 🌟 프로젝트 개요

**Autonomous Agent Development Workflow**는 Go 언어로 처음부터 구축된 초경량, 고성능, 단일 바이너리 자율 AI 소프트웨어 개발 에이전트입니다. 개발자와 엔지니어링 팀에게 로컬 페어 프로그래밍 파트너를 제공하여 작업 계획 수립, 코드 작성 및 수정, 쉘 명령어 안전 실행, 컴파일 오류 자율 복구, 히스토리 분기 탐색을 원활하게 수행합니다.

Python의 거대한 가상 환경이나 복잡한 Node.js 의존성을 필요로 하지 않으며, Go의 `//go:embed` 기술을 통해 모던 웹 인터페이스를 단일 실행 파일에 완전히 포함하여 **약 9MB 단일 파일**로 배포됩니다.

---

## ✨ 핵심 기능

1. **외부 런타임 종속성 제로 단일 바이너리 (`//go:embed`)**
   - HTML, Tailwind CSS, 오프라인 마크다운 렌더러, 아이콘이 바이너리에 내장되어 있어 Node.js나 Python 설치 없이 즉시 실행할 수 있습니다.
2. **ReAct 자율 루프 및 자기 치유**
   - `Think -> Act -> Observe -> Reflect` 표준 사이클을 구현하여, 테스트 실패나 컴파일 에러 발생 시 오류 원인을 분석하고 자율적으로 코드를 수정합니다.
3. **Pi Agent 공식 사양 기반 7단계 추론 깊이 제어**
   - 입력창에서 실시간으로 조절 가능한 추론 깊이: `끄기 (off)`, `최소 (minimal)`, `낮음 (low)`, `중간 (medium)`, `깊음 (high)`, `매우 깊음 (xhigh)`, `최대 (max)`.
   - 모델별(DeepSeek, Xiaomi MiMo, Claude 3.7, OpenAI o-시리즈, Gemini 3 등) 지원 레벨을 자동으로 감지합니다.
4. **3단계 인간 참여형(HITL) 안전 검증**
   - **완전 자율 (Auto)**: 격리된 샌드박스 환경용.
   - **제한 승인 (Guarded)**: 읽기 및 진단 명령은 자동 실행, 파일 덮어쓰기/삭제/네트워크 호출 등은 1클릭 수동 승인 모달 표시.
   - **수동 승인 (Manual)**: 모든 도구 실행에 대해 개발자의 명시적 승인 요구.
5. **터미널 실시간 스트리밍**
   - 긴 명령어도 RFC 6455 WebSocket을 통해 웹 콘솔에 실시간 줄 단위로 스트리밍됩니다.
6. **Git 감지 및 1초 안전 롤백**
   - 변경 사항을 실시간 모니터링하며, 원클릭으로 수정 사항을 취소하고 미추적 파일을 안전하게 정리합니다.
7. **비파괴 세션 분기 트리 (타임머신)**
   - 턴마다 스냅샷을 생성하여 이전 분기로 안전하게 되돌아가고 새로운 탐색을 이어갈 수 있습니다.
8. **실시간 토큰 및 비용 계량기**
   - 입출력 토큰과 예상 비용을 실시간으로 계산하여 표시합니다.

---

## 🧠 Pi 표준 7단계 추론 깊이 선택기

| 레벨 | UI 표기 | 토큰 예산 | 주 사용처 |
| :--- | :--- | :--- | :--- |
| `off` | **추론 끄기** | 0 토큰 | 단순 질의응답, 포맷 변환 |
| `minimal` | **최소 추론** | ~1k 토큰 | 한 줄 수정, 빠른 검색 |
| `low` | **낮은 추론** | ~2k 토큰 | 짧은 스크립트 작성, 가벼운 리팩토링 |
| `medium` | **중간 추론** | ~8k 토큰 | 일반적인 기능 구현 및 디버깅 |
| `high` | **깊은 추론** | ~16k 토큰 | 복잡한 알고리즘 설계 및 테스트 루프 (권장) |
| `xhigh` | **매우 깊은 추론**| ~32k 토큰 | 대규모 아키텍처 분석 및 모듈 재설계 |
| `max` | **최대 추론** | 최대 한도 | 심각한 버그의 근본 원인 규명 및 극단 추론 |

*추론을 지원하지 않는 모델(예: GPT-4o)을 선택하면 버튼이 자동으로 비활성화됩니다.*

---

## 🚀 빠른 시작 가이드

### 사전 요구 사항
- [Go 1.22+](https://golang.org/dl/) (소스 빌드 시에만 필요)
- Git (선택 사항)

### 1. 소스 코드에서 실행
```bash
git clone https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow.git
cd Autonomous-Agent-Development-Workflow

# 에이전트 서버 시작 (기본 포트: 8080)
go run cmd/agent/main.go -port 8080 -dir .
```

브라우저에서 `http://127.0.0.1:8080` 접속.

### 2. 단일 바이너리 빌드
```bash
# Linux / macOS
go build -ldflags "-s -w" -o agent ./cmd/agent

# Windows
go build -ldflags "-s -w" -o agent.exe ./cmd/agent
```

---

## ⚙️ 설정 가이드

`config.example.json`을 `config.json`으로 복사하여 사용합니다:

```bash
cp config.example.json config.json
```

설정 예시:
```json
{
  "active_profile_id": "deepseek-flash",
  "profiles": [
    {
      "id": "deepseek-flash",
      "name": "DeepSeek V4.1 Flash",
      "provider": "deepseek",
      "base_url": "https://api.deepseek.com",
      "api_key": "YOUR_API_KEY",
      "model": "deepseek-flash",
      "max_steps": 15
    }
  ],
  "workspace_dir": ".",
  "approval_mode": "guarded"
}
```

---

## 🧪 테스트 및 품질 검증

```bash
go test -v ./...
```

모든 단위 테스트, 통합 테스트 및 보안 가드레일 테스트 100% 통과를 보장합니다.

---

## 👤 작성자 및 라이선스

- **작성자**: [yifangdong](https://github.com/yfddxwx)
- **저장소**: [Autonomous-Agent-Development-Workflow](https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow)
- **라이선스**: [MIT License](LICENSE) © 2026 yifangdong

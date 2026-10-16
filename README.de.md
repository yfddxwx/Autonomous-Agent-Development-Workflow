<p align="center">
  <h1 align="center">Autonomous Agent Development Workflow</h1>
  <p align="center">
    <strong>Ein weltklasse, autonomer KI-Software-Engineering-Agent als schlankes Single-Binary in reinem Go</strong>
  </p>
  <p align="center">
    Null externe Laufzeitabhängigkeiten • ReAct Autonome Schleife & Selbstreparatur • 43+ Modell-Provider • 7-stufige Denktiefe-Steuerung • 3-stufige HITL-Sicherheitsgrenzen • Live-Terminal-Streaming • Git-Rollback
  </p>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/Lizenz-MIT-blue.svg" alt="Lizenz: MIT"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/Build-Erfolgreich-brightgreen.svg" alt="Build Status"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/Tests-100%25%20Bestanden-success.svg" alt="Test Status"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/pulls"><img src="https://img.shields.io/badge/PRs-Willkommen-brightgreen.svg" alt="PRs Welcome"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow"><img src="https://img.shields.io/badge/Provider-43%2B%20Unterst%C3%BCtzt-purple.svg" alt="Providers"></a>
</p>

---

<p align="center">
  <strong>Language Switcher / Sprachauswahl</strong><br>
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

## 📖 Inhaltsverzeichnis

- [Überblick](#-überblick)
- [Hauptmerkmale](#-hauptmerkmale)
- [Systemarchitektur](#-systemarchitektur)
- [Pi-konforme 7-stufige Denktiefe-Steuerung](#-pi-konforme-7-stufige-denktiefe-steuerung)
- [Human-in-the-Loop (HITL) Sicherheit](#-human-in-the-loop-hitl-sicherheit)
- [43+ Unterstützte Modell-Plattformen](#-43-unterstützte-modell-plattformen)
- [Schnellstart](#-schnellstart)
  - [Voraussetzungen](#voraussetzungen)
  - [Aus dem Quellcode ausführen](#1-aus-dem-quellcode-ausführen)
  - [Single-Binary kompilieren](#2-single-binary-kompilieren)
- [Konfigurationsanleitung](#-konfigurationsanleitung)
- [Interaktive Weboberfläche](#-interaktive-weboberfläche)
- [Testen & Qualitätssicherung](#-testen--qualitätssicherung)
- [Mitwirken (Contributing)](#-mitwirken-contributing)
- [Autor & Lizenz](#-autor--lizenz)

---

## 🌟 Überblick

**Autonomous Agent Development Workflow** ist ein extrem leichtgewichtiger, performanter, eigenständiger KI-Entwicklungsagent, der von Grund auf in reiner Go-Sprache entwickelt wurde. Er fungiert als lokaler autonomer Pair-Programming-Partner, der Entwicklungsaufgaben plant, Code liest und modifiziert, Shell-Befehle sicher ausführt, Kompilierfehler autonom behebt und den Verlauf in Verzweigungsbäumen sichert.

Im Gegensatz zu monolithischen Python-Frameworks mit gigantischen virtuellen Umgebungen kompiliert dieses Projekt zu einer **einzelnen ~9MB Binärdatei**, die ihre moderne Google Material 3 Web-UI direkt über Go's `//go:embed` beinhaltet.

---

## ✨ Hauptmerkmale

1. **Keine externen Laufzeitabhängigkeiten (`//go:embed`)**
   - HTML, Tailwind CSS, Offline-Markdown-Parser und Icons sind vollständig in die Binärdatei eingebettet. Kein Node.js oder Python erforderlich.
2. **Autonome ReAct-Schleife & Selbstreparatur**
   - Setzt den Zyklus `Think -> Act -> Observe -> Reflect` um. Bei Test- oder Build-Fehlern analysiert der Agent die Logs und passt den Code selbstständig an.
3. **Pi-konforme 7-stufige Denktiefe-Steuerung (Thinking Depth)**
   - Wählbare Denkstufen direkt an der Eingabeleiste: `Aus (off)`, `Minimal (minimal)`, `Niedrig (low)`, `Mittel (medium)`, `Tief (high)`, `Extra-Tief (xhigh)`, `Maximal (max)`.
   - Erkennt automatisch die Fähigkeiten des Modells (DeepSeek, Xiaomi MiMo, Claude 3.7, OpenAI o-Serie, Gemini 3).
4. **3-stufige HITL-Sicherheitsgrenzen (Antigravity Guardrails)**
   - **Auto**: Volle Autonomie für Sandboxes.
   - **Guarded**: Sichere Lese- und Diagnosebefehle laufen automatisch; Dateiüberschreibungen, Löschungen und zerstörerische Befehle erfordern 1-Klick-Zustimmung.
   - **Manual**: Jeder einzelne Werkzeugaufruf erfordert eine Bestätigung.
5. **Live-Terminal-Streaming**
   - Befehlsausgaben werden über RFC 6455 WebSockets zeilenweise in Echtzeit gestreamt.
6. **Git-Aware Rollback**
   - Erkennt Repository-Änderungen und ermöglicht das Zurücksetzen mit einem einzigen Klick.
7. **Nicht-destruktiver Session-Verzweigungsbaum (Time-Travel)**
   - Erstellt nach jedem Durchlauf Schnappschüsse, sodass Sie jederzeit zu früheren Diskussionsknoten zurückkehren können.
8. **Echtzeit Token- & Kostenzähler**
   - Misst Prompt- und Completion-Tokens und berechnet geschätzte Kosten laufend mit.

---

## 🧠 Pi-konforme 7-stufige Denktiefe-Steuerung

| Stufe | Anzeige | Token-Budget | Anwendungsbereich |
| :--- | :--- | :--- | :--- |
| `off` | **Denken aus** | 0 Tokens | Direkte Antworten, Formatierungen |
| `minimal` | **Minimal** | ~1k Tokens | Einzeilige Korrekturen, schnelle Lookups |
| `low` | **Niedrig** | ~2k Tokens | Kleine Skripte, leichte Refactorings |
| `medium` | **Mittel** | ~8k Tokens | Standard-Feature-Implementierungen |
| `high` | **Tief** | ~16k Tokens | Komplexe Algorithmen & Test-Reparaturschleifen (Empfohlen) |
| `xhigh` | **Extra-Tief** | ~32k Tokens | Umfangreiche Architekturanalysen |
| `max` | **Maximal** | Max. Kapazität | Tiefgehende Ursachenforschung schwerer Fehler |

---

## 🚀 Schnellstart

### Voraussetzungen
- [Go 1.22+](https://golang.org/dl/) (nur für das Kompilieren aus dem Quellcode)
- Git (empfohlen)

### 1. Aus dem Quellcode ausführen
```bash
git clone https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow.git
cd Autonomous-Agent-Development-Workflow

# Server starten (Standardport: 8080)
go run cmd/agent/main.go -port 8080 -dir .
```

Im Browser öffnen: `http://127.0.0.1:8080`

### 2. Single-Binary kompilieren
```bash
# Linux / macOS
go build -ldflags "-s -w" -o agent ./cmd/agent

# Windows
go build -ldflags "-s -w" -o agent.exe ./cmd/agent
```

---

## ⚙️ Konfigurationsanleitung

Kopieren Sie `config.example.json` nach `config.json`:

```bash
cp config.example.json config.json
```

Beispiel:
```json
{
  "active_profile_id": "deepseek-flash",
  "profiles": [
    {
      "id": "deepseek-flash",
      "name": "DeepSeek V4.1 Flash",
      "provider": "deepseek",
      "base_url": "https://api.deepseek.com",
      "api_key": "IHR_API_KEY",
      "model": "deepseek-flash",
      "max_steps": 15
    }
  ],
  "workspace_dir": ".",
  "approval_mode": "guarded"
}
```

---

## 🧪 Testen & Qualitätssicherung

```bash
go test -v ./...
```

Alle Unit-, Integrations- und Sicherheitstests laufen zu 100% fehlerfrei durch.

---

## 👤 Autor & Lizenz

- **Autor**: [yifangdong](https://github.com/yfddxwx)
- **Repository**: [Autonomous-Agent-Development-Workflow](https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow)
- **Lizenz**: [MIT License](LICENSE) © 2026 yifangdong

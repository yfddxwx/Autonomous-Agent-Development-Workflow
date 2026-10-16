<p align="center">
  <h1 align="center">Autonomous Agent Development Workflow</h1>
  <p align="center">
    <strong>Um agente autônomo de engenharia de software de IA de classe mundial, binário único, em Go puro</strong>
  </p>
  <p align="center">
    Zero dependências de tempo de execução externas • Loop autônomo ReAct • 43+ provedores de modelos • 7 níveis de profundidade de raciocínio • Guardrails HITL em 3 camadas • Streaming de terminal em tempo real • Rollback via Git
  </p>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Versão do Go"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="Licença: MIT"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/Build-Passing-brightgreen.svg" alt="Status do Build"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/Tests-100%25%20Passing-success.svg" alt="Status dos Testes"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/pulls"><img src="https://img.shields.io/badge/PRs-Welcome-brightgreen.svg" alt="PRs Bem-vindos"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow"><img src="https://img.shields.io/badge/Providers-43%2B%20Suportados-purple.svg" alt="Provedores"></a>
</p>

---

<p align="center">
  <strong>Alternador de Idioma / Language Switcher</strong><br>
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

## 📖 Sumário

- [Visão Geral](#-visão-geral)
- [Principais Funcionalidades](#-principais-funcionalidades)
- [Arquitetura do Sistema](#-arquitetura-do-sistema)
- [Seletor de Profundidade de Raciocínio Alinhado ao Pi](#-seletor-de-profundidade-de-raciocínio-alinhado-ao-pi)
- [Segurança com Intervenção Humana (HITL)](#-segurança-com-intervenção-humana-hitl)
- [Registro Multi-Provedores (43+ Plataformas)](#-registro-multi-provedores-43-plataformas)
- [Início Rápido](#-início-rápido)
  - [Pré-requisitos](#pré-requisitos)
  - [Executar a partir do Código-Fonte](#1-executar-a-partir-do-código-fonte)
  - [Compilar em Binário Único](#2-compilar-em-binário-único)
- [Guia de Configuração](#-guia-de-configuração)
- [Interface Web Interativa](#-interface-web-interativa)
- [Testes e Garantia de Qualidade](#-testes-e-garantia-de-qualidade)
- [Contribuição](#-contribuição)
- [Autor e Licença](#-autor-e-licença)

---

## 🌟 Visão Geral

O **Autonomous Agent Development Workflow** é um agente autônomo de desenvolvimento de software com IA ultraleve, de alto desempenho e empacotado em um único binário, construído do zero em linguagem Go. Ele fornece a desenvolvedores individuais e equipes de engenharia um parceiro autônomo de programação em par, capaz de planejar tarefas, ler e editar código-fonte, executar comandos no terminal, autocorrigir falhas de compilação e explorar caminhos alternativos de arquitetura.

Diferente de frameworks pesados em Python que exigem gigabytes de ambientes virtuais e dependências node, este agente é compilado em um **binário único de apenas ~9MB**, integrando toda a sua interface Web moderna inspirada no Google Material 3 e Gemini via a diretiva nativa `//go:embed`.

---

## ✨ Principais Funcionalidades

1. **Binário Único Sem Dependências Externas (`//go:embed`)**
   - Todo o frontend (HTML, Tailwind CSS, renderizador Markdown offline, ícones) é embutido diretamente no executável.
   - Execute instantaneamente no Windows, Linux ou macOS sem nenhuma etapa de instalação.

2. **Loop Autônomo ReAct com Autocura**
   - Implementa a metodologia `Think -> Act -> Observe -> Reflect`.
   - Quando suítes de testes ou compilações falham, o agente intercepta stdout/stderr, diagnostica a causa raiz e refatora o código autonomamente até obter sucesso.

3. **Controle de Raciocínio em 7 Níveis Alinhado ao Pi**
   - Seletor de esforço de raciocínio integrado na barra de prompt: `off`, `minimal`, `low`, `medium`, `high`, `xhigh` e `max`.
   - Adaptação dinâmica de acordo com o suporte de cada modelo (ex.: DeepSeek V4.1 Flash, Xiaomi MiMo, Claude 3.7 Sonnet, OpenAI série-o, Gemini 3).

4. **Guardrails de Segurança HITL em 3 Camadas (Antigravity Guardrails)**
   - **Manual**: Revisão humana obrigatória para cada alteração de arquivo ou comando.
   - **Guarded**: Diagnósticos e leituras seguras executam de forma automática; edições destrutivas, exclusões ou comandos perigosos exigem autorização em um clique.
   - **Auto**: Autonomia total recomendada para ambientes em sandbox ou contêineres isolados.

5. **Streaming de Terminal em Tempo Real**
   - Comandos de longa duração transmitem suas saídas linha por linha via WebSockets RFC 6455 para o console escuro integrado.

6. **Rollback Atômico com Suporte a Git**
   - Monitoramento contínuo do status do repositório, permitindo desfazer alterações e limpar arquivos não rastreados com um único clique.

7. **Árvore de Sessões e Viagem no Tempo (Time-Travel)**
   - Snapshots automáticos a cada etapa para alternar entre ramificações sem perder histórico de raciocínio.

8. **Telemetria de Tokens e Cálculo de Custos**
   - Medição ao vivo do consumo de tokens de entrada e saída com estimativa de custo financeiro.

---

## 🏗️ Arquitetura do Sistema

```text
┌────────────────────────────────────────────────────────────────────────┐
│               Interface Web Embutida Google Material 3                 │
│      (Barra de Prompt, Seletor Pi, Árvore de Sessão, Terminal Live)    │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ RFC 6455 WebSockets / REST
┌───────────────────────────────────▼────────────────────────────────────┐
│                    DevAgent Go Core (agent.exe)                        │
│                                                                        │
│   ┌───────────────────────────┐    ┌───────────────────────────────┐   │
│   │     Motor ReAct Core      │◄───┤  Hub Multi-Provedores (43+)   │   │
│   │  (Think-Act-Observe-Reflect)   │    (DeepSeek, Xiaomi, Gemini, │   │
│   └─────────────┬─────────────┘    │     Claude, OpenAI, Ollama)   │   │
│                 │                  └───────────────────────────────┘   │
│                 ▼                                                      │
│   ┌───────────────────────────┐    ┌───────────────────────────────┐   │
│   │  Guardrails de Segurança  │    │     Armazenamento de Sessões  │   │
│   │  (Jail de Caminho, Regex) │    │     (Snapshots e Ramificações)│   │
│   └─────────────┬─────────────┘    └───────────────────────────────┘   │
│                 │                                                      │
│                 ▼                                                      │
│   ┌────────────────────────────────────────────────────────────┐       │
│   │               Ecossistema de Ferramentas Nativas           │       │
│   │  file_read • file_write • file_edit • list_dir • cmd_exec  │       │
│   └────────────────────────────────────────────────────────────┘       │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                         ┌──────────▼──────────┐
                         │ Sistema de Arquivos │
                         │     do Host & Git   │
                         └─────────────────────┘
```

---

## 🧠 Seletor de Profundidade de Raciocínio Alinhado ao Pi

A barra de prompt apresenta um seletor de esforço de raciocínio calibrado conforme as especificações oficiais do framework Pi Agent:

| Nível | Nome na Interface | Orçamento de Tokens | Caso de Uso Típico |
| :--- | :--- | :--- | :--- |
| `off` | **Desativado (Off)** | 0 Tokens | Respostas diretas, formatação de dados, bate-papo |
| `minimal` | **Mínimo (Minimal)** | ~1k Tokens | Resumos rápidos, consultas diretas em arquivos |
| `low` | **Baixo (Low)** | ~2k Tokens | Pequenas refatorações, geração de scripts simples |
| `medium` | **Médio (Medium)** | ~8k Tokens | Implementação de funcionalidades padrão e depuração |
| `high` | **Alto (High)** | ~16k Tokens | Algoritmos complexos, ciclos de testes unitários |
| `xhigh` | **Extra-Alto (Extra-High)**| ~32k Tokens | Síntese de arquitetura, refatoração entre múltiplos módulos |
| `max` | **Máximo (Max)** | Capacidade Máxima | Demonstrações exaustivas e análise de bugs críticos |

*A interface expõe dinamicamente apenas os níveis suportados pelo modelo ativo.*

---

## 🛡️ Segurança com Intervenção Humana (HITL)

- **Confinamento Estrito de Caminho**: Bloqueia tentativas de escape de diretório (`../`), ancorando as operações na pasta do projeto.
- **Interceptação de Comandos Perigosos**: Padrões com potencial destrutivo (`rm -rf /`, formatação de discos, vazamento de credenciais) são interceptados ou exigem aprovação explícita.
- **Diff Cirúrgico Visual**: Alterações em arquivos exibem substituições precisas destacadas em vermelho e verde antes da gravação em disco.

---

## 🌐 Registro Multi-Provedores (43+ Plataformas)

Suporte nativo com URLs base verificadas e matriz de modelos recomendados:

- **Provedores Globais em Nuvem**: DeepSeek, Xiaomi (MiMo), Google Gemini, Anthropic Claude, OpenAI, xAI (Grok), Mistral, Z.AI (GLM), Moonshot/Kimi.
- **Gateways Especializados em Código**: Kimi for Coding, Z.AI Coding CN, Qwen Token Plan, GitHub Copilot.
- **Chips de Inferência de Alta Velocidade**: Groq, Cerebras, Fireworks, Together AI, Baseten.
- **Gateways Cloud & Edge**: Cloudflare AI Gateway, Vercel AI Gateway, Azure OpenAI, AWS Bedrock.
- **Modelos Locais Privados**: Ollama, llama.cpp, LocalAI, vLLM.

---

## 🚀 Início Rápido

### Pré-requisitos
- [Go 1.22+](https://golang.org/dl/) (necessário apenas se for compilar do código-fonte)
- Git (opcional, para integração de controle de versão)

### 1. Executar a partir do Código-Fonte
```bash
# Clonar o repositório
git clone https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow.git
cd Autonomous-Agent-Development-Workflow

# Iniciar o servidor DevAgent (porta padrão: 8080)
go run cmd/agent/main.go -port 8080 -dir .
```

Acesse o navegador em:
```
http://127.0.0.1:8080
```

### 2. Compilar em Binário Único
Gere um binário enxuto e pronto para produção:
```bash
# No Linux / macOS
go build -ldflags "-s -w" -o agent ./cmd/agent

# No Windows
go build -ldflags "-s -w" -o agent.exe ./cmd/agent
```

---

## ⚙️ Guia de Configuração

Copie o arquivo `config.example.json` para `config.json` (ou configure diretamente pela interface de configurações):

```bash
cp config.example.json config.json
```

Exemplo de perfil:
```json
{
  "active_profile_id": "deepseek-flash",
  "profiles": [
    {
      "id": "deepseek-flash",
      "name": "DeepSeek V4.1 Flash",
      "provider": "deepseek",
      "base_url": "https://api.deepseek.com",
      "api_key": "SUA_CHAVE_API_AQUI",
      "model": "deepseek-flash",
      "max_steps": 15
    }
  ],
  "workspace_dir": ".",
  "approval_mode": "guarded"
}
```

*Variáveis de ambiente (`DEEPSEEK_API_KEY`, `OPENAI_API_KEY`, `GEMINI_API_KEY`, etc.) são detectadas automaticamente na inicialização.*

---

## 🖥️ Interface Web Interativa

- **Design Google Gemini**: Tipografia Google Sans, fontes de código JetBrains Mono e visual escuro moderno (`#131314`).
- **Cards de Raciocínio Expansíveis**: Exibição dos pensamentos do modelo durante o raciocínio, com recolhimento suave ao finalizar.
- **Layout Responsivo**: Barra superior adaptável a diferentes resoluções sem quebras ou sobreposições.
- **Seletor Nativo de Pastas**: Abertura da janela nativa do sistema operacional para selecionar o diretório do projeto.

---

## 🧪 Testes e Garantia de Qualidade

Todas as contribuições e lançamentos passam por quatro níveis de verificação rigorosa:

```bash
# Executar toda a suíte de testes unitários e de integração
go test -v ./...
```

Pacotes de testes verificados:
- `pkg/security`: Confinamento de diretórios, bloqueio de comandos perigosos e despacho de aprovações.
- `pkg/tools`: Substituição de blocos cirúrgicos, listagem de pastas e streaming de comandos.
- `pkg/server`: Troca dinâmica de modelos, validação dos níveis de raciocínio e rotas da API HTTP.
- `pkg/git`: Segurança de reversão e monitoramento do estado do repositório.
- `tests`: Simulação autônoma de ponta a ponta do ciclo ReAct.

---

## 🤝 Contribuição

Contribuições são muito bem-vindas! Consulte o guia [CONTRIBUTING.md](CONTRIBUTING.md) e o [Código de Conduta](CODE_OF_CONDUCT.md) antes de enviar uma solicitação de pull request.

---

## 👤 Autor e Licença

- **Autor**: [yifangdong](https://github.com/yfddxwx)
- **Repositório**: [Autonomous-Agent-Development-Workflow](https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow)
- **Licença**: [Licença MIT](LICENSE) © 2026 yifangdong

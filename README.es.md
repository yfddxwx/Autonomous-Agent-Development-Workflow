<p align="center">
  <h1 align="center">Autonomous Agent Development Workflow</h1>
  <p align="center">
    <strong>Un agente autónomo de ingeniería de software con IA de clase mundial, binario único, en Go puro</strong>
  </p>
  <p align="center">
    Cero dependencias externas en tiempo de ejecución • Bucle autónomo ReAct • 43+ proveedores de modelos • 7 niveles de profundidad de razonamiento • Guardarraíles de seguridad HITL de 3 niveles • Transmisión de terminal en tiempo real • Reversión con Git
  </p>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Versión de Go"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="Licencia: MIT"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/Build-Passing-brightgreen.svg" alt="Estado de Compilación"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/Tests-100%25%20Passing-success.svg" alt="Estado de Pruebas"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/pulls"><img src="https://img.shields.io/badge/PRs-Welcome-brightgreen.svg" alt="PRs Bienvenidos"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow"><img src="https://img.shields.io/badge/Providers-43%2B%20Soportados-purple.svg" alt="Proveedores"></a>
</p>

---

<p align="center">
  <strong>Selector de Idioma / Language Switcher</strong><br>
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

## 📖 Tabla de Contenidos

- [Descripción General](#-descripción-general)
- [Características Principales](#-características-principales)
- [Arquitectura del Sistema](#-arquitectura-del-sistema)
- [Selector de Profundidad de Razonamiento Alineado con Pi](#-selector-de-profundidad-de-razonamiento-alineado-con-pi)
- [Seguridad con Intervención Humana (HITL)](#-seguridad-con-intervención-humana-hitl)
- [Registro de Múltiples Proveedores (43+ Plataformas)](#-registro-de-múltiples-proveedores-43-plataformas)
- [Inicio Rápido](#-inicio-rápido)
  - [Requisitos Previos](#requisitos-previos)
  - [Ejecutar desde el Código Fuente](#1-ejecutar-desde-el-código-fuente)
  - [Compilar Binario Único](#2-compilar-binario-único)
- [Guía de Configuración](#-guía-de-configuración)
- [Interfaz Web Interactiva](#-interfaz-web-interactiva)
- [Pruebas y Control de Calidad](#-pruebas-y-control-de-calidad)
- [Contribución](#-contribución)
- [Autor y Licencia](#-autor-y-licencia)

---

## 🌟 Descripción General

**Autonomous Agent Development Workflow** es un agente autónomo de desarrollo de software con IA ultraligero, de alto rendimiento y empaquetado en un único archivo binario, construido íntegramente en lenguaje Go. Proporciona a desarrolladores y equipos de ingeniería un compañero de programación en pareja capaz de planificar tareas, examinar y modificar código fuente, ejecutar comandos de terminal, autocorregir fallos de compilación y evaluar vías arquitectónicas alternativas.

A diferencia de los entornos pesados de Python que demandan gigabytes de dependencias virtuales y módulos npm, este agente se compila en un **único binario de aproximadamente 9MB**, incorporando directamente su moderna interfaz Web inspirada en Google Material 3 y Gemini mediante la directiva `//go:embed`.

---

## ✨ Características Principales

1. **Binario Único sin Dependencias Externas (`//go:embed`)**
   - Todo el frontend (HTML, Tailwind CSS, motor de renderizado Markdown sin conexión, iconos) está integrado en el ejecutable.
   - Funciona de inmediato en Windows, Linux y macOS sin necesidad de configurar entornos adicionales.

2. **Bucle Autónomo ReAct con Autorrecuperación**
   - Aplica el ciclo `Think -> Act -> Observe -> Reflect`.
   - Cuando las pruebas unitarias o las compilaciones fallan, el agente captura la salida estándar y de errores, diagnostica la raíz del problema y refactoriza el código de manera autónoma hasta resolverlo.

3. **Control de Razonamiento en 7 Niveles Alineado con Pi**
   - Selector visual de intensidad de pensamiento integrado en la barra de entrada: `off`, `minimal`, `low`, `medium`, `high`, `xhigh` y `max`.
   - Adaptación dinámica de niveles disponibles según las capacidades del modelo (p. ej., DeepSeek V4.1 Flash, Xiaomi MiMo, Claude 3.7 Sonnet, serie OpenAI o, Gemini 3).

4. **Guardarraíles de Seguridad HITL de 3 Niveles (Antigravity Guardrails)**
   - **Manual**: Revisión y aprobación obligatoria para cada edición de archivo o ejecución de comando.
   - **Guarded**: Las operaciones seguras de solo lectura y diagnóstico se ejecutan automáticamente; las escrituras destructivas o comandos de riesgo solicitan autorización en un solo clic.
   - **Auto**: Autonomía total para entornos aislados o contenedores de prueba seguros.

5. **Transmisión de Terminal en Tiempo Real**
   - Los comandos de larga duración envían su salida línea a línea a través de WebSockets RFC 6455 hacia la terminal integrada.

6. **Reversión Atómica con Git**
   - Supervisión en tiempo real del estado del repositorio con función de reversión en un clic para descartar cambios y limpiar archivos no deseados con total seguridad.

7. **Árbol de Sesiones y Viaje en el Tiempo (Time-Travel)**
   - Creación automática de instantáneas (snapshots) en cada paso para alternar entre diferentes enfoques sin perder el historial de razonamiento.

8. **Telemetría de Tokens y Medidor de Costos**
   - Seguimiento exacto de tokens de entrada y salida con estimación económica en vivo.

---

## 🏗️ Arquitectura del Sistema

```text
┌────────────────────────────────────────────────────────────────────────┐
│               Interfaz Web Integrada Google Material 3                 │
│    (Barra de Prompt, Selector Pi, Árbol de Sesión, Terminal en Vivo)   │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ RFC 6455 WebSockets / REST
┌───────────────────────────────────▼────────────────────────────────────┐
│                    DevAgent Go Core (agent.exe)                        │
│                                                                        │
│   ┌───────────────────────────┐    ┌───────────────────────────────┐   │
│   │     Motor ReAct Core      │◄───┤  Hub Multiproveedor (43+)     │   │
│   │  (Think-Act-Observe-Reflect)   │    (DeepSeek, Xiaomi, Gemini, │   │
│   └─────────────┬─────────────┘    │     Claude, OpenAI, Ollama)   │   │
│                 │                  └───────────────────────────────┘   │
│                 ▼                                                      │
│   ┌───────────────────────────┐    ┌───────────────────────────────┐   │
│   │  Guardarraíles HITL       │    │     Gestor de Sesiones        │   │
│   │  (Aislamiento de rutas)   │    │     (Instantáneas y ramas)    │   │
│   └─────────────┬─────────────┘    └───────────────────────────────┘   │
│                 │                                                      │
│                 ▼                                                      │
│   ┌────────────────────────────────────────────────────────────┐       │
│   │                Ecosistema de Herramientas Nativas          │       │
│   │  file_read • file_write • file_edit • list_dir • cmd_exec  │       │
│   └────────────────────────────────────────────────────────────┘       │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                         ┌──────────▼──────────┐
                         │ Sistema de Archivos │
                         │     y Shell/Git     │
                         └─────────────────────┘
```

---

## 🧠 Selector de Profundidad de Razonamiento Alineado con Pi

La barra de comandos incorpora un control de esfuerzo de razonamiento ajustado a las especificaciones oficiales del ecosistema Pi Agent:

| Nivel | Nombre en Interfaz | Presupuesto de Tokens | Caso de Uso Típico |
| :--- | :--- | :--- | :--- |
| `off` | **Desactivado (Off)** | 0 Tokens | Respuestas directas, formateo de texto, conversación |
| `minimal` | **Mínimo (Minimal)** | ~1k Tokens | Resúmenes breves, búsquedas puntuales en código |
| `low` | **Bajo (Low)** | ~2k Tokens | Modificaciones menores, generación de scripts rápidos |
| `medium` | **Medio (Medium)** | ~8k Tokens | Desarrollo estándar de funciones y depuración |
| `high` | **Alto (High)** | ~16k Tokens | Diseño algorítmico avanzado, bucles de pruebas completas |
| `xhigh` | **Extra-Alto (Extra-High)**| ~32k Tokens | Síntesis arquitectónica, refactorización modular compleja |
| `max` | **Máximo (Max)** | Capacidad Máxima | Deducción matemática exhaustiva y resolución de fallos críticos |

*La interfaz visualiza dinámicamente únicamente los niveles compatibles con el modelo en uso.*

---

## 🛡️ Seguridad con Intervención Humana (HITL)

- **Aislamiento Estricto de Rutas**: Neutraliza cualquier intento de escape de directorios (`../`), bloqueando la ejecución al espacio de trabajo asignado.
- **Intercepción de Comandos Críticos**: Patrones potencialmente destructivos (`rm -rf /`, formateo de unidades, lectura de credenciales) son retenidos o elevados para autorización explícita.
- **Comparación Visual de Diferencias (Diff)**: Las sustituciones quirúrgicas de código se visualizan con claridad en rojo y verde antes de ser aplicadas en el disco.

---

## 🌐 Registro de Múltiples Proveedores (43+ Plataformas)

Soporte integrado con URLs base verificadas y matrices de modelos recomendados:

- **Plataformas Cloud Principales**: DeepSeek, Xiaomi (MiMo), Google Gemini, Anthropic Claude, OpenAI, xAI (Grok), Mistral, Z.AI (GLM), Moonshot/Kimi.
- **Pasarelas de Programación Dedicadas**: Kimi for Coding, Z.AI Coding CN, Qwen Token Plan, GitHub Copilot.
- **Aceleradores de Inferencia de Alta Velocidad**: Groq, Cerebras, Fireworks, Together AI, Baseten.
- **Pasarelas Cloud & Edge**: Cloudflare AI Gateway, Vercel AI Gateway, Azure OpenAI, AWS Bedrock.
- **Modelos Privados Locales**: Ollama, llama.cpp, LocalAI, vLLM.

---

## 🚀 Inicio Rápido

### Requisitos Previos
- [Go 1.22+](https://golang.org/dl/) (necesario únicamente para compilar desde las fuentes)
- Git (opcional, para funciones de control de versiones)

### 1. Ejecutar desde el Código Fuente
```bash
# Clonar el repositorio
git clone https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow.git
cd Autonomous-Agent-Development-Workflow

# Iniciar el servidor DevAgent (puerto por defecto: 8080)
go run cmd/agent/main.go -port 8080 -dir .
```

Abra su navegador en:
```
http://127.0.0.1:8080
```

### 2. Compilar Binario Único
Genere un binario optimizado para producción:
```bash
# En Linux / macOS
go build -ldflags "-s -w" -o agent ./cmd/agent

# En Windows
go build -ldflags "-s -w" -o agent.exe ./cmd/agent
```

---

## ⚙️ Guía de Configuración

Copie `config.example.json` a `config.json` (o defina los parámetros desde el panel de ajustes):

```bash
cp config.example.json config.json
```

Ejemplo de definición de perfil:
```json
{
  "active_profile_id": "deepseek-flash",
  "profiles": [
    {
      "id": "deepseek-flash",
      "name": "DeepSeek V4.1 Flash",
      "provider": "deepseek",
      "base_url": "https://api.deepseek.com",
      "api_key": "SU_CLAVE_API_AQUI",
      "model": "deepseek-flash",
      "max_steps": 15
    }
  ],
  "workspace_dir": ".",
  "approval_mode": "guarded"
}
```

*Las variables de entorno (`DEEPSEEK_API_KEY`, `OPENAI_API_KEY`, `GEMINI_API_KEY`, etc.) se importan automáticamente al iniciar la aplicación.*

---

## 🖥️ Interfaz Web Interactiva

- **Estética Google Gemini**: Fuentes tipográficas Google Sans y JetBrains Mono sobre paneles oscuros ergonómicos (`#131314`).
- **Bloques Desplegables de Pensamiento**: El proceso de razonamiento se despliega en tiempo real y se contrae elegantemente al finalizar.
- **Diseño Totalmente Adaptable**: Barra superior optimizada para cualquier pantalla sin recortes.
- **Selector Nativo de Carpetas**: Acceso directo al cuadro de diálogo nativo del sistema operativo para seleccionar carpetas de proyecto.

---

## 🧪 Pruebas y Control de Calidad

Cada solicitud de extracción (PR) y lanzamiento se evalúa rigurosamente frente a cuatro barreras de calidad:

```bash
# Ejecutar todas las pruebas unitarias y de integración
go test -v ./...
```

Suites de verificación integradas:
- `pkg/security`: Confinamiento de rutas, intercepción de comandos de riesgo y gestión de autorizaciones.
- `pkg/tools`: Reemplazo quirúrgico de bloques, navegación de directorios y streaming de comandos.
- `pkg/server`: Conmutación dinámica de modelos, validación de profundidad de razonamiento y endpoints HTTP.
- `pkg/git`: Seguridad de rollback y seguimiento del estado del repositorio.
- `tests`: Simulación completa de extremo a extremo del bucle ReAct.

---

## 🤝 Contribución

¡Todas las contribuciones son bienvenidas! Por favor, revise el documento [CONTRIBUTING.md](CONTRIBUTING.md) y el [Código de Conducta](CODE_OF_CONDUCT.md) antes de enviar una solicitud de cambio.

---

## 👤 Autor y Licencia

- **Autor**: [yifangdong](https://github.com/yfddxwx)
- **Repositorio**: [Autonomous-Agent-Development-Workflow](https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow)
- **Licencia**: [Licencia MIT](LICENSE) © 2026 yifangdong

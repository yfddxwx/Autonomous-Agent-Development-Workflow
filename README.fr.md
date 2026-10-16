<p align="center">
  <h1 align="center">Autonomous Agent Development Workflow</h1>
  <p align="center">
    <strong>Un agent autonome d'ingénierie logicielle IA de classe mondiale, binaire unique, en pur Go</strong>
  </p>
  <p align="center">
    Zéro dépendance d'exécution externe • Boucle autonome ReAct • 43+ fournisseurs de modèles • 7 niveaux de profondeur de réflexion • 3 niveaux de garde-fous HITL • Diffusion de terminal en direct • Restauration Git
  </p>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Version Go"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="Licence : MIT"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/Build-Passing-brightgreen.svg" alt="Statut du build"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/actions"><img src="https://img.shields.io/badge/Tests-100%25%20Passing-success.svg" alt="Statut des tests"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow/pulls"><img src="https://img.shields.io/badge/PRs-Welcome-brightgreen.svg" alt="Contributions bienvenues"></a>
  <a href="https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow"><img src="https://img.shields.io/badge/Providers-43%2B%20Supportés-purple.svg" alt="Fournisseurs"></a>
</p>

---

<p align="center">
  <strong>Sélecteur de langue / Language Switcher</strong><br>
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

## 📖 Table des matières

- [Présentation](#-présentation)
- [Fonctionnalités clés](#-fonctionnalités-clés)
- [Architecture du système](#-architecture-du-système)
- [Sélecteur de profondeur de réflexion aligné sur Pi](#-sélecteur-de-profondeur-de-réflexion-aligné-sur-pi)
- [Sécurité avec intervention humaine (HITL)](#-sécurité-avec-intervention-humaine-hitl)
- [Registre multi-fournisseurs (43+ plateformes)](#-registre-multi-fournisseurs-43-plateformes)
- [Démarrage rapide](#-démarrage-rapide)
  - [Prérequis](#prérequis)
  - [Exécution depuis les sources](#1-exécution-depuis-les-sources)
  - [Compilation en binaire unique](#2-compilation-en-binaire-unique)
- [Guide de configuration](#-guide-de-configuration)
- [Interface Web interactive](#-interface-web-interactive)
- [Tests et assurance qualité](#-tests-et-assurance-qualité)
- [Contribution](#-contribution)
- [Auteur et Licence](#-auteur-et-licence)

---

## 🌟 Présentation

**Autonomous Agent Development Workflow** est un agent de développement logiciel IA ultra-léger, haute performance et distribué sous forme de binaire unique, conçu intégralement en langage Go. Il offre aux développeurs et équipes d'ingénierie un partenaire de pair-programming autonome capable de planifier des tâches, d'explorer et modifier du code source, d'exécuter des commandes système, d'auto-corriger les erreurs de compilation et d'explorer des branches d'architecture alternatives.

Contrairement aux frameworks Python encombrants nécessitant des gigaoctets d'environnements virtuels et de modules npm, cet agent se compile en un **binaire unique d'environ 9 Mo**, intégrant directement son interface Web moderne inspirée de Google Material 3 et Gemini grâce à la directive `//go:embed`.

---

## ✨ Fonctionnalités clés

1. **Binaire unique sans dépendance (`//go:embed`)**
   - L'ensemble du frontend (HTML, Tailwind CSS, moteur Markdown hors-ligne, icônes) est embarqué dans un exécutable autonome.
   - S'exécute instantanément sur Windows, Linux ou macOS sans installation complexe de runtime.

2. **Boucle autonome ReAct et auto-guérison**
   - Implémente le paradigme `Think -> Act -> Observe -> Reflect`.
   - Lors de l'échec d'un test ou d'une compilation, l'agent intercepte stdout/stderr, analyse la cause racine et refactorise le code de manière itérative jusqu'au succès.

3. **Contrôle de profondeur de réflexion à 7 niveaux aligné sur Pi**
   - Sélecteur et jauge d'effort de raisonnement intégrés : `off`, `minimal`, `low`, `medium`, `high`, `xhigh` et `max`.
   - Adaptation dynamique en temps réel selon les capacités du modèle (ex. DeepSeek V4.1 Flash, Xiaomi MiMo, Claude 3.7 Sonnet, OpenAI o-series, Gemini 3).

4. **Garde-fous de sécurité HITL à 3 niveaux (Antigravity Guardrails)**
   - **Manual** : Examen et approbation humaine obligatoires pour chaque écriture de fichier ou exécution de commande.
   - **Guarded** : Diagnostics sûrs et lectures automatisées ; les modifications destructives, suppressions ou commandes sensibles requièrent une validation en un clic.
   - **Auto** : Autonomie complète adaptée aux environnements sandbox ou conteneurisés sécurisés.

5. **Streaming de terminal en temps réel**
   - Exécution asynchrone des commandes avec diffusion des sorties ligne par ligne via WebSockets RFC 6455 dans une console intégrée.

6. **Restauration atomique assistée par Git**
   - Surveillance continue du statut du dépôt Git avec retour arrière immédiat en un clic pour annuler les modifications indésirables.

7. **Arbre de sessions et voyage dans le temps (Time-Travel)**
   - Instantanés automatiques à chaque étape d'exécution pour basculer facilement entre différentes hypothèses de conception sans perte d'historique.

8. **Télémétrie des tokens et suivi des coûts**
   - Mesure précise en direct des tokens de prompt et de complétion avec estimation des coûts.

---

## 🏗️ Architecture du système

```text
┌────────────────────────────────────────────────────────────────────────┐
│               Interface Web embarquée Google Material 3                │
│    (Barre de prompt, sélecteur de réflexion, sessions, terminal live)   │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ RFC 6455 WebSockets / REST
┌───────────────────────────────────▼────────────────────────────────────┐
│                    DevAgent Go Core (agent.exe)                        │
│                                                                        │
│   ┌───────────────────────────┐    ┌───────────────────────────────┐   │
│   │    Moteur ReAct Core      │◄───┤  Hub multi-fournisseurs (43+) │   │
│   │  (Think-Act-Observe-Reflect)   │    (DeepSeek, Xiaomi, Gemini, │   │
│   └─────────────┬─────────────┘    │     Claude, OpenAI, Ollama)   │   │
│                 │                  └───────────────────────────────┘   │
│                 ▼                                                      │
│   ┌───────────────────────────┐    ┌───────────────────────────────┐   │
│   │ Garde-fous HITL & Prison  │    │  Gestionnaire d'historique    │   │
│   │  (Jail de chemin, Regex)  │    │  (Snapshots & branches)       │   │
│   └─────────────┬─────────────┘    └───────────────────────────────┘   │
│                 │                                                      │
│                 ▼                                                      │
│   ┌────────────────────────────────────────────────────────────┐       │
│   │                 Écosystème d'outils intégrés               │       │
│   │  file_read • file_write • file_edit • list_dir • cmd_exec  │       │
│   └────────────────────────────────────────────────────────────┘       │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
                         ┌──────────▼──────────┐
                         │ Système de fichiers │
                         │     & Shell/Git     │
                         └─────────────────────┘
```

---

## 🧠 Sélecteur de profondeur de réflexion aligné sur Pi

La barre de saisie intègre un sélecteur de profondeur calibré selon les spécifications officielles de l'écosystème Pi Agent :

| Niveau | Nom de l'interface | Budget de tokens | Cas d'utilisation type |
| :--- | :--- | :--- | :--- |
| `off` | **Désactivé (Off)** | 0 Token | Réponses directes, questions conversationnelles |
| `minimal` | **Minimal** | ~1k Tokens | Résumé concis, consultations ponctuelles |
| `low` | **Faible (Low)** | ~2k Tokens | Petits refactorings, scripts légers |
| `medium` | **Moyen (Medium)** | ~8k Tokens | Implémentation standard de fonctionnalités et débogage |
| `high` | **Élevé (High)** | ~16k Tokens | Conception algorithmique complexe, boucles de tests |
| `xhigh` | **Très élevé (Extra-High)**| ~32k Tokens | Synthèse d'architecture, refactorisation inter-modules |
| `max` | **Maximal (Max)** | Capacité Max | Démonstration formelle et diagnostic d'anomalies critiques |

*L'interface n'affiche dynamiquement que les niveaux pris en charge par le modèle sélectionné.*

---

## 🛡️ Sécurité avec intervention humaine (HITL)

- **Confinement de chemin strict** : Empêche toute tentative d'évasion d'arborescence (`../`), isolant l'agent dans son espace de travail.
- **Interception de commandes à risque** : Détection des modèles destructeurs (`rm -rf /`, modifications disques, accès trousseaux) et soumission à validation expresse.
- **Diff visuel chirurgical** : Aperçu précis en rouge/vert des substitutions de code avant application effective sur le disque.

---

## 🌐 Registre multi-fournisseurs (43+ plateformes)

Prise en charge intégrée avec URLs de base validées et matrices de modèles recommandées :

- **Grands modèles Cloud** : DeepSeek, Xiaomi (MiMo), Google Gemini, Anthropic Claude, OpenAI, xAI (Grok), Mistral, Z.AI (GLM), Moonshot/Kimi.
- **Passerelles dédiées au code** : Kimi for Coding, Z.AI Coding CN, Qwen Token Plan, GitHub Copilot.
- **Accélérateurs matériels haute vitesse** : Groq, Cerebras, Fireworks, Together AI, Baseten.
- **Passerelles Cloud & Edge** : Cloudflare AI Gateway, Vercel AI Gateway, Azure OpenAI, AWS Bedrock.
- **Modèles locaux privés** : Ollama, llama.cpp, LocalAI, vLLM.

---

## 🚀 Démarrage rapide

### Prérequis
- [Go 1.22+](https://golang.org/dl/) (uniquement en cas de compilation source)
- Git (optionnel, pour les fonctions de versionnement)

### 1. Exécution depuis les sources
```bash
# Cloner le dépôt
git clone https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow.git
cd Autonomous-Agent-Development-Workflow

# Lancer le serveur DevAgent (port par défaut : 8080)
go run cmd/agent/main.go -port 8080 -dir .
```

Ouvrez votre navigateur sur :
```
http://127.0.0.1:8080
```

### 2. Compilation en binaire unique
Générez un exécutable allégé et optimisé pour la production :
```bash
# Sur Linux / macOS
go build -ldflags "-s -w" -o agent ./cmd/agent

# Sur Windows
go build -ldflags "-s -w" -o agent.exe ./cmd/agent
```

---

## ⚙️ Guide de configuration

Copiez `config.example.json` vers `config.json` (ou configurez directement via l'interface des paramètres) :

```bash
cp config.example.json config.json
```

Exemple de profil :
```json
{
  "active_profile_id": "deepseek-flash",
  "profiles": [
    {
      "id": "deepseek-flash",
      "name": "DeepSeek V4.1 Flash",
      "provider": "deepseek",
      "base_url": "https://api.deepseek.com",
      "api_key": "VOTRE_CLE_API_ICI",
      "model": "deepseek-flash",
      "max_steps": 15
    }
  ],
  "workspace_dir": ".",
  "approval_mode": "guarded"
}
```

*Les variables d'environnement (`DEEPSEEK_API_KEY`, `OPENAI_API_KEY`, `GEMINI_API_KEY`, etc.) sont automatiquement détectées au démarrage.*

---

## 🖥️ Interface Web interactive

- **Thème Google Gemini** : Typographie Google Sans, police de code JetBrains Mono et surfaces sombres ergonomiques (`#131314`).
- **Volets de réflexion rétractables** : Les étapes d'introspection s'affichent en direct puis se replient harmonieusement.
- **Navigation adaptable** : Barre supérieure réactive s'ajustant à toutes les définitions d'écran.
- **Explorateur de dossiers natif** : Sélection de répertoire via la boîte de dialogue native du système d'exploitation.

---

## 🧪 Tests et assurance qualité

Chaque demande d'intégration (PR) et version est validée à travers une suite rigoureuse de tests :

```bash
# Exécuter tous les tests unitaires et de sécurité
go test -v ./...
```

Suites de tests incluses :
- `pkg/security` : Confinement de chemin, blocage de commandes à risque et gestion des autorisations.
- `pkg/tools` : Modification chirurgicale de blocs, exploration de répertoires et streaming terminal.
- `pkg/server` : Changement dynamique de profil, validation des niveaux de réflexion et routes API HTTP.
- `pkg/git` : Intégrité des rollbacks et surveillance de l'état du dépôt.
- `tests` : Simulation autonome de bout en bout de la boucle ReAct.

---

## 🤝 Contribution

Les contributions sont chaleureusement accueillies ! Veuillez consulter le guide [CONTRIBUTING.md](CONTRIBUTING.md) et le [Code de conduite](CODE_OF_CONDUCT.md) avant de soumettre une pull request.

---

## 👤 Auteur et Licence

- **Auteur** : [yifangdong](https://github.com/yfddxwx)
- **Dépôt** : [Autonomous-Agent-Development-Workflow](https://github.com/yfddxwx/Autonomous-Agent-Development-Workflow)
- **Licence** : [Licence MIT](LICENSE) © 2026 yifangdong

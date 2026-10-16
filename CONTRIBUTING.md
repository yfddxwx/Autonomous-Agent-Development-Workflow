# Contributing to Autonomous Agent Development Workflow

Thank you for your interest in contributing to **Autonomous Agent Development Workflow**! We are committed to building a world-class, high-reliability autonomous software engineering assistant.

## Code of Conduct

This project and everyone participating in it is governed by the [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## How to Contribute

### 1. Reporting Bugs
- Check existing issues before opening a new one.
- Use the [Bug Report Template](.github/ISSUE_TEMPLATE/bug_report.md).
- Include operating system, Go version, model provider in use, reproduction steps, and error logs.

### 2. Suggesting Enhancements
- Use the [Feature Request Template](.github/ISSUE_TEMPLATE/feature_request.md).
- Clearly describe the problem, proposed solution, and alternative considerations.

### 3. Pull Requests
1. Fork the repository and create your branch from `main`.
2. Ensure Go 1.22+ is installed.
3. Keep code modular and adhere to the project's **Four-Fold Gatekeeping Standards**:
   - Functional completeness.
   - 100% unit test coverage for new logic.
   - Integration verification with mock SSE / HTTP server.
   - Security audit (no path traversal, strict workspace isolation).
4. Run all tests locally before submitting:
   ```bash
   go test -v ./...
   ```
5. Follow conventional commits (e.g., `feat:`, `fix:`, `docs:`, `refactor:`, `test:`).
6. Submit a Pull Request referencing relevant issues.

## Development Guidelines

- **Zero External Runtime Dependencies**: The distribution binary must remain single, lightweight, and self-contained using `//go:embed`.
- **Security-First**: Any tool that can modify files or execute shell commands must respect the Human-in-the-Loop (HITL) gate and workspace directory confinement.
- **Documentation**: All public APIs and interfaces must include clear Go doc comments.

Thank you for helping make autonomous coding safer and more productive!

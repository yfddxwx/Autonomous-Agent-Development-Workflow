## Description
Briefly describe the change and rationale behind it.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Performance improvement / Refactoring

## Four-Fold Gatekeeping Checklist
- [ ] My code follows the clean architecture and coding standards of this project.
- [ ] I have added unit tests to prove my fix is effective or that my feature works.
- [ ] All new and existing unit & integration tests pass locally (`go test -v ./...`).
- [ ] Security check: no workspace path escape, no command injection vulnerabilities, all dangerous ops respect HITL.
- [ ] Documentation has been updated accordingly.

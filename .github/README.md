# GitHub Configuration

This directory contains GitHub-related configuration for the TLS Agent project.

## 📂 Directory Structure

```
.github/
├── workflows/
│   └── pre-commit.yml              # GitHub Actions workflow for pre-commit hooks
├── ISSUE_TEMPLATE/
│   └── bug_report.md               # Template for bug reports
├── CODEOWNERS                       # Code ownership configuration
├── INTEGRATION_CHECKLIST.md        # Setup checklist
└── README.md                       # This file
```

## ⚙️ Configuration Files

### CODEOWNERS
Defines code ownership and review requirements:
- Owner: sbusanelli (sbusanelli@gmail.com)
- Coverage: All files in repository
- Auto-requests reviews: Yes

### .github/workflows/pre-commit.yml
GitHub Actions workflow that:
- Runs pre-commit hooks on pull requests
- Ensures code quality standards
- Validates Go code formatting, linting, and tests

## 🔍 Development Workflow

### Code Quality Checks
- ✅ Go code formatting (gofmt, gofumpt)
- ✅ Linting (golangci-lint, revive, go vet)
- ✅ Security scanning (gosec, detect-secrets)
- ✅ Tests pass (go test -race)
- ✅ Compilation succeeds (go build)
- ✅ Dependencies are tidy (go mod tidy)

### Review Process
1. Create feature branch
2. Make changes with commit messages
3. Open pull request
4. Automatic code quality checks run
5. CODEOWNERS automatically requested for review
6. Merge after approval and checks pass

## 📞 Support Resources

- **GitHub CODEOWNERS**: https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners
- **GitHub Actions**: https://docs.github.com/en/actions
- **Pre-commit Hooks**: https://pre-commit.com/
- **Go Best Practices**: https://golang.org/doc/effective_go

---

**Project Status**: 🟢 Active Development

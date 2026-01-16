# 🐰 GitHub Configuration & CodeRabbit Setup

This directory contains all GitHub-related configuration for the TLS Agent project, including CodeRabbit AI code review integration.

## 📂 Directory Structure

```
.github/
├── workflows/
│   └── coderabbit.yml              # GitHub Actions workflow for CodeRabbit reviews
├── ISSUE_TEMPLATE/
│   └── coderabbit-review.md        # Template for CodeRabbit review issues
├── coderabbit.yaml                 # Detailed CodeRabbit configuration
├── CODERABBIT_SETUP.md            # Complete setup guide
├── INTEGRATION_CHECKLIST.md        # Activation checklist
└── README.md                       # This file
```

## 🚀 Quick Start

### 1. Install CodeRabbit
- Visit: https://github.com/apps/coderabbit-ai
- Click "Install"
- Select this repository

### 2. Add API Key
- Go to repo Settings → Secrets and variables → Actions
- Create new secret: `CODERABBIT_API_KEY`
- Paste your CodeRabbit API key

### 3. Create Test PR
```bash
git checkout -b test/coderabbit
echo "# Test" >> test.txt
git add test.txt
git commit -m "Test CodeRabbit"
git push origin test/coderabbit
```

CodeRabbit will review automatically!

## 📖 Documentation Files

### CODERABBIT_SETUP.md
Complete guide covering:
- Feature overview
- Installation steps
- Configuration details
- Usage examples
- Troubleshooting

👉 **Start here for detailed setup**

### INTEGRATION_CHECKLIST.md
Step-by-step checklist with:
- Setup completion status
- Next steps
- Configuration details
- Success criteria

👉 **Use this to track activation progress**

## ⚙️ Configuration Files

### coderabbit.yaml (Root Level)
Main configuration specifying:
- Model: GPT-4
- Language: Go
- Auto-review settings
- File patterns to review
- Project context

### .github/coderabbit.yaml
Detailed configuration with:
- Exclusion patterns (generated files, vendor, etc.)
- Inclusion patterns (Go, proto, yaml files)
- Review rules (quality, security, performance, best practices)
- Focus areas (shutdown, certs, TLS, concurrency, error handling)
- Severity levels
- Notification settings

### .github/workflows/coderabbit.yml
GitHub Actions workflow that:
- Triggers on PR events
- Runs CodeRabbit reviews
- Posts comments with results
- Integrates with GitHub checks

## 🎯 What Gets Reviewed

### ✅ Reviewed Files
- `*.go` - Go source code
- `*.proto` - Protocol buffer definitions
- `*.yaml`/`*.yml` - Configuration files
- `go.mod`/`go.sum` - Module files
- `*.md` - Documentation

### ❌ Skipped Files
- `*.pb.go` - Generated protobuf
- `*.pb.gw.go` - Generated gateway
- `vendor/**` - Dependencies
- `bin/**` - Binaries
- `node_modules/**` - Node packages

## 🔍 Review Focus Areas

| Area | Description |
|------|-------------|
| **Graceful Shutdown** | Signal handling, timeouts, clean exits |
| **Certificate Mgmt** | TLS config, cert loading, hot reload |
| **Security** | TLS best practices, input validation |
| **Concurrency** | Goroutines, channels, race conditions |
| **Error Handling** | Error wrapping, propagation, recovery |

## �� Severity Levels

- 🔴 **Critical** - Security vulnerabilities, memory safety
- 🟠 **High** - Race conditions, goroutine leaks
- 🟡 **Medium** - Code quality, performance issues
- 🔵 **Low** - Style suggestions, documentation

## 🛠️ Customization

### Change Review Rules
Edit `coderabbit.yaml` or `.github/coderabbit.yaml`

### Exclude More Files
```yaml
exclude_patterns:
  - "path/to/exclude/**"
```

### Modify Focus Areas
```yaml
focus_areas:
  - your_focus_area
```

## 📞 Support

- **CodeRabbit Docs**: https://coderabbit.ai/docs
- **GitHub Actions**: https://docs.github.com/en/actions
- **Setup Guide**: See `CODERABBIT_SETUP.md`
- **Activation Checklist**: See `INTEGRATION_CHECKLIST.md`

## ✨ Status

- ✅ Configuration created
- ✅ Workflow configured
- ✅ Documentation written
- ⏳ Awaiting API key setup
- ⏳ Awaiting first PR review

---

**Next Step**: Add `CODERABBIT_API_KEY` to GitHub Secrets and create a test PR!

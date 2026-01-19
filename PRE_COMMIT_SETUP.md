# Pre-Commit Hooks Setup Guide

This document explains how to set up and use pre-commit hooks for Go linting and code quality checks.

## Overview

Pre-commit hooks automatically run code quality checks before each commit, preventing poor-quality code from being committed. This project includes:

- **Go Linting** - golangci-lint with 15+ linters
- **Code Formatting** - gofmt and gofumpt
- **Code Analysis** - go vet and staticcheck
- **Security Scanning** - gosec and detect-secrets
- **File Validation** - YAML, JSON, merge conflicts, large files
- **Build & Test Verification** - go build and go test
- **Dependency Management** - go mod tidy

## Quick Start

### 1. Install Pre-Commit

Choose your preferred installation method:

**Option A: Using pip (recommended)**
```bash
pip install pre-commit
```

**Option B: Using Homebrew (macOS)**
```bash
brew install pre-commit
```

**Option C: Using conda**
```bash
conda install -c conda-forge pre-commit
```

**Option D: Using apt (Ubuntu/Debian)**
```bash
sudo apt-get install pre-commit
```

### 2. Install Go Tools

Make sure you have the required Go tools installed:

```bash
# golangci-lint (linter aggregator)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
  sh -s -- -b $(go env GOPATH)/bin

# gosec (security scanner - optional)
curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | \
  sh -s -- -b $(go env GOPATH)/bin
```

### 3. Install Git Hooks

Run the setup script:

```bash
chmod +x setup-pre-commit-hooks.sh
./setup-pre-commit-hooks.sh
```

Or manually install:

```bash
pre-commit install
pre-commit install --hook-type commit-msg
```

## Configuration Files

### `.pre-commit-config.yaml`

Main configuration file that defines which pre-commit hooks to run. Key sections:

```yaml
repos:
  - repo: https://github.com/golangci/golangci-lint
    hooks:
      - id: golangci-lint
```

**Key Features:**
- Specifies hook sources and versions
- Defines when hooks run (commit, push, etc.)
- Configures hook behavior

### `.golangci.yaml`

Golangci-lint configuration that defines which linters to run:

```yaml
linters:
  enable:
    - errcheck      # Finds unchecked errors
    - govet         # Go vet analysis
    - staticcheck   # Static analysis
    - gosec         # Security issues
    - revive        # Code style
    - gofmt         # Code formatting
    # ... and more
```

**Available Linters:**

| Linter | Purpose | Severity |
|--------|---------|----------|
| `errcheck` | Find unchecked errors | High |
| `govet` | Suspicious constructs | High |
| `staticcheck` | Advanced static analysis | High |
| `gosec` | Security problems | Critical |
| `unused` | Unused code | Medium |
| `ineffassign` | Ineffectual assignments | Medium |
| `gosimple` | Simplify code | Low |
| `typecheck` | Type checking | High |
| `goimports` | Unused imports | Low |
| `gofmt` | Code formatting | Low |
| `revive` | Code style issues | Medium |

### `.secrets.baseline`

Baseline file for the `detect-secrets` hook. Tracks known secrets that should be ignored.

## Using Pre-Commit Hooks

### Run Hooks Before Commit (Automatic)

Hooks run automatically when you commit:

```bash
git add .
git commit -m "my changes"
# Hooks run automatically
```

If any hook fails, the commit is blocked. Fix the issues and retry.

### Run Hooks Manually

Run all hooks on staged files:

```bash
pre-commit run
```

Run all hooks on all files:

```bash
pre-commit run --all-files
```

Run a specific hook:

```bash
pre-commit run golangci-lint --all-files
```

### Skip Hooks (Not Recommended)

If you must skip hooks:

```bash
git commit --no-verify
```

**Warning:** Bypassing hooks should be rare. They exist to maintain code quality!

### Update Hooks

Update all hooks to their latest versions:

```bash
pre-commit autoupdate
```

Update a specific hook:

```bash
pre-commit autoupdate --repo https://github.com/golangci/golangci-lint
```

## What Each Hook Does

### Go-Specific Hooks

#### `golangci-lint`
Runs 15+ linters in parallel to check for:
- Unchecked errors
- Unused code
- Security vulnerabilities
- Code style issues
- Performance problems

**Auto-fixes:** Yes (with `--fix` flag)

#### `go fmt`
Formats Go code according to Go standards.

**Auto-fixes:** Yes

#### `go vet`
Reports suspicious constructs in Go code.

**Auto-fixes:** No (requires manual fixes)

#### `go build`
Verifies that your code compiles.

**Auto-fixes:** No

#### `go test`
Runs unit tests to ensure functionality.

**Auto-fixes:** No

#### `go mod tidy`
Ensures `go.mod` and `go.sum` are properly maintained.

**Auto-fixes:** Yes

### General Hooks

#### `check-merge-conflict`
Detects merge conflict markers in files.

#### `check-yaml`
Validates YAML syntax.

#### `check-json`
Validates JSON syntax.

#### `check-added-large-files`
Prevents committing files larger than 500KB.

#### `trailing-whitespace`
Removes trailing whitespace.

**Auto-fixes:** Yes

#### `end-of-file-fixer`
Ensures files end with exactly one newline.

**Auto-fixes:** Yes

#### `detect-secrets`
Scans for hardcoded secrets, API keys, passwords, etc.

**Auto-fixes:** No (requires manual review)

#### `gosec`
Security scanner for Go code.

**Auto-fixes:** No (requires manual fixes)

## Common Issues and Solutions

### Issue: "pre-commit: command not found"

**Solution:** Install pre-commit:
```bash
pip install pre-commit
```

### Issue: "golangci-lint: command not found"

**Solution:** Install golangci-lint:
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
  sh -s -- -b $(go env GOPATH)/bin
```

### Issue: Hooks are slow

**Solution:**
1. Skip slow hooks in CI:
   ```bash
   SKIP=gosec,go-test pre-commit run --all-files
   ```

2. Only run on changed files:
   ```bash
   pre-commit run
   ```

### Issue: "go mod tidy" hook is making unwanted changes

**Solution:** Run manually to see what it wants to change:
```bash
go mod tidy
git diff go.mod go.sum
```

### Issue: "detect-secrets" is flagging a false positive

**Solution:** Add to the baseline:
```bash
detect-secrets scan --baseline .secrets.baseline
```

### Issue: Hooks failed but I want to commit anyway

**Solution:** Only for urgent cases:
```bash
git commit --no-verify
```

Then fix the issues and recommit without `--no-verify`.

## CI/CD Integration

### GitHub Actions

Example workflow file (`.github/workflows/pre-commit.yml`):

```yaml
name: pre-commit

on:
  pull_request:
  push:

jobs:
  pre-commit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-python@v4
      - uses: pre-commit/action@v3
```

### GitLab CI

Example configuration (`.gitlab-ci.yml`):

```yaml
pre-commit:
  image: python:3.11
  before_script:
    - pip install pre-commit
    - apt-get update && apt-get install -y golang-go
  script:
    - pre-commit run --all-files
```

## Project-Specific Notes

### For Contributors

1. Install pre-commit before making changes
2. All hooks must pass before commits
3. Hooks auto-fix formatting and style issues
4. Manual review required for security issues

### For Maintainers

1. Regularly update hooks: `pre-commit autoupdate`
2. Review security findings from `gosec` and `detect-secrets`
3. Update `.secrets.baseline` when appropriate
4. Monitor hook performance and adjust as needed

## Advanced Configuration

### Skip Hooks for Specific Commits

```bash
# Skip all hooks
SKIP=* git commit -m "WIP: debugging"

# Skip specific hooks
SKIP=gosec,go-test git commit -m "Temporary changes"
```

### Custom Hook Behavior

Edit `.pre-commit-config.yaml` to:

```yaml
- repo: https://github.com/golangci/golangci-lint
  rev: v1.55.2
  hooks:
    - id: golangci-lint
      args: ['--fix', '--issues-exit-code=1']
      stages: [commit]
```

### Different Hooks for Different Branches

Use `stages` to control when hooks run:

```yaml
stages:
  - commit      # Runs on `git commit`
  - push        # Runs on `git push`
  - manual      # Only when explicitly run
```

## Resources

- [Pre-commit Documentation](https://pre-commit.com/)
- [Golangci-lint Documentation](https://golangci-lint.run/)
- [Go Security Checker (gosec)](https://github.com/securego/gosec)
- [Detect Secrets](https://github.com/Yelp/detect-secrets)

## Troubleshooting Commands

```bash
# Check pre-commit version
pre-commit --version

# Check hook configuration
pre-commit run --all-files --verbose

# Check what files would be processed
pre-commit run golangci-lint --all-files --dry-run

# Clean up pre-commit cache
pre-commit clean

# Uninstall hooks
pre-commit uninstall
pre-commit uninstall --hook-type commit-msg
```

## Next Steps

1. ✅ Install pre-commit
2. ✅ Install Go tools (golangci-lint, gosec)
3. ✅ Run setup script: `./setup-pre-commit-hooks.sh`
4. ✅ Test hooks: `pre-commit run --all-files`
5. ✅ Start developing with automatic quality checks!

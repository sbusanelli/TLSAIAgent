# Pre-Commit Hooks Quick Reference

## Installation

```bash
# Install pre-commit
pip install pre-commit

# Install Go linter tools
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
  sh -s -- -b $(go env GOPATH)/bin

# Install hooks in your local repo
./setup-pre-commit-hooks.sh
# OR
pre-commit install
```

## Common Commands

| Command | Purpose |
|---------|---------|
| `make install-hooks` | Install pre-commit hooks |
| `make run-hooks` | Run hooks on staged files |
| `make run-hooks-all` | Run hooks on all files |
| `make update-hooks` | Update hooks to latest versions |
| `make lint` | Run linter |
| `make fmt` | Format code |
| `make test` | Run tests |

## Pre-Commit Hooks Included

### Linting & Analysis
- ✅ **golangci-lint** - 15+ Go linters
- ✅ **go fmt** - Go code formatting
- ✅ **go vet** - Go code analysis
- ✅ **revive** - Style checker

### Security
- ✅ **gosec** - Go security scanner
- ✅ **detect-secrets** - Secret detection

### Build & Test
- ✅ **go build** - Compilation check
- ✅ **go test** - Unit tests
- ✅ **go mod tidy** - Dependency management

### File Validation
- ✅ **check-yaml** - YAML syntax
- ✅ **check-json** - JSON syntax
- ✅ **check-merge-conflict** - Merge markers
- ✅ **check-added-large-files** - File size limit
- ✅ **trailing-whitespace** - Whitespace cleanup
- ✅ **end-of-file-fixer** - File endings

## What Happens on Commit

When you run `git commit`:

1. ✅ Code is formatted (gofmt, gofumpt)
2. ✅ Code is analyzed (golangci-lint, go vet, revive)
3. ✅ Security is scanned (gosec, detect-secrets)
4. ✅ Code compiles (go build)
5. ✅ Tests pass (go test -race)
6. ✅ Dependencies are tidy (go mod tidy)
7. ✅ Files are validated (YAML, JSON, merge conflicts)
8. ✅ Whitespace is cleaned up

If any check fails, the commit is blocked. Fix issues and try again.

## Skip Hooks (Use Sparingly)

```bash
# Skip all hooks
git commit --no-verify

# Skip in Makefile
SKIP=gosec,go-test git commit -m "WIP"
```

## Auto-Fix Commands

Commands that automatically fix issues:

```bash
make fmt              # Format code
make lint-fix         # Fix linting issues
make run-hooks-all    # Some hooks auto-fix
```

## Troubleshooting

**Hooks not running?**
```bash
pre-commit install
```

**Update hooks**
```bash
pre-commit autoupdate
```

**Run verbose output**
```bash
pre-commit run --all-files --verbose
```

**Check installed hooks**
```bash
pre-commit run --all-files --dry-run
```

## Configuration Files

- `.pre-commit-config.yaml` - Hook configuration
- `.golangci.yaml` - Linter settings
- `.secrets.baseline` - Secret detection baseline
- `Makefile` - Development commands
- `.github/workflows/pre-commit.yml` - CI/CD

## CI/CD Integration

Pre-commit hooks also run in GitHub Actions on:
- ✅ Pull requests
- ✅ Pushes to main/develop

See `.github/workflows/pre-commit.yml` for details.

## See Also

- Full documentation: [PRE_COMMIT_SETUP.md](PRE_COMMIT_SETUP.md)
- Makefile targets: `make help`
- Pre-commit docs: https://pre-commit.com/
- Golangci-lint docs: https://golangci-lint.run/

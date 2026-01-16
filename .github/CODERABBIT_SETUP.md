# CodeRabbit Integration Setup Guide

## Overview

CodeRabbit is an AI-powered code review agent that automatically reviews pull requests and provides intelligent feedback on code quality, security, performance, and best practices.

## Features

✅ **Automated Code Review** - AI-powered reviews on every PR
✅ **Security Analysis** - Detects security vulnerabilities  
✅ **Performance Review** - Identifies performance issues
✅ **Go Best Practices** - Checks Go-specific patterns
✅ **Line-by-Line Comments** - Contextual feedback on specific lines
✅ **Summary Reports** - Comprehensive review summaries

## Setup Instructions

### 1. Install CodeRabbit App

1. Go to [CodeRabbit](https://coderabbit.ai/)
2. Sign in with your GitHub account
3. Click "Install App" and authorize the application
4. Select the repositories to enable CodeRabbit for (including this one)

### 2. Add API Key to GitHub Secrets

1. Get your CodeRabbit API key from your CodeRabbit account settings
2. Go to your GitHub repository
3. Navigate to **Settings** → **Secrets and variables** → **Actions**
4. Click **New repository secret**
5. Name: `CODERABBIT_API_KEY`
6. Value: Your CodeRabbit API key
7. Click **Add secret**

### 3. Configure CodeRabbit (Already Done!)

The following files have been created:

- **`.github/coderabbit.yaml`** - Main configuration file with rules and preferences
- **`.github/workflows/coderabbit.yml`** - GitHub Actions workflow that triggers reviews

## Configuration

### Excluded Files

The following are automatically excluded from review:
- Generated Protocol Buffer files (`*.pb.go`, `*.pb.gw.go`)
- Vendored dependencies
- Minified assets
- Built binaries

### Included Files

The following are reviewed:
- Go source files (`*.go`)
- Protocol Buffer definitions (`*.proto`)
- YAML configuration files
- Markdown documentation
- Module files

### Focus Areas

CodeRabbit will pay special attention to:
- **Graceful Shutdown** - Review shutdown implementation
- **Certificate Management** - Monitor cert loading patterns
- **TLS Security** - Verify TLS configurations
- **Concurrency** - Check goroutines and channels
- **Error Handling** - Validate error patterns

### Review Rules

#### Code Quality
- Unused variables and imports
- Code style and formatting
- Documentation completeness

#### Security
- Potential vulnerabilities
- TLS/SSL configurations
- Error handling patterns
- Input validation

#### Performance
- Goroutine management
- Memory allocation patterns
- Concurrent access issues
- Resource leaks

#### Go Best Practices
- Error handling completeness
- Defer usage patterns
- Goroutine leak prevention
- Nil pointer checks
- Race condition detection

## How It Works

### Workflow

1. **Pull Request Created/Updated**
   - GitHub Actions workflow is triggered
   - CodeRabbit reviews the changes

2. **Review Analysis**
   - Files are analyzed against configured rules
   - Security, performance, and quality issues identified
   - Best practices validated

3. **Feedback**
   - Comments added to specific lines
   - Summary comment posted on PR
   - Issues categorized by severity

### Severity Levels

- **🔴 Critical**: Security vulnerabilities, data loss risks, memory safety issues
- **🟠 High**: Race conditions, goroutine leaks, improper error handling
- **🟡 Medium**: Code quality issues, performance degradation, unused code
- **🔵 Low**: Style issues, documentation improvements, minor suggestions

## Usage Examples

### What CodeRabbit Reviews

✅ **Go Code Quality**
```go
// Will flag: unused variables, improper error handling
var unusedVar = 42  // ← Will be flagged
result := someFunction()  // ← Will check for error handling
```

✅ **Concurrency Issues**
```go
// Will flag: potential race conditions
go func() {
    // ...concurrent access...
}()
```

✅ **Resource Management**
```go
// Will flag: proper defer usage, resource cleanup
defer watcher.Close()
```

✅ **Security**
```go
// Will flag: TLS configurations, input validation
tlsCfg := &tls.Config{
    MinVersion: tls.VersionTLS12,  // ← Good practice
}
```

## Customization

### To Add New Rules

Edit `.github/coderabbit.yaml`:

```yaml
focus_areas:
  - your_new_focus_area
```

### To Modify Excluded Files

Edit `.github/coderabbit.yaml`:

```yaml
exclude_patterns:
  - "path/to/exclude/**"
```

### To Change Severity Levels

Edit `.github/coderabbit.yaml`:

```yaml
severity_levels:
  critical:
    - "your issue type"
```

## Troubleshooting

### CodeRabbit Not Reviewing

1. Verify the API key is set in GitHub Secrets
2. Check the workflow is enabled (`.github/workflows/coderabbit.yml`)
3. Ensure CodeRabbit app is installed on the repository
4. Check GitHub Actions permissions allow PRs to be reviewed

### Too Many/Few Comments

Adjust in `.github/coderabbit.yaml`:
- Change `min_changed_lines` to filter smaller changes
- Modify `max_files` to limit review scope
- Update severity settings

### Slow Reviews

Reduce `maxConcurrentReviews` or increase `timeout` in the workflow file.

## Example Review Output

When CodeRabbit reviews a PR, it will:

1. **Add line comments** on specific issues:
```
🐰 CodeRabbit: This error is not being handled. Consider adding error handling.
```

2. **Post a summary** comment:
```
## 🐰 CodeRabbit Review

- Files Reviewed: 3
- Issues Found: 2
  - Security: 0
  - Performance: 1
  - Quality: 1
```

## Best Practices

1. **Respond to Feedback** - Address CodeRabbit suggestions in subsequent commits
2. **Review the Suggestions** - Not all suggestions may apply to your use case
3. **Override When Needed** - You can dismiss suggestions that don't fit your project
4. **Keep Configuration Updated** - Update rules as project needs evolve
5. **Monitor Performance** - Check review times and adjust limits as needed

## Resources

- [CodeRabbit Documentation](https://coderabbit.ai/docs)
- [CodeRabbit GitHub App](https://github.com/apps/coderabbit-ai)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

## Next Steps

1. ✅ Configuration files created
2. ⏳ Add API key to GitHub Secrets
3. ⏳ Create a test PR to verify it's working
4. ⏳ Review and adjust rules based on PR feedback
5. ⏳ Share setup with team members

---

**Questions?** Refer to the configuration files or check the CodeRabbit documentation.

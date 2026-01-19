# CODEOWNERS Configuration

## Overview

The `.github/CODEOWNERS` file has been created to define code ownership and review requirements for the TLS Agent project.

## Owner Information

**Primary Owner**: sbusanelli (sbusanelli@gmail.com)

---

## What is CODEOWNERS?

The CODEOWNERS file tells GitHub who should be requested as a reviewer for pull requests that modify files in your repository.

### Benefits:
- ✅ Automatically request reviews from code owners
- ✅ Enforce code review before merge (with branch protection)
- ✅ Clear accountability for different parts of the codebase
- ✅ Prevent unauthorized changes to critical files
- ✅ Track code ownership across the project

---

## File Organization

### Coverage by Category

| Category | Pattern | Owner |
|----------|---------|-------|
| **All Files** | `*` | @sbusanelli |
| **Go Source** | `*.go` | @sbusanelli |
| **Tests** | `*_test.go` | @sbusanelli |
| **Configuration** | `*.yaml`, `*.yml` | @sbusanelli |
| **Workflows** | `.github/workflows/` | @sbusanelli |
| **Build** | `build.sh`, `Dockerfile` | @sbusanelli |
| **Dependencies** | `go.mod`, `go.sum` | @sbusanelli |
| **Protocol Buffers** | `*.proto`, `*.pb.go` | @sbusanelli |
| **Documentation** | `*.md` | @sbusanelli |
| **TLS/Certs** | `certs/`, `internal/tlsstore/` | @sbusanelli |
| **Agent** | `internal/agent/` | @sbusanelli |

---

## How It Works

### On Pull Request Creation

1. **File Detection**: GitHub detects which files were changed in the PR
2. **Owner Lookup**: Checks `.github/CODEOWNERS` for matching patterns
3. **Review Request**: Automatically requests review from matching owners
4. **Status Check**: PR shows that review is required

### Review Process

1. **Owner Notified**: @sbusanelli is notified of the review request
2. **Code Review**: Owner reviews the changes
3. **Approval**: Owner approves or requests changes
4. **Merge**: PR can be merged once approved (if branch protection enabled)

---

## GitHub Integration

### Automatic Review Requests

When a PR is created that modifies any files matching patterns in CODEOWNERS:
- ✅ @sbusanelli is automatically requested as a reviewer
- ✅ The PR shows "Review requested" status
- ✅ Notifications sent to the owner

### With Branch Protection

If you enable branch protection rules requiring CODEOWNERS approval:
- ✅ At least one code owner must approve before merge
- ✅ Code owner dismissals prevent merge
- ✅ Provides enforcement of review requirements

---

## Setting Up Branch Protection

To enforce code owner reviews:

1. Go to **Settings** → **Branches**
2. Click **Add rule** under "Branch protection rules"
3. Enter `main` as the branch name pattern
4. Enable **Require a pull request before merging**
5. Enable **Require approval from Code Owners**
6. Enable **Dismiss stale pull request approvals**
7. Click **Create**

---

## Patterns in This Configuration

### Universal Pattern
```
* @sbusanelli
```
All files in the repository are owned by sbusanelli.

### Specific File Types
```
*.go @sbusanelli
*_test.go @sbusanelli
```
All Go files and test files require sbusanelli review.

### Directory Pattern
```
.github/workflows/ @sbusanelli
certs/ @sbusanelli
internal/tlsstore/ @sbusanelli
```
Entire directories are assigned to sbusanelli.

### Nested Pattern
```
**/test/** @sbusanelli
**/*.pb.go @sbusanelli
```
Files matching patterns at any depth.

---

## For Future Team Expansion

When adding team members, update patterns like:

```markdown
# Multiple owners example:
*.go @sbusanelli @new-member

# Different owners for different areas:
internal/agent/ @sbusanelli
internal/tlsstore/ @security-team
.github/workflows/ @devops-team

# Team-based ownership:
*.go @tlsagent/core-reviewers
```

---

## Tips & Best Practices

1. **Keep It Simple**: Don't over-complicate ownership rules
2. **Match Reality**: Ensure owners actually maintain those files
3. **Document Changes**: Update CODEOWNERS when responsibilities change
4. **Review Regularly**: Audit ownership rules quarterly
5. **Use Teams**: For larger teams, use GitHub teams (@org/team-name)

---

## Verification

### Check If CODEOWNERS Is Working

1. Create a test PR on a non-main branch
2. Observe if code owners are automatically requested
3. Check PR status shows review requirements

### Common Issues

**Issue**: Code owners not requested
- ✅ Check CODEOWNERS syntax
- ✅ Verify GitHub username is correct (@username)
- ✅ Ensure file patterns match changed files

**Issue**: CODEOWNERS not found
- ✅ File must be at `.github/CODEOWNERS`
- ✅ No file extension
- ✅ Must be on the default branch

---

## Resources

- [About Code Owners](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)
- [CODEOWNERS Syntax](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners#codeowners-file-location)
- [Branch Protection Rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)

---

## Current Configuration Summary

✅ **CODEOWNERS file created**
✅ **Owner**: sbusanelli (sbusanelli@gmail.com)
✅ **Coverage**: All files and directories
✅ **Ready for**: Pull request reviews
⏳ **Next Step**: Enable branch protection if desired

---

**Status**: 🟢 ACTIVE

The CODEOWNERS file is now active and will automatically request reviews from sbusanelli on all pull requests.

cd /Users/sbusan01/go

# Install gosec if missing
if ! command -v gosec &> /dev/null; then
  go install github.com/securego/gosec/v2/cmd/gosec@latest
fi

# Run security scan
echo "=== Running gosec security scan ==="
gosec -fmt=json -out=gosec-report.json ./... || true
gosec ./...

# Check for hardcoded secrets / credentials
echo -e "\n=== Searching for hardcoded secrets ==="
if command -v gitleaks &> /dev/null; then
  gitleaks detect --source . --verbose
else
  echo "gitleaks not installed (optional). Install with: brew install gitleaks"
fi

# Check dependencies for known vulnerabilities
echo -e "\n=== Checking Go module vulnerabilities ==="
go list -json -m all | nancy sleuth || echo "nancy not installed (optional)"

# Manual checks
echo -e "\n=== Manual security checks ==="
echo "1. Checking for hardcoded passwords/keys in .go files:"
grep -rn "password\|secret\|apikey\|token" --include="*.go" . || echo "None found (good)"

echo -e "\n2. Checking for unsafe use of crypto:"
grep -rn "crypto/md5\|crypto/sha1\|DES\|RC4" --include="*.go" . || echo "None found (good)"

echo -e "\n3. Checking certificate validation:"
grep -rn "InsecureSkipVerify" --include="*.go" . || echo "None found (good)"

# Git security check
echo -e "\n=== Git history check ==="
echo "Checking for sensitive patterns in commits (last 20 commits):"
git log --all -p -S "password\|secret\|apikey" -20 --oneline || echo "None found (good)"

# Verify git status
git status

# Build the project
go mod tidy
go build -v -o bin/tls-agent ./main.go

# List final structure
echo "=== Project ready ==="
find . -type f \( -name "*.go" -o -name "*.sh" -o -name "go.mod" -o -name "README.md" \) | sort

#!/bin/bash
# setup-pre-commit-hooks.sh
# Script to install and configure pre-commit hooks for the project

set -e

echo "🔧 Setting up pre-commit hooks for TLS Agent..."
echo ""

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    echo "❌ pre-commit is not installed"
    echo ""
    echo "Install pre-commit using one of these methods:"
    echo ""
    echo "1. Using pip (recommended):"
    echo "   pip install pre-commit"
    echo ""
    echo "2. Using Homebrew (macOS):"
    echo "   brew install pre-commit"
    echo ""
    echo "3. Using conda:"
    echo "   conda install -c conda-forge pre-commit"
    echo ""
    exit 1
fi

echo "✅ pre-commit is installed"
echo "   Version: $(pre-commit --version)"
echo ""

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "❌ golangci-lint is not installed"
    echo ""
    echo "Install golangci-lint using one of these methods:"
    echo ""
    echo "1. Using curl:"
    echo "   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
    echo ""
    echo "2. Using Homebrew (macOS):"
    echo "   brew install golangci-lint"
    echo ""
    echo "3. Using snap (Linux):"
    echo "   snap install golangci-lint"
    echo ""
    exit 1
fi

echo "✅ golangci-lint is installed"
echo "   Version: $(golangci-lint --version)"
echo ""

# Check if gosec is installed
if ! command -v gosec &> /dev/null; then
    echo "⚠️  gosec is not installed (optional for enhanced security checks)"
    echo ""
    echo "Install gosec using:"
    echo "   curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin"
    echo ""
fi

echo "⏳ Initializing git hooks..."
cd "$(dirname "$0")"

# Install pre-commit hooks
pre-commit install
pre-commit install --hook-type commit-msg

echo ""
echo "✅ Pre-commit hooks installed successfully!"
echo ""
echo "📋 Configured hooks:"
echo "   • golangci-lint     - Go code linting"
echo "   • go fmt            - Go code formatting"
echo "   • go vet            - Go code analysis"
echo "   • check-merge-conflict - Detects merge conflicts"
echo "   • check-yaml        - YAML syntax validation"
echo "   • check-json        - JSON syntax validation"
echo "   • check-added-large-files - Prevents large files"
echo "   • trailing-whitespace - Removes trailing whitespace"
echo "   • end-of-file-fixer - Ensures files end with newline"
echo "   • detect-secrets    - Detects hardcoded secrets"
echo "   • gosec             - Security vulnerability scanner"
echo "   • go build          - Verifies code compilation"
echo "   • go test           - Runs unit tests"
echo "   • go mod tidy       - Ensures go.mod is tidy"
echo ""

# Run hooks on all files to verify setup
echo "🧪 Testing hooks on all files..."
if pre-commit run --all-files; then
    echo ""
    echo "✅ All pre-commit hooks passed!"
else
    echo ""
    echo "⚠️  Some hooks found issues. Review and fix them:"
    echo "   • Run: pre-commit run --all-files"
    echo "   • Or for a specific hook: pre-commit run <hook-id> --all-files"
fi

echo ""
echo "📚 Usage:"
echo "   • Run pre-commit on staged files: pre-commit run"
echo "   • Run on all files: pre-commit run --all-files"
echo "   • Bypass hooks: git commit --no-verify (not recommended!)"
echo "   • Update hooks: pre-commit autoupdate"
echo ""

# Verify git status
git status

# Build the project
go mod tidy
go build -v -o bin/tls-agent ./main.go

# List final structure
echo "=== Project ready ==="
find . -type f \( -name "*.go" -o -name "*.sh" -o -name "go.mod" -o -name "README.md" \) | sort

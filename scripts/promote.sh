#!/bin/bash

# Environment Promotion Script
# Usage: ./scripts/promote.sh [environment] [options]

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
ENVIRONMENT=""
FORCE=false
SKIP_TESTS=false
VERSION=""
DRY_RUN=false

# Help function
show_help() {
    cat << EOF
Environment Promotion Script

USAGE:
    ./scripts/promote.sh [ENVIRONMENT] [OPTIONS]

ENVIRONMENTS:
    development     Deploy to development environment
    staging        Deploy to staging environment (requires dev validation)
    production     Deploy to production environment (requires staging validation)

OPTIONS:
    --force         Force promotion (bypasses quality gates)
    --skip-tests    Skip automated tests
    --version VER   Specify version to deploy
    --dry-run       Show what would be deployed without actually deploying
    --help          Show this help message

EXAMPLES:
    ./scripts/promote.sh development
    ./scripts/promote.sh staging --version v1.2.0
    ./scripts/promote.sh production --force
    ./scripts/promote.sh staging --dry-run

ENVIRONMENT VARIABLES:
    DOCKER_REGISTRY    Container registry (default: ghcr.io)
    IMAGE_NAME        Image name (default: tlsai-agent)
    KUBECONFIG       Kubernetes config file (for k8s deployments)
EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        development|staging|production)
            ENVIRONMENT="$1"
            shift
            ;;
        --force)
            FORCE=true
            shift
            ;;
        --skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}Error: Unknown option $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# Validate environment argument
if [[ -z "$ENVIRONMENT" ]]; then
    echo -e "${RED}Error: Environment is required${NC}"
    show_help
    exit 1
fi

# Set default values
DOCKER_REGISTRY="${DOCKER_REGISTRY:-ghcr.io}"
IMAGE_NAME="${IMAGE_NAME:-tlsai-agent}"

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "INFO")
            echo -e "${BLUE}[INFO]${NC} $message"
            ;;
        "SUCCESS")
            echo -e "${GREEN}[SUCCESS]${NC} $message"
            ;;
        "WARNING")
            echo -e "${YELLOW}[WARNING]${NC} $message"
            ;;
        "ERROR")
            echo -e "${RED}[ERROR]${NC} $message"
            ;;
    esac
}

# Function to check prerequisites
check_prerequisites() {
    print_status "INFO" "Checking prerequisites..."
    
    # Check if Docker is running
    if ! docker info > /dev/null 2>&1; then
        print_status "ERROR" "Docker is not running or not accessible"
        exit 1
    fi
    
    # Check if docker-compose is available
    if ! command -v docker-compose &> /dev/null; then
        print_status "ERROR" "docker-compose is not installed"
        exit 1
    fi
    
    # Check if required files exist
    local compose_file="docker-compose.${ENVIRONMENT}.yml"
    if [[ ! -f "$compose_file" ]]; then
        print_status "ERROR" "Compose file not found: $compose_file"
        exit 1
    fi
    
    # Check environment-specific requirements
    case $ENVIRONMENT in
        "staging")
            if [[ ! -f "docker-compose.dev.yml" ]]; then
                print_status "ERROR" "Development environment must be set up first"
                exit 1
            fi
            ;;
        "production")
            if [[ ! -f "docker-compose.staging.yml" ]]; then
                print_status "ERROR" "Staging environment must be set up first"
                exit 1
            fi
            ;;
    esac
    
    print_status "SUCCESS" "Prerequisites check passed"
}

# Function to run quality gates
run_quality_gates() {
    if [[ "$SKIP_TESTS" == "true" ]]; then
        print_status "WARNING" "Skipping quality gates due to --skip-tests flag"
        return 0
    fi
    
    print_status "INFO" "Running quality gates..."
    
    # Run pre-commit hooks
    print_status "INFO" "Running pre-commit hooks..."
    if ! pre-commit run --all-files; then
        if [[ "$FORCE" != "true" ]]; then
            print_status "ERROR" "Pre-commit hooks failed. Use --force to bypass."
            exit 1
        else
            print_status "WARNING" "Pre-commit hooks failed, but proceeding due to --force flag"
        fi
    fi
    
    # Run unit tests
    print_status "INFO" "Running unit tests..."
    if ! go test -v -race -coverprofile=coverage.out ./...; then
        if [[ "$FORCE" != "true" ]]; then
            print_status "ERROR" "Unit tests failed. Use --force to bypass."
            exit 1
        else
            print_status "WARNING" "Unit tests failed, but proceeding due to --force flag"
        fi
    fi
    
    # Check coverage
    if [[ -f "coverage.out" ]]; then
        COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        print_status "INFO" "Test coverage: ${COVERAGE}%"
        
        if (( $(echo "$COVERAGE < 60" | bc -l) )); then
            if [[ "$FORCE" != "true" ]]; then
                print_status "ERROR" "Coverage below 60%. Use --force to bypass."
                exit 1
            else
                print_status "WARNING" "Coverage below 60%, but proceeding due to --force flag"
            fi
        fi
    fi
    
    # Run security scan
    print_status "INFO" "Running security scan..."
    if command -v gosec &> /dev/null; then
        if ! gosec ./...; then
            if [[ "$FORCE" != "true" ]]; then
                print_status "ERROR" "Security scan failed. Use --force to bypass."
                exit 1
            else
                print_status "WARNING" "Security scan failed, but proceeding due to --force flag"
            fi
        fi
    else
        print_status "WARNING" "gosec not found, skipping security scan"
    fi
    
    print_status "SUCCESS" "Quality gates passed"
}

# Function to build Docker image
build_image() {
    print_status "INFO" "Building Docker image..."
    
    local image_tag="${VERSION:-latest}"
    local full_image="${DOCKER_REGISTRY}/${IMAGE_NAME}:${image_tag}"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_status "INFO" "DRY RUN: Would build image $full_image"
        return 0
    fi
    
    # Build the image
    if ! docker build -t "$full_image" .; then
        print_status "ERROR" "Docker build failed"
        exit 1
    fi
    
    print_status "SUCCESS" "Docker image built: $full_image"
}

# Function to deploy to environment
deploy_environment() {
    local compose_file="docker-compose.${ENVIRONMENT}.yml"
    
    print_status "INFO" "Deploying to $ENVIRONMENT environment..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_status "INFO" "DRY RUN: Would deploy using $compose_file"
        return 0
    fi
    
    # Set environment variables
    export VERSION="${VERSION:-latest}"
    export BUILD_DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    export COMMIT_SHA="$(git rev-parse HEAD)"
    
    # Deploy using docker-compose
    if ! docker-compose -f "$compose_file" up -d; then
        print_status "ERROR" "Deployment to $ENVIRONMENT failed"
        exit 1
    fi
    
    print_status "SUCCESS" "Deployed to $ENVIRONMENT environment"
}

# Function to run health checks
run_health_checks() {
    print_status "INFO" "Running health checks..."
    
    if [[ "$DRY_RUN" == "true" ]]; then
        print_status "INFO" "DRY RUN: Would run health checks"
        return 0
    fi
    
    # Wait for services to start
    sleep 30
    
    # Check if main service is healthy
    local container_name="tlsai-agent-${ENVIRONMENT}"
    if [[ "$ENVIRONMENT" == "development" ]]; then
        container_name="tlsai-agent-dev"
    fi
    
    if docker ps --filter "name=$container_name" --filter "status=running" | grep -q .; then
        print_status "SUCCESS" "Service $container_name is running"
    else
        print_status "ERROR" "Service $container_name is not running"
        docker logs "$container_name" --tail 50
        exit 1
    fi
    
    # Check health endpoint
    local health_url="https://localhost:8443/health"
    if [[ "$ENVIRONMENT" == "development" ]]; then
        health_url="https://localhost:8443/health"
    fi
    
    # Wait a bit more for health endpoint
    sleep 10
    
    if curl -f -k -s "$health_url" > /dev/null; then
        print_status "SUCCESS" "Health check passed for $ENVIRONMENT"
    else
        print_status "WARNING" "Health check failed for $ENVIRONMENT"
    fi
}

# Function to show deployment summary
show_summary() {
    print_status "INFO" "Deployment Summary"
    echo "=================="
    echo "Environment: $ENVIRONMENT"
    echo "Version: ${VERSION:-latest}"
    echo "Registry: $DOCKER_REGISTRY"
    echo "Image: $IMAGE_NAME"
    echo "Force: $FORCE"
    echo "Skip Tests: $SKIP_TESTS"
    echo "Dry Run: $DRY_RUN"
    echo "=================="
}

# Main execution
main() {
    print_status "INFO" "Starting environment promotion to $ENVIRONMENT"
    
    # Check prerequisites
    check_prerequisites
    
    # Run quality gates
    run_quality_gates
    
    # Build image
    build_image
    
    # Deploy environment
    deploy_environment
    
    # Run health checks
    run_health_checks
    
    # Show summary
    show_summary
    
    print_status "SUCCESS" "Environment promotion to $ENVIRONMENT completed successfully!"
    
    if [[ "$ENVIRONMENT" == "staging" ]]; then
        print_status "INFO" "Next step: Promote to production with: ./scripts/promote.sh production"
    fi
}

# Run main function
main "$@"

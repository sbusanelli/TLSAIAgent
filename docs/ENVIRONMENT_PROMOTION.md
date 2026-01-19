# Environment Promotion Guide

This guide explains how to promote the TLSAIAgent through different environments (Development → Staging → Production) using the automated CI/CD pipeline.

## 🏗️ Environment Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Development   │───▶│     Staging     │───▶│   Production   │
│                 │    │                 │    │                 │
│ • Debug mode    │    │ • Integration   │    │ • Production   │
│ • Local dev     │    │ • Full stack    │    │ • Monitoring   │
│ • Hot reload    │    │ • Performance  │    │ • Security     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### Local Development

```bash
# Start development environment
make dev

# View logs
make dev-logs

# Stop development environment
make dev-stop
```

### Environment Promotion

```bash
# Promote to development
make promote-dev

# Promote to staging (after dev validation)
make promote-staging

# Promote to production (after staging validation)
make promote-prod

# Or use the generic command
make promote ENV=staging
```

## 📋 Promotion Process

### 1. Development Environment

**Purpose**: Local development and testing

**Triggers**:
- Push to `develop` branch
- Manual trigger
- Local `make promote-dev`

**Quality Gates**:
- Pre-commit hooks
- Unit tests (≥60% coverage)
- Security scan
- Code formatting

**Services**:
- TLSAIAgent (debug mode)
- Redis
- PostgreSQL
- Grafana (admin/admin123)

**Access**:
- Agent: https://localhost:8443
- Grafana: http://localhost:3000

### 2. Staging Environment

**Purpose**: Integration testing and performance validation

**Triggers**:
- Push to `develop` branch (auto-promotion)
- Manual trigger with quality gates
- `make promote-staging`

**Prerequisites**:
- Development environment validated
- All quality gates passed
- Integration tests successful

**Quality Gates**:
- All development gates
- Integration tests
- Performance benchmarks
- Health checks

**Services**:
- TLSAIAgent (info logging)
- Redis
- PostgreSQL
- Prometheus
- Grafana (admin/staging123)

**Access**:
- Agent: https://staging.tlsai-agent.example.com
- Grafana: https://staging.tlsai-agent.example.com/grafana

### 3. Production Environment

**Purpose**: Production deployment with full monitoring

**Triggers**:
- Push to `main` branch (after staging validation)
- Version tags (v*)
- Manual trigger with approval
- `make promote-prod`

**Prerequisites**:
- Staging environment validated
- Manual approval required
- All quality gates passed
- Security review complete

**Quality Gates**:
- All staging gates
- Manual approval
- Production health checks
- Rollback capability verified

**Services**:
- TLSAIAgent (warn logging)
- Redis
- PostgreSQL
- Prometheus
- Grafana
- Nginx reverse proxy

**Access**:
- Agent: https://tlsai-agent.example.com
- Grafana: https://tlsai-agent.example.com/grafana

## 🔧 Configuration

### Environment Variables

Each environment has its own configuration file:

- `.env.dev` - Development settings
- `.env.staging` - Staging settings  
- `.env.production` - Production settings (secrets required)

### Docker Compose Files

- `docker-compose.dev.yml` - Development stack
- `docker-compose.staging.yml` - Staging stack
- `docker-compose.prod.yml` - Production stack

### Loading Environment Variables

```bash
# Load development environment
make env-dev

# Load staging environment  
make env-staging

# Load production environment (requires secrets)
make env-prod
```

## 📊 Quality Gates

### Code Quality

- **Coverage**: ≥60% test coverage required
- **Security**: ≤5 high-priority security issues
- **Formatting**: All code must pass gofmt
- **Linting**: All linting rules must pass

### Testing

- **Unit Tests**: All tests must pass with race detection
- **Integration Tests**: Full stack integration validation
- **Performance Tests**: Benchmark regression testing
- **Security Tests**: Vulnerability scanning

### Health Checks

- **Application Health**: HTTPS endpoint must respond
- **Service Health**: All dependencies must be healthy
- **Resource Health**: CPU/Memory within limits
- **Security Health**: No critical vulnerabilities

## 🔄 Rollback Process

### Automatic Rollback

Rollback is triggered automatically if:
- Health checks fail
- Critical security issues detected
- Performance degradation >50%

### Manual Rollback

```bash
# Via GitHub Actions
1. Go to Actions → Environment Promotion
2. Click "Run workflow"
3. Select "Rollback deployment"
4. Choose target environment

# Via CLI
gh workflow run environment-promotion.yml \
  --field rollback=true \
  --field environment=production
```

## 📈 Monitoring

### Metrics Collection

- **Application Metrics**: Response time, error rate, throughput
- **System Metrics**: CPU, memory, disk, network
- **Business Metrics**: User activity, feature usage
- **Security Metrics**: Authentication failures, suspicious activity

### Alerting

- **Critical**: Service down, security breach
- **Warning**: High latency, resource usage >80%
- **Info**: Deployments, configuration changes

### Dashboards

- **Overview**: System health and performance
- **Application**: Feature-specific metrics
- **Infrastructure**: Resource utilization
- **Security**: Threat detection and compliance

## 🛠️ Troubleshooting

### Common Issues

**Build Failures**:
```bash
# Check build logs
docker-compose logs tlsai-agent

# Rebuild with debug
docker-compose build --no-cache tlsai-agent
```

**Health Check Failures**:
```bash
# Check service status
docker-compose ps

# Test health endpoint
curl -k https://localhost:8443/health
```

**Permission Issues**:
```bash
# Fix certificate permissions
sudo chown -R $USER:$USER certs/
chmod 600 certs/server.key
```

### Debug Mode

Enable debug logging for troubleshooting:

```bash
# Development
export LOG_LEVEL=debug
make dev

# Staging
export LOG_LEVEL=debug
make promote-staging --force
```

## 📝 Best Practices

### Development

1. **Always run pre-commit hooks** before committing
2. **Write tests** for new features
3. **Update documentation** for API changes
4. **Use feature flags** for experimental features

### Promotion

1. **Validate in development** before staging promotion
2. **Monitor staging** for at least 24 hours before production
3. **Use semantic versioning** for releases
4. **Document breaking changes** in release notes

### Production

1. **Monitor health** continuously after deployment
2. **Have rollback plan** ready before deployment
3. **Use canary deployments** for major changes
4. **Review security** logs regularly

## 🔗 Additional Resources

- [GitHub Actions Workflows](.github/workflows/)
- [Docker Configuration](docker-compose.*.yml)
- [Environment Variables](.env.*)
- [Monitoring Setup](config/)
- [Security Guidelines](docs/security.md)

## 📞 Support

For issues with environment promotion:

1. Check the [troubleshooting guide](#-troubleshooting)
2. Review [GitHub Actions logs](https://github.com/sbusanelli/TLSAIAgent/actions)
3. Create an [issue](https://github.com/sbusanelli/TLSAIAgent/issues)
4. Contact the DevOps team

---

**Note**: Always ensure you have the necessary permissions and approvals before promoting to production environments.

# 🚀 TLSAIAgent Enterprise Deployment Guide

## 📋 Overview

This comprehensive guide covers enterprise-grade deployment of TLSAIAgent using modern containerization, orchestration, and monitoring technologies. The deployment supports multiple environments including development, staging, and production with full security, monitoring, and observability capabilities.

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           TLSAIAgent Enterprise Architecture                      │
├─────────────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │
│  │   Ingress   │  │   Service   │  │  Deployment │  │   ConfigMap │           │
│  │   (HTTPS)   │  │   (LB)     │  │   (3 Pods)  │  │   (Config)  │           │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘           │
│         │               │               │               │                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │
│  │   Cert-Mgr  │  │ Prometheus  │  │  Grafana    │  │   Vault     │           │
│  │   (TLS)     │  │ (Metrics)   │  │ (Dashboards)│  │ (Secrets)   │           │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘           │
│         │               │               │               │                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐           │
│  │ ELK Stack   │  │   Redis     │  │ PostgreSQL  │  │   Nginx     │           │
│  │ (Logging)   │  │ (Cache)     │  │ (Database)  │  │ (Proxy)     │           │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘           │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## 🎯 Deployment Options

### 1. Docker Compose (Development/Small Scale)
- **Use Case**: Development, testing, small-scale deployments
- **Complexity**: Low
- **Scalability**: Limited
- **Monitoring**: Basic

### 2. Kubernetes (Production/Enterprise)
- **Use Case**: Production, large-scale, multi-tenant
- **Complexity**: High
- **Scalability**: Unlimited
- **Monitoring**: Advanced

### 3. Cloud Native (Managed Services)
- **Use Case**: Cloud-first, serverless, auto-scaling
- **Complexity**: Medium
- **Scalability**: Auto
- **Monitoring**: Cloud-native

## 🐳 Docker Compose Deployment

### Prerequisites

```bash
# Install Docker and Docker Compose
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### Quick Start

```bash
# Clone repository
git clone https://github.com/sbusanelli/TLSAIAgent.git
cd TLSAIAgent

# Create environment file
cp .env.example .env

# Edit environment variables
nano .env

# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f tlsai-agent
```

### Environment Configuration

```bash
# .env file
VERSION=v1.0.0
BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)
COMMIT_SHA=$(git rev-parse HEAD)

# Database Configuration
POSTGRES_DB=tlsai_agent
POSTGRES_USER=tlsai
POSTGRES_PASSWORD=your_secure_password

# Redis Configuration
REDIS_PASSWORD=your_redis_password

# Vault Configuration
VAULT_ADDR=https://vault:8200
VAULT_DEV_ROOT_TOKEN_ID=your_vault_token

# Grafana Configuration
GRAFANA_ADMIN_USER=admin
GRAFANA_ADMIN_PASSWORD=your_grafana_password

# Logging Configuration
LOG_LEVEL=info
METRICS_ENABLED=true
```

### Service Access

| Service | URL | Credentials |
|---------|-----|-------------|
| TLSAIAgent | http://localhost:8080 | - |
| TLSAIAgent (HTTPS) | https://localhost:8443 | - |
| Grafana | http://localhost:3000 | admin/admin123 |
| Prometheus | http://localhost:9091 | - |
| Kibana | http://localhost:5601 | - |
| Vault | https://localhost:8200 | token-based |

## ☸️ Kubernetes Deployment

### Prerequisites

```bash
# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Install ingress-nginx
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install ingress-nginx ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace
```

### Namespace Setup

```bash
# Create namespace
kubectl apply -f k8s/namespace.yaml

# Verify namespace
kubectl get namespace tlsai-agent
```

### Secrets Configuration

```bash
# Create TLS certificates
kubectl create secret tls tlsai-agent-certs \
  --cert=certs/tls.crt \
  --key=certs/tls.key \
  --namespace=tlsai-agent

# Create application secrets
kubectl create secret generic tlsai-agent-secrets \
  --from-literal=vault-token=your_vault_token \
  --from-literal=db-password=your_db_password \
  --namespace=tlsai-agent
```

### Deployment

```bash
# Deploy all components
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/ingress.yaml

# Verify deployment
kubectl get pods -n tlsai-agent
kubectl get services -n tlsai-agent
kubectl get ingress -n tlsai-agent
```

### Scaling and Updates

```bash
# Scale deployment
kubectl scale deployment tlsai-agent --replicas=5 -n tlsai-agent

# Update deployment
kubectl set image deployment/tlsai-agent tlsai-agent=tlsai-agent:v1.1.0 -n tlsai-agent

# Rollback deployment
kubectl rollout undo deployment/tlsai-agent -n tlsai-agent

# Check rollout status
kubectl rollout status deployment/tlsai-agent -n tlsai-agent
```

## 🔒 Security Configuration

### TLS/SSL Setup

```bash
# Generate self-signed certificates (development)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout certs/tls.key \
  -out certs/tls.crt \
  -subj "/CN=tlsai-agent.local"

# Use Let's Encrypt (production)
# Cert-manager will automatically handle certificate issuance
```

### Network Policies

```yaml
# Network policy for restricted access
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: tlsai-agent-netpol
  namespace: tlsai-agent
spec:
  podSelector:
    matchLabels:
      app: tlsai-agent
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: ingress-nginx
    ports:
    - protocol: TCP
      port: 8080
```

### RBAC Configuration

```yaml
# Service account and permissions
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tlsai-agent
  namespace: tlsai-agent

---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: tlsai-agent
  namespace: tlsai-agent
rules:
- apiGroups: [""]
  resources: ["configmaps", "secrets"]
  verbs: ["get", "list", "watch"]
```

## 📊 Monitoring and Observability

### Prometheus Metrics

```yaml
# Service monitor configuration
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: tlsai-agent
  namespace: tlsai-agent
spec:
  selector:
    matchLabels:
      app: tlsai-agent
  endpoints:
  - port: metrics
    path: /metrics
    interval: 30s
```

### Grafana Dashboards

```bash
# Import pre-built dashboards
kubectl create configmap grafana-dashboards \
  --from-file=grafana/dashboards/ \
  --namespace=monitoring

# Configure dashboard provisioning
kubectl apply -f grafana/provisioning/
```

### Logging Configuration

```yaml
# Logstash configuration
input {
  beats {
    port => 5044
  }
}

filter {
  if [kubernetes][namespace] == "tlsai-agent" {
    json {
      source => "message"
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "tlsai-agent-%{+YYYY.MM.dd}"
  }
}
```

## 🔧 Configuration Management

### ConfigMap Updates

```bash
# Update configuration
kubectl edit configmap tlsai-agent-config -n tlsai-agent

# Restart pods to apply changes
kubectl rollout restart deployment/tlsai-agent -n tlsai-agent
```

### Secret Management

```bash
# Update secrets
kubectl create secret generic tlsai-agent-secrets \
  --from-literal=new-secret=value \
  --dry-run=client -o yaml | kubectl apply -f -

# Rotate certificates
kubectl delete secret tlsai-agent-certs -n tlsai-agent
kubectl create secret tls tlsai-agent-certs \
  --cert=certs/new-tls.crt \
  --key=certs/new-tls.key \
  --namespace=tlsai-agent
```

## 🚀 CI/CD Integration

### GitHub Actions

```yaml
# .github/workflows/deploy.yml
name: Deploy TLSAIAgent

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Deploy to Kubernetes
      run: |
        kubectl apply -f k8s/
        kubectl rollout status deployment/tlsai-agent -n tlsai-agent
```

### Helm Chart

```yaml
# Chart.yaml
apiVersion: v2
name: tlsai-agent
description: TLSAIAgent Helm Chart
type: application
version: 1.0.0
appVersion: "1.0.0"

# values.yaml
replicaCount: 3
image:
  repository: tlsai-agent
  tag: "v1.0.0"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 8080

ingress:
  enabled: true
  className: nginx
  annotations: {}
  hosts:
    - host: tlsai-agent.example.com
      paths:
        - path: /
          pathType: Prefix
```

## 🔍 Troubleshooting

### Common Issues

1. **Pod Not Starting**
   ```bash
   kubectl describe pod <pod-name> -n tlsai-agent
   kubectl logs <pod-name> -n tlsai-agent
   ```

2. **Service Not Accessible**
   ```bash
   kubectl get svc -n tlsai-agent
   kubectl get endpoints -n tlsai-agent
   kubectl port-forward svc/tlsai-agent 8080:8080 -n tlsai-agent
   ```

3. **Certificate Issues**
   ```bash
   kubectl get certificates -n tlsai-agent
   kubectl describe certificate tlsai-agent-tls -n tlsai-agent
   kubectl logs -n cert-manager deployment/cert-manager
   ```

4. **High Memory Usage**
   ```bash
   kubectl top pods -n tlsai-agent
   kubectl exec -it <pod-name> -n tlsai-agent -- top
   ```

### Debug Commands

```bash
# Check pod status
kubectl get pods -n tlsai-agent -o wide

# Check events
kubectl get events -n tlsai-agent --sort-by=.metadata.creationTimestamp

# Check resource usage
kubectl top nodes
kubectl top pods -n tlsai-agent

# Port forward for debugging
kubectl port-forward svc/tlsai-agent 8080:8080 -n tlsai-agent

# Exec into pod
kubectl exec -it <pod-name> -n tlsai-agent -- /bin/sh

# Check logs
kubectl logs -f deployment/tlsai-agent -n tlsai-agent
```

## 📈 Performance Tuning

### Resource Limits

```yaml
# Optimize resource requests and limits
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

### Horizontal Pod Autoscaling

```yaml
# Configure HPA
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: tlsai-agent-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: tlsai-agent
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Database Optimization

```sql
-- PostgreSQL optimization
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
SELECT pg_reload_conf();
```

## 🔄 Backup and Recovery

### Database Backup

```bash
# PostgreSQL backup
kubectl exec -it postgres-pod -n tlsai-agent -- \
  pg_dump -U tlsai tlsai_agent > backup.sql

# Restore backup
kubectl exec -i postgres-pod -n tlsai-agent -- \
  psql -U tlsai tlsai_agent < backup.sql
```

### Configuration Backup

```bash
# Backup all configurations
kubectl get all -n tlsai-agent -o yaml > backup.yaml

# Restore configuration
kubectl apply -f backup.yaml
```

## 📚 Best Practices

### Security
- Use least privilege access
- Enable network policies
- Rotate secrets regularly
- Use TLS for all communications
- Enable audit logging

### Performance
- Set appropriate resource limits
- Use horizontal pod autoscaling
- Monitor key metrics
- Optimize database queries
- Use caching where appropriate

### Reliability
- Use multiple replicas
- Implement health checks
- Use graceful shutdown
- Monitor error rates
- Implement circuit breakers

### Observability
- Collect all metrics
- Centralize logging
- Set up alerting
- Use distributed tracing
- Monitor SLA compliance

## 🆘 Support and Maintenance

### Regular Tasks
- Daily: Check pod health and resource usage
- Weekly: Review logs and metrics
- Monthly: Update dependencies and certificates
- Quarterly: Review and update security policies

### Emergency Procedures
1. **Service Outage**
   - Check pod status
   - Review recent changes
   - Rollback if necessary
   - Communicate with stakeholders

2. **Security Incident**
   - Isolate affected systems
   - Collect evidence
   - Notify security team
   - Document timeline

3. **Performance Degradation**
   - Check resource usage
   - Review recent deployments
   - Scale resources if needed
   - Optimize configuration

---

## 📞 Contact Information

- **Maintainer**: Sreedhar Busanelli
- **Email**: sbusanelli@example.com
- **GitHub**: https://github.com/sbusanelli/TLSAIAgent
- **Documentation**: https://github.com/sbusanelli/TLSAIAgent/blob/main/docs/

---

*This deployment guide is part of the TLSAIAgent enterprise suite. For more information, please refer to the main documentation.*

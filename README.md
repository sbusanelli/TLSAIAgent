# TLS Hot Reload AI Agent (Go)

## Run
```
go mod tidy
go run main.go
```

## Test
Replace certs/server.crt and server.key with a new cert.
New TLS connections will immediately use the new cert without restart.

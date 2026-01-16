package main

import (
    "crypto/tls"
    "log"
    "net/http"

    "tls-agent/internal/agent"
    "tls-agent/internal/tlsstore"
)

func main() {
    cert, err := tlsstore.Load("certs/server.crt", "certs/server.key")
    if err != nil {
        log.Fatal(err)
    }

    store := tlsstore.New(cert)

    tlsCfg := &tls.Config{
        GetCertificate: store.GetCertificate,
        MinVersion:     tls.VersionTLS12,
    }

    state := agent.NewState(cert)
    go agent.Run(store, state)

    server := &http.Server{
        Addr:      ":8443",
        TLSConfig: tlsCfg,
    }

    log.Println("TLS Agent server running on :8443")
    log.Fatal(server.ListenAndServeTLS("", ""))
}

package agent

import (
    "crypto/tls"
    "log"
    "time"

    "tls-agent/internal/tlsstore"
)

type State struct {
    Current  *tls.Certificate
    Previous *tls.Certificate
    LastRun  time.Time
}

func NewState(cert *tls.Certificate) *State {
    return &State{Current: cert}
}

func Run(store *tlsstore.Store, state *State) {
    ticker := time.NewTicker(30 * time.Second)

    for range ticker.C {
        if state.Current.Leaf != nil && time.Until(state.Current.Leaf.NotAfter) < 7*24*time.Hour {
            log.Println("Agent: cert nearing expiry, attempting reload")

            cert, err := tlsstore.Load("certs/server.crt", "certs/server.key")
            if err != nil {
                log.Println("Agent reload failed:", err)
                continue
            }

            state.Previous = state.Current
            state.Current = cert
            store.Update(cert)

            log.Println("Agent: certificate rotated successfully")
        }

        state.LastRun = time.Now()
    }
}

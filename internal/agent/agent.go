package agent

import (
    "crypto/tls"
    "log"
    "time"

    "github.com/fsnotify/fsnotify"
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
    // Create file watcher for certificate files
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Println("Agent: failed to create watcher:", err)
        return
    }
    defer watcher.Close()

    // Watch certificate files
    if err := watcher.Add("certs/server.crt"); err != nil {
        log.Println("Agent: failed to watch server.crt:", err)
    }
    if err := watcher.Add("certs/server.key"); err != nil {
        log.Println("Agent: failed to watch server.key:", err)
    }

    log.Println("Agent: watching certs/server.crt and certs/server.key for changes")

    // Also run periodic checks (fallback, every 60 seconds)
    ticker := time.NewTicker(60 * time.Second)
    defer ticker.Stop()

    // Track recent reloads to avoid duplicate processing
    lastReloadTime := time.Now()
    reloadDebounce := 2 * time.Second

    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            // Ignore remove/rename events, only process write events
            if event.Has(fsnotify.Write) {
                now := time.Now()
                // Debounce: ignore reload if last reload was < 2 seconds ago
                if now.Sub(lastReloadTime) < reloadDebounce {
                    log.Println("Agent: debouncing rapid file changes")
                    continue
                }

                log.Println("Agent: detected certificate file change:", event.Name)
                if reloadCert(store, state) {
                    lastReloadTime = now
                }
            }

        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            log.Println("Agent: watcher error:", err)

        case <-ticker.C:
            // Periodic fallback check (e.g., detect external changes)
            if state.Current.Leaf != nil && time.Until(state.Current.Leaf.NotAfter) < 7*24*time.Hour {
                log.Println("Agent: cert nearing expiry (7 days), attempting reload")
                reloadCert(store, state)
            }
        }

        state.LastRun = time.Now()
    }
}

func reloadCert(store *tlsstore.Store, state *State) bool {
    cert, err := tlsstore.Load("certs/server.crt", "certs/server.key")
    if err != nil {
        log.Println("Agent: reload failed:", err)
        return false
    }

    state.Previous = state.Current
    state.Current = cert
    store.Update(cert)

    log.Println("Agent: certificate reloaded successfully")
    return true
}

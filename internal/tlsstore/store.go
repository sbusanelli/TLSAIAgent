package tlsstore

import (
    "crypto/tls"
    "sync/atomic"
)

type Store struct {
    cert atomic.Value
}

func New(initial *tls.Certificate) *Store {
    s := &Store{}
    s.cert.Store(initial)
    return s
}

func (s *Store) GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
    return s.cert.Load().(*tls.Certificate), nil
}

func (s *Store) Update(cert *tls.Certificate) {
    s.cert.Store(cert)
}

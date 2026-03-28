package tlsstore

import (
	"crypto/tls"
	"sync/atomic"
	"time"
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

// IsValid checks if the current certificate is valid and not expired
func (s *Store) IsValid() bool {
	cert := s.cert.Load().(*tls.Certificate)
	if cert == nil {
		return false
	}

	// If Leaf is not parsed, we consider the certificate valid
	// (it would be parsed during TLS handshake)
	if cert.Leaf == nil {
		return true
	}

	// Check if certificate is still valid (not expired)
	return time.Now().Before(cert.Leaf.NotAfter)
}

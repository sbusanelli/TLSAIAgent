package tlsstore

import (
	"crypto/tls"
	"os"
	"testing"
	"time"
)

// TestLoad tests certificate loading functionality
func TestLoad(t *testing.T) {
	// Test loading valid certificates
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	if cert == nil {
		t.Fatal("Certificate should not be nil")
	}

	if cert.Certificate == nil {
		t.Error("Certificate certificate should not be nil")
	}

	if cert.PrivateKey == nil {
		t.Error("Certificate private key should not be nil")
	}

	// Test certificate properties
	if len(cert.Certificate.Raw) == 0 {
		t.Error("Certificate raw bytes should not be empty")
	}

	if cert.Certificate.NotBefore.After(time.Now()) {
		t.Error("Certificate should be valid")
	}

	if cert.PrivateKey.D != 0 {
		t.Error("Private key should not be zero")
	}
}

// TestLoadInvalidFiles tests loading invalid certificate files
func TestLoadInvalidFiles(t *testing.T) {
	// Test non-existent files
	_, err := Load("nonexistent.crt", "nonexistent.key")
	if err == nil {
		t.Error("Loading non-existent files should fail")
	}

	// Test empty files
	tempDir := t.TempDir()
	emptyCert := tempDir + "/empty.crt"
	emptyKey := tempDir + "/empty.key"

	// Create empty files
	err = os.WriteFile(emptyCert, []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to create empty cert file: %v", err)
	}

	err = os.WriteFile(emptyKey, []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to create empty key file: %v", err)
	}

	// Try to load empty files
	_, err = Load(emptyCert, emptyKey)
	if err == nil {
		t.Error("Loading empty certificate files should fail")
	}

	// Clean up
	os.Remove(emptyCert)
	os.Remove(emptyKey)
}

// TestLoadMismatchedFiles tests loading mismatched certificate/key pairs
func TestLoadMismatchedFiles(t *testing.T) {
	tempDir := t.TempDir()
	certFile := tempDir + "/cert.crt"
	keyFile := tempDir + "/key.key"

	// Create test certificate and key files
	testCert := `-----BEGIN CERTIFICATE-----
MIIDdzCCAn+gAwIBAgI...
-----END CERTIFICATE-----`

	testKey := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQ...
-----END PRIVATE KEY-----`

	// Write certificate
	err := os.WriteFile(certFile, []byte(testCert), 0644)
	if err != nil {
		t.Fatalf("Failed to write certificate file: %v", err)
	}

	// Write different key
	err = os.WriteFile(keyFile, []byte("different-key"), 0644)
	if err != nil {
		t.Fatalf("Failed to write key file: %v", err)
	}

	// Try to load mismatched files
	_, err = Load(certFile, keyFile)
	if err == nil {
		t.Error("Loading mismatched certificate/key files should fail")
	}

	// Clean up
	os.Remove(certFile)
	os.Remove(keyFile)
}

// TestNew tests certificate store creation
func TestNew(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	if store == nil {
		t.Fatal("Store should not be nil")
	}

	// Test GetCertificate method
	retrievedCert := store.GetCertificate()
	if retrievedCert == nil {
		t.Error("Retrieved certificate should not be nil")
	}

	if retrievedCert != cert {
		t.Error("Retrieved certificate should match original")
	}
}

// TestGetCertificateWithValidFiles tests certificate retrieval
func TestGetCertificateWithValidFiles(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	// Test multiple retrievals
	for i := 0; < 10; i++ {
		retrievedCert := store.GetCertificate()
		if retrievedCert == nil {
			t.Errorf("Retrieved certificate should not be nil (iteration %d)", i)
		}

		if retrievedCert != cert {
			t.Errorf("Retrieved certificate should match original (iteration %d)", i)
		}
	}
}

// TestGetCertificateWithNilStore tests certificate retrieval with nil store
func TestGetCertificateWithNilStore(t *testing.T) {
	var store *CertificateStore

	// Test with nil store
	retrievedCert := store.GetCertificate()
	if retrievedCert != nil {
		t.Error("Retrieved certificate should be nil for nil store")
	}
}

// TestGetCertificateConcurrentAccess tests concurrent certificate retrieval
func TestGetCertificateConcurrentAccess(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			retrievedCert := store.GetCertificate()
			if retrievedCert == nil {
				t.Errorf("Retrieved certificate should not be nil (goroutine %d)", id)
			}
			if retrievedCert != cert {
				t.Errorf("Retrieved certificate should match original (goroutine %d)", id)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestCertificateValidation tests certificate validation
func TestCertificateValidation(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Test certificate validity
	if cert.Certificate == nil {
		t.Fatal("Certificate should not be nil")
	}

	// Test certificate expiration
	if cert.Certificate.NotBefore.After(time.Now()) {
		t.Error("Certificate should not be expired")
	}

	// Test certificate subject
	if len(cert.Certificate.Subject.CommonName) == 0 {
		t.Error("Certificate should have a common name")
	}

	// Test certificate issuer
	if len(cert.Certificate.Issuer.CommonName) == 0 {
		t.Error("Certificate should have an issuer")
	}

	// Test private key
	if cert.PrivateKey == nil {
		t.Error("Private key should not be nil")
	}

	if cert.PrivateKey.D == 0 {
		t.Error("Private key should not be zero")
	}
}

// TestPrivateKeyValidation tests private key validation
func TestPrivateKeyValidation(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Test private key properties
	if cert.PrivateKey == nil {
		t.Fatal("Private key should not be nil")
	}

	// Test private key algorithm
	if cert.PrivateKey.Algorithm.String() == "" {
		t.Error("Private key should have an algorithm")
	}

	// Test private key size
	if cert.PrivateKey.N == 0 {
		t.Error("Private key should have non-zero modulus")
	}

	// Test private key public exponent
	if cert.PrivateKey.E == 0 {
		t.Error("Private key should have non-zero public exponent")
	}
}

// TestCertificateExpiration tests certificate expiration handling
func TestCertificateExpiration(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Test not before time
	if cert.Certificate.NotBefore.IsZero() {
		t.Error("Certificate not before time should not be zero")
	}

	// Test not after time
	if cert.Certificate.NotAfter.IsZero() {
		t.Error("Certificate not after time should not be zero")
	}

	// Test certificate is currently valid
	now := time.Now()
	if now.Before(cert.Certificate.NotBefore) {
		t.Error("Certificate should be valid now")
	}

	if now.After(cert.Certificate.NotAfter) {
		t.Error("Certificate should not be expired")
	}

	// Test certificate validity period
	validityPeriod := cert.Certificate.NotAfter.Sub(cert.Certificate.NotBefore)
	if validityPeriod <= 0 {
		t.Error("Certificate validity period should be positive")
	}

	// Test certificate validity period is reasonable (at least 1 day)
	if validityPeriod < 24*time.Hour {
		t.Error("Certificate validity period should be at least 24 hours")
	}
}

// TestCertificateSubject tests certificate subject information
func TestCertificateSubject(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	subject := cert.Certificate.Subject

	// Test common name
	if len(subject.CommonName) == 0 {
		t.Error("Certificate should have a common name")
	}

	// Test organization
	if len(subject.Organization) == 0 {
		t.Error("Certificate should have an organization")
	}

	// Test organizational unit
	if len(subject.OrganizationalUnit) == 0 {
		t.Error("Certificate should have an organizational unit")
	}

	// Test country
	if len(subject.Country) == 0 {
		t.Error("Certificate should have a country")
	}
}

// TestCertificateIssuer tests certificate issuer information
func TestCertificateIssuer(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	issuer := cert.Certificate.Issuer

	// Test issuer common name
	if len(issuer.CommonName) == 0 {
		t.Error("Certificate should have an issuer common name")
	}

	// Test issuer organization
	if len(issuer.Organization) == 0 {
		t.Error("Certificate should have an issuer organization")
	}

	// Test issuer organizational unit
	if len(issuer.OrganizationalUnit) == 0 {
		t.Error("Certificate should have an issuer organizational unit")
	}
}

// TestCertificateExtensions tests certificate extensions
func TestCertificateExtensions(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	extensions := cert.Certificate.Extensions

	// Test extensions are present
	if len(extensions) == 0 {
		t.Error("Certificate should have extensions")
	}

	// Test basic constraints extension
	basicConstraints := cert.Certificate.BasicConstraints
	if basicConstraints == nil {
		t.Error("Certificate should have basic constraints")
	}

	// Test key usage extension
	keyUsage := cert.Certificate.KeyUsage
	if keyUsage == nil {
		t.ExtKeyUsage = []tls.KeyUsage{tls.KeyUsageDigitalSignature, tls.KeyUsageKeyEncipherment}
	}

	// Test extended key usage extension
	extKeyUsage := cert.Certificate.ExtKeyUsage
	if extKeyUsage == nil {
		t.ExtKeyUsage = []tls.ExtKeyUsage{tls.ExtKeyUsageServerAuth, tls.ExtKeyUsageClientAuth}
	}
}

// TestCertificateSanity tests certificate sanity checks
func TestCertificateSanity(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Test certificate is not self-signed (unless it's a root CA cert)
	if cert.Certificate.Issuer.CommonName == cert.Certificate.Subject.CommonName &&
		len(cert.Certificate.Authority) == 0 {
		t.Log("Certificate appears to be self-signed (may be acceptable for root CA)")
	}

	// Test certificate has appropriate key usage for TLS
	keyUsage := cert.Certificate.KeyUsage
	if keyUsage != nil {
		if keyUsage&tls.KeyUsageDigitalSignature == 0 {
			t.Error("Certificate should support digital signatures")
		}
		if keyUsage&tls.KeyUsageKeyEncipherment == 0 {
			t.Error("Certificate should support key encipherment")
		}
	}

	// Test certificate has appropriate extended key usage for TLS
	extKeyUsage := cert.Certificate.ExtKeyUsage
	if extKeyUsage != nil {
		if extKeyUsage&tls.ExtKeyUsageServerAuth == 0 {
			t.Error("Certificate should support server authentication")
		}
	}

	// Test certificate has appropriate basic constraints
	basicConstraints := cert.Certificate.BasicConstraints
	if basicConstraints != nil {
		if basicConstraints.IsCA {
			t.Error("Certificate should not be a CA certificate for server use")
		}
	}
}

// TestCertificateReload tests certificate reloading functionality
func TestCertificateReload(t *testing.T) {
	// Load initial certificates
	cert1, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load initial certificates: %v", err)
	}

	// Create store
	store := New(cert1)

	// Load certificates again
	cert2, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to reload certificates: %v", err)
	}

	// Create new store with reloaded certificates
	store2 := New(cert2)

	// Verify both stores have the same certificate
	if store.GetCertificate() != store2.GetCertificate() {
		t.Error("Reloaded certificate should match original")
	}
}

// TestCertificateMemoryUsage tests memory usage of certificate operations
func TestCertificateMemoryUsage(t *testing.T) {
	// Load certificates
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Create multiple stores to test memory usage
	stores := make([]*CertificateStore, 100)
	for i := 0; i < 100; i++ {
		stores[i] = New(cert)
	}

	// Test all stores can retrieve certificates
	for i, store := range stores {
		retrievedCert := store.GetCertificate()
		if retrievedCert == nil {
			t.Errorf("Store %d should have valid certificate", i)
		}
		if retrievedCert != cert {
			t.Errorf("Store %d certificate should match original", i)
		}
	}
}

// TestCertificateThreadSafety tests thread safety of certificate operations
func TestCertificateThreadSafety(t *testing.T) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	// Test concurrent certificate retrieval
	done := make(chan bool, 50)
	for i := 0; i < 50; i++ {
		go func(id int) {
			retrievedCert := store.GetCertificate()
			if retrievedCert == nil {
				t.Errorf("Thread %d: Retrieved certificate should not be nil", id)
			}
			if retrievedCert != cert {
				t.Errorf("Thread %d: Retrieved certificate should match original", id)
			}
			done <- true
		}(i)
	}

	// Wait for all threads to complete
	for i := 0; i < 50; i++ {
		<-done
	}
}

// BenchmarkCertificateLoad benchmarks certificate loading performance
func BenchmarkCertificateLoad(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Load("certs/server.crt", "certs/server.key")
	}
}

// BenchmarkCertificateNew benchmarks certificate store creation
func BenchmarkCertificateNew(b *testing.B) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(cert)
	}
}

// BenchmarkGetCertificate benchmarks certificate retrieval
func BenchmarkGetCertificate(b *testing.B) {
	cert, err := Load("certs/server.crt", "certs/server.key")
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetCertificate()
	}
}

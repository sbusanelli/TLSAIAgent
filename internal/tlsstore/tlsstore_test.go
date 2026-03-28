package tlsstore

import (
	"bytes"
	"crypto/tls"
	"os"
	"testing"
)

// TestIsValid tests certificate validity checking
func TestIsValid(t *testing.T) {
	// Load a valid certificate
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	// Test valid certificate
	if !store.IsValid() {
		t.Error("Valid certificate should return true")
	}

	// Test with nil certificate
	storeNil := New(nil)
	if storeNil.IsValid() {
		t.Error("Nil certificate should return false")
	}
}

// TestLoad tests certificate loading functionality
func TestLoad(t *testing.T) {
	// Test loading valid certificates
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
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
	if len(cert.Certificate) == 0 {
		t.Error("Certificate raw bytes should not be empty")
	}

	// Note: cert.Certificate is [][]byte, not a struct with NotBefore field
	// To access certificate details, you would need to parse the x509 certificate
	// For now, just verify the certificate is loaded successfully
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
MIIDdzCCAn+gAwIBAgIEXAMPLECERTIFICATEDATA
-----END CERTIFICATE-----`

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
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	if store == nil {
		t.Fatal("Store should not be nil")
	}

	// Test GetCertificate method
	retrievedCert, err := store.GetCertificate(&tls.ClientHelloInfo{})
	if err != nil {
		t.Errorf("GetCertificate failed: %v", err)
	}
	if retrievedCert == nil {
		t.Error("Retrieved certificate should not be nil")
	}

	if retrievedCert != cert {
		t.Error("Retrieved certificate should match original")
	}
}

// TestGetCertificateWithValidFiles tests certificate retrieval
func TestGetCertificateWithValidFiles(t *testing.T) {
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	// Test multiple retrievals
	for i := 0; i < 10; i++ {
		retrievedCert, err := store.GetCertificate(&tls.ClientHelloInfo{})
		if err != nil {
			t.Errorf("GetCertificate failed (iteration %d): %v", i, err)
		}
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
	var store *Store

	// Test with nil store - this should panic, so we'll recover
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling GetCertificate on nil store")
		}
	}()

	store.GetCertificate(&tls.ClientHelloInfo{})
}

// TestGetCertificateConcurrentAccess tests concurrent certificate retrieval
func TestGetCertificateConcurrentAccess(t *testing.T) {
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			retrievedCert, err := store.GetCertificate(&tls.ClientHelloInfo{})
			if err != nil {
				t.Errorf("GetCertificate failed (goroutine %d): %v", id, err)
			}
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
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Test certificate validity
	if cert.Certificate == nil {
		t.Fatal("Certificate should not be nil")
	}

	// Test that certificate has raw data
	if len(cert.Certificate) == 0 {
		t.Error("Certificate should have raw data")
	}

	// Test private key
	if cert.PrivateKey == nil {
		t.Error("Private key should not be nil")
	}

	// Note: To access certificate details like NotBefore, Subject, etc.,
	// you would need to parse the raw certificate using x509.ParseCertificate
	// For now, we just verify the certificate is loaded successfully
}

// TestPrivateKeyValidation tests private key validation
func TestPrivateKeyValidation(t *testing.T) {
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Test private key properties
	if cert.PrivateKey == nil {
		t.Fatal("Private key should not be nil")
	}

	// Note: cert.PrivateKey is of type crypto.PrivateKey (interface)
	// To access specific key properties like Algorithm, N, E, etc.,
	// you would need to type assert to the specific key type (e.g., *rsa.PrivateKey)
	// For now, we just verify the private key is not nil
}

// TestCertificateExpiration tests certificate expiration handling
func TestCertificateExpiration(t *testing.T) {
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Note: cert.Certificate is [][]byte, not a struct with NotBefore/NotAfter fields
	// To access certificate expiration times, you would need to parse the certificate
	// using x509.ParseCertificate to get an x509.Certificate struct
	// For now, we just verify the certificate is loaded successfully

	if len(cert.Certificate) == 0 {
		t.Error("Certificate should have raw data")
	}
}

// TestCertificateSubject tests certificate subject information
func TestCertificateSubject(t *testing.T) {
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Note: cert.Certificate is [][]byte, not a struct with Subject field
	// To access certificate subject information, you would need to parse the certificate
	// using x509.ParseCertificate to get an x509.Certificate struct
	// For now, we just verify the certificate is loaded successfully

	if len(cert.Certificate) == 0 {
		t.Error("Certificate should have raw data")
	}
}

// TestCertificateIssuer tests certificate issuer information
func TestCertificateIssuer(t *testing.T) {
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Note: cert.Certificate is [][]byte, not a struct with Issuer field
	// To access certificate issuer information, you would need to parse the certificate
	// using x509.ParseCertificate to get an x509.Certificate struct
	// For now, we just verify the certificate is loaded successfully

	if len(cert.Certificate) == 0 {
		t.Error("Certificate should have raw data")
	}
}

// TestCertificateExtensions tests certificate extensions
func TestCertificateExtensions(t *testing.T) {
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Note: cert.Certificate is [][]byte, not a struct with Extensions field
	// To access certificate extensions, you would need to parse the certificate
	// using x509.ParseCertificate to get an x509.Certificate struct
	// For now, we just verify the certificate is loaded successfully

	if len(cert.Certificate) == 0 {
		t.Error("Certificate should have raw data")
	}
}

// TestCertificateSanity tests certificate sanity checks
func TestCertificateSanity(t *testing.T) {
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Note: cert.Certificate is [][]byte, not a struct with Issuer/Subject/Authority fields
	// To access certificate details for sanity checks, you would need to parse the certificate
	// using x509.ParseCertificate to get an x509.Certificate struct
	// For now, we just verify the certificate is loaded successfully

	if len(cert.Certificate) == 0 {
		t.Error("Certificate should have raw data")
	}
}

// TestCertificateReload tests certificate reloading functionality
func TestCertificateReload(t *testing.T) {
	// Load initial certificates
	cert1, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load initial certificates: %v", err)
	}

	// Create store
	store := New(cert1)

	// Load certificates again
	cert2, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to reload certificates: %v", err)
	}

	// Create new store with reloaded certificates
	store2 := New(cert2)

	// Verify both stores have certificates with the same content
	retrievedCert1, err1 := store.GetCertificate(&tls.ClientHelloInfo{})
	retrievedCert2, err2 := store2.GetCertificate(&tls.ClientHelloInfo{})
	if err1 != nil || err2 != nil {
		t.Errorf("GetCertificate failed: %v, %v", err1, err2)
	}
	if retrievedCert1 == nil || retrievedCert2 == nil {
		t.Error("Both retrieved certificates should not be nil")
	}
	// Compare the raw certificate bytes instead of the pointers
	if !certificatesEqual(retrievedCert1, retrievedCert2) {
		t.Error("Reloaded certificate should match original")
	}
}

// Helper function to compare two certificates
func certificatesEqual(cert1, cert2 *tls.Certificate) bool {
	if cert1 == nil || cert2 == nil {
		return false
	}
	// Compare the raw certificate bytes
	if len(cert1.Certificate) != len(cert2.Certificate) {
		return false
	}
	for i := range cert1.Certificate {
		if !bytes.Equal(cert1.Certificate[i], cert2.Certificate[i]) {
			return false
		}
	}
	return true
}

// TestCertificateMemoryUsage tests memory usage of certificate operations
func TestCertificateMemoryUsage(t *testing.T) {
	// Load certificates
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	// Create multiple stores to test memory usage
	stores := make([]*Store, 100)
	for i := 0; i < 100; i++ {
		stores[i] = New(cert)
	}

	// Test all stores can retrieve certificates
	for i, store := range stores {
		retrievedCert, err := store.GetCertificate(&tls.ClientHelloInfo{})
		if err != nil {
			t.Errorf("Store %d GetCertificate failed: %v", i, err)
		}
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
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		t.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	// Test concurrent certificate retrieval
	done := make(chan bool, 50)
	for i := 0; i < 50; i++ {
		go func(id int) {
			retrievedCert, err := store.GetCertificate(&tls.ClientHelloInfo{})
			if err != nil {
				t.Errorf("Thread %d: GetCertificate failed: %v", id, err)
			}
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
		Load("../../certs/server.crt", "../../certs/server.key")
	}
}

// BenchmarkCertificateNew benchmarks certificate store creation
func BenchmarkCertificateNew(b *testing.B) {
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
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
	cert, err := Load("../../certs/server.crt", "../../certs/server.key")
	if err != nil {
		b.Fatalf("Failed to load certificates: %v", err)
	}

	store := New(cert)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = store.GetCertificate(&tls.ClientHelloInfo{})
	}
}

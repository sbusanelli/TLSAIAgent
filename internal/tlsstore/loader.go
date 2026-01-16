package tlsstore

import "crypto/tls"

func Load(certFile, keyFile string) (*tls.Certificate, error) {
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return nil, err
    }
    return &cert, nil
}

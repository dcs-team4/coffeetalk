package server

import (
	"crypto/tls"
	"embed"
	"net/http"
)

// For production environment: expects tls-cert.pem and tls-key.pem in tls directory.
//go:embed all:tls
var tlsFiles embed.FS

// Listens and serves on the given address, using embedded TLS certificate and key files.
// If handler is nil, uses http.DefaultServeMux.
func listenAndServeTLS(address string, handler http.Handler) error {
	certFile, err := tlsFiles.ReadFile("tls/tls-cert.pem")
	if err != nil {
		return err
	}

	keyFile, err := tlsFiles.ReadFile("tls/tls-key.pem")
	if err != nil {
		return err
	}

	certificate, err := tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:      address,
		Handler:   handler,
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{certificate}},
	}

	return server.ListenAndServeTLS("", "")
}

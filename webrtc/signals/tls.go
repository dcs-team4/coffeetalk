package signals

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
// Limited, alternate implementation of http.ListenAndServeTLS in order to use embedded files.
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

	tlsConfig := tls.Config{
		NextProtos:   []string{"http/1.1"},
		Certificates: []tls.Certificate{certificate},
	}

	server := &http.Server{Addr: address, Handler: handler, TLSConfig: &tlsConfig}
	return server.ListenAndServeTLS("", "")
}

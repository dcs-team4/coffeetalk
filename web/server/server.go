package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

// Serves a website from the given static and template directories on the given address.
// Keeps running until an error occurs. Uses TLS if in production environment.
func Serve(address string, staticDir fs.FS, templatesDir fs.FS) error {
	err := setupRoutes(staticDir, templatesDir)
	if err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	if os.Getenv("ENV") == "production" {
		return listenAndServeTLS(address, nil)
	} else {
		return http.ListenAndServe(address, nil)
	}
}

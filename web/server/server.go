package server

import (
	"io/fs"
	"log"
	"net/http"
	"os"
)

// Serves a web app from the given static and template directories on the given address.
func Serve(address string, staticDir fs.FS, templatesDir fs.FS) error {
	err := setupRoutes(staticDir, templatesDir)
	if err != nil {
		log.Fatalf("Error setting up routes: %v\n", err)
	}

	// Runs the web server (with TLS in if in production) until an error occurs.
	if os.Getenv("ENV") == "production" {
		return listenAndServeTLS(address, nil)
	} else {
		return http.ListenAndServe(address, nil)
	}
}

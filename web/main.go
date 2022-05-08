package main

import (
	"embed"
	"log"
	"os"

	"github.com/dcs-team4/coffeetalk/web/server"
)

// Embeds the "templates" directory in this Go binary, for processing and serving of HTML.
//go:embed templates
var templatesDir embed.FS

// Embeds the "static" directory in this Go binary, for serving of JS and CSS files.
//go:embed static
var staticDir embed.FS

func main() {
	// Gets PORT environment variable, defaulting to 3000 if not present.
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}

	log.Printf("Web server listening on port %v...\n", port)

	// Serves the web app from the templates and static folders, until an error occurs.
	err := server.Serve(":"+port, staticDir, templatesDir)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Web server closed.")
}

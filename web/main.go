package main

import (
	"embed"
	"log"
	"os"

	"github.com/dcs-team4/coffeetalk/web/server"
)

// Embeds the "templates" directory, for processing and serving HTML templates.
//go:embed templates
var templatesDir embed.FS

// Embeds the "static" directory, for serving JS and CSS files.
// Excludes directories starting with _, such as _types, which we do not want to serve.
//go:embed static
var staticDir embed.FS

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}

	log.Printf("Web server listening on port %v...\n", port)

	err := server.Serve(":"+port, staticDir, templatesDir)
	if err != nil {
		log.Println(err)
	}

	log.Println("Web server closed.")
}

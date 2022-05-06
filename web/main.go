package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
)

// Embeds the "public" folder in this Go binary, for easy serving of the website.
//go:embed public
var publicFolder embed.FS

func main() {
	// Gets PORT environment variable, defaulting to 3000 if not present.
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}

	// Serves the HTML/CSS/JavaScript files inside the "public" folder from the root URL.
	publicFiles, _ := fs.Sub(publicFolder, "public")
	http.Handle("/", http.FileServer(http.FS(publicFiles)))

	// Runs the web server (with TLS in if in production) until an error occurs.
	var err error
	if os.Getenv("ENV") == "production" {
		err = listenAndServeTLS(":"+port, nil)
	} else {
		err = http.ListenAndServe(":"+port, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Web server closed.")
}

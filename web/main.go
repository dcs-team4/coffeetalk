package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

const webServerPort = "3000"

// Embeds the "public" folder in this Go binary, for easy serving of the website.
//go:embed public
var publicFolder embed.FS

func main() {
	// Serves the HTML/CSS/JavaScript files inside the "public" folder from the root URL.
	publicFiles, _ := fs.Sub(publicFolder, "public")
	http.Handle("/", http.FileServer(http.FS(publicFiles)))

	err := http.ListenAndServe(":"+webServerPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}

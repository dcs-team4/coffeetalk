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
	// Exposes the files inside the "public" folder.
	public, err := fs.Sub(publicFolder, "public")
	if err != nil {
		log.Fatal(err)
	}

	// Serves the files inside the "public" folder from the root URL.
	http.Handle("/", http.FileServer(http.FS(public)))

	err = http.ListenAndServe(":"+webServerPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}

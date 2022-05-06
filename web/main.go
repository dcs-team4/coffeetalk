package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Gets PORT environment variable, defaulting to 3000 if not present.
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}

	err := setupRoutes()
	if err != nil {
		log.Fatalf("Error setting up routes: %v\n", err)
	}

	// Runs the web server (with TLS in if in production) until an error occurs.
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

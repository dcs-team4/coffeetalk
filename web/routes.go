package main

import (
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

// Configures which HTML template files correspond to which routes on the website.
var filesToRoutes = map[string]string{
	"index.html":  "/",
	"office.html": "/office",
}

// Embeds the "templates" folder in this Go binary, for easy processing and serving of HTML.
//go:embed templates
var templatesFolder embed.FS

const templatesFolderName = "templates"

// Embeds the "static" folder in this Go binary, for easy serving JS and CSS files.
//go:embed static
var staticFolder embed.FS

func setupRoutes() error {
	// Serves the HTML/CSS/JavaScript files from inside the "static" folder.
	http.Handle("/static/", http.FileServer(http.FS(staticFolder)))

	// Reads template files from the "templates" folder.
	templates, err := fs.ReadDir(templatesFolder, templatesFolderName)
	if err != nil {
		return err
	}

	// For each template, registers an HTTP handler for serving the template from the appropriate
	// route.
	for _, templateFile := range templates {
		templateName := templateFile.Name()

		parsedTemplate, err := template.ParseFS(
			templatesFolder, templatesFolderName+"/"+templateName,
		)
		if err != nil {
			return err
		}

		route, ok := filesToRoutes[templateName]
		if !ok {
			return errors.New("invalid route config")
		}

		http.HandleFunc(route, func(res http.ResponseWriter, req *http.Request) {
			err := parsedTemplate.Execute(res, frontendEnv)
			if err != nil {
				log.Printf("Template execution failed: %v\n", err)
			}

			log.Printf("Request to %v handled.\n", route)
		})
	}

	return nil
}

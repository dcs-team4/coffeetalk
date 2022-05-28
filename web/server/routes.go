package server

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

// Configures which HTML template files correspond to which routes on the website,
// and the client type to pass to the web app as an environment variable.
var routeConfig = map[string]struct {
	route      string
	clientType string
}{
	"index.html": {
		route:      "/",
		clientType: "home",
	},
	"office.html": {
		route:      "/office",
		clientType: "office",
	},
}

// Names of served static and template directories.
const (
	staticDirName    = "static"
	templatesDirName = "templates"
)

// Sets up HTTP handler for the given static directory,
// and adds handlers for templates as configured in routeConfig.
func setupRoutes(staticDir fs.FS, templatesDir fs.FS) error {
	http.Handle("/"+staticDirName+"/", http.FileServer(http.FS(staticDir)))

	templates, err := fs.ReadDir(templatesDir, templatesDirName)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	// For each template, registers HTTP handler for serving it from the configured route.
	for _, templateFile := range templates {
		templateName := templateFile.Name()

		parsedTemplate, err := template.ParseFS(templatesDir, templatesDirName+"/"+templateName)
		if err != nil {
			return fmt.Errorf("failed to parse template %v: %w", templateName, err)
		}

		config, ok := routeConfig[templateName]
		if !ok {
			return fmt.Errorf("missing route config for template %v", templateName)
		}

		http.HandleFunc(config.route, func(res http.ResponseWriter, req *http.Request) {
			log.Println("Received request to route:", config.route)

			err := parsedTemplate.Execute(res, makeClientEnv(config.clientType))
			if err != nil {
				log.Println("Template execution failed:", err)
			}
		})
	}

	return nil
}

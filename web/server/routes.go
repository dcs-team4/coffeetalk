package server

import (
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

// Configures which HTML template files correspond to which routes on the website, and the client
// type to pass to the web app as an environment variable.
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

// Sets up HTTP handler for the given static directory (expects it to be named "static"),
// and adds handlers for templates as configured in routeConfig.
func setupRoutes(staticDir fs.FS, templatesDir fs.FS) error {
	// Serves the HTML/CSS/JavaScript files from inside the "static" folder.
	http.Handle("/"+staticDirName+"/", http.FileServer(http.FS(staticDir)))

	// Reads template files from the "templates" folder.
	templates, err := fs.ReadDir(templatesDir, templatesDirName)
	if err != nil {
		return err
	}

	// For each template, registers HTTP handler for serving it from the configured route.
	for _, templateFile := range templates {
		templateName := templateFile.Name()

		parsedTemplate, err := template.ParseFS(
			templatesDir, templatesDirName+"/"+templateName,
		)
		if err != nil {
			return err
		}

		config, ok := routeConfig[templateName]
		if !ok {
			return errors.New("invalid route config")
		}

		http.HandleFunc(config.route, func(res http.ResponseWriter, req *http.Request) {
			err := parsedTemplate.Execute(res, makeClientEnv(config.clientType))
			if err != nil {
				log.Printf("Template execution failed: %v\n", err)
			}

			log.Printf("Request to %v handled.\n", config.route)
		})
	}

	return nil
}

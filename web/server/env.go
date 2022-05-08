package server

import (
	"os"
)

// Client environment variables configured on the server.
const (
	// Type of web app: "home" / "office"
	envClientType = "CLIENT_TYPE"
)

// Names of server environment variables that should be passed on to the web application.
var forwardedEnvs = []string{
	"ENV",
	"WEBRTC_HOST",
	"WEBRTC_PORT",
	"MQTT_HOST",
	"MQTT_PORT",
}

// Environment variables passed to every web app client type.
var baseClientEnv map[string]string = makeBaseClientEnv(forwardedEnvs)

// Goes through the given envNames, checks if an environment variable is defined for each,
// and if it is, adds it to the returned map.
func makeBaseClientEnv(envNames []string) map[string]string {
	env := make(map[string]string)

	for _, envName := range envNames {
		envVar, ok := os.LookupEnv(envName)
		if ok {
			env[envName] = envVar
		}
	}

	return env
}

// Returns a map of environment variables, using the variables from baseClientEnv, and adding
// extra environment variables passed as arguments.
func makeClientEnv(clientType string) map[string]string {
	env := make(map[string]string)
	for key, value := range baseClientEnv {
		env[key] = value
	}
	env[envClientType] = clientType
	return env
}

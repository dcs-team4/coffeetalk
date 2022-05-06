package main

import (
	"os"
)

// Names of environment variables passed to the client.
var frontendEnvNames = []string{
	"ENV",
	"WEBRTC_HOST",
	"WEBRTC_PORT",
	"MQTT_HOST",
	"MQTT_PORT",
}

// Environment variables passed to the frontend.
var frontendEnv map[string]string = prepareEnv(frontendEnvNames)

// Goes through the given envNames, checks if an environment variable is defined for each,
// and if it is, adds it to the returned map.
func prepareEnv(envNames []string) map[string]string {
	env := make(map[string]string)

	for _, envName := range envNames {
		envVar, ok := os.LookupEnv(envName)
		if ok {
			env[envName] = envVar
		}
	}

	return env
}

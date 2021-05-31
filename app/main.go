package main

import (
	internal "jirm.cz/traefik-basic-auth-manager/internal"
)

func main() {

	// Parse options
	config := internal.NewGlobalConfig()

	// Setup logger
	log := internal.NewDefaultLogger()

	internal.InitDB()

	// Start
	//log.WithField("config", config).Debug("Starting with config")
	log.Info("Starting Traefik Basic Auth Manager")
	log.Infof("Listening on :%d", config.Webserver.Port)

	internal.ApiServer()
}

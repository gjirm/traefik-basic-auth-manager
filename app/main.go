package main

import (
	"github.com/sirupsen/logrus"
	internal "jirm.cz/traefik-basic-auth-manager/internal"
)

var (
	version string
	commit  string
	date    string
)

func main() {

	// Parse options
	internal.NewGlobalConfig()

	// Setup logger
	log := internal.NewDefaultLogger()

	// Start
	//log.WithField("config", config).Debug("Starting with config")
	log.WithFields(logrus.Fields{
		"version":    version,
		"commitHash": commit,
		"BuildDate":  date,
	}).Info("Starting Traefik Basic Auth Manager")

	// Init NutsDB
	internal.InitDB()

	// Start webserver
	internal.ApiServer()
}

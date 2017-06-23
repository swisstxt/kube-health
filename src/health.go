// kube-health
//
// Configurable health checker with built-in web server
//
// Supports a single argument, the location of the configuration file.
// Will read /etc/kube-health/config.json if run without arguments.

package main

import (
	"os"
	"log"
	"health"
	"net/http"
)

const (
	defaultConfigName = "/etc/kube-health/config.json"
)

func main() {
	var configname string
	if len(os.Args) > 1 {
		configname = os.Args[1]
	} else {
		configname = defaultConfigName
	}
	
	config, err := health.LoadConfiguration(configname)
	if err != nil {
		log.Fatal("Error loading configuration: ", err.Error())
	}
	
	server := &health.Server{Config: config}
	log.Fatal(http.ListenAndServe(config.Listen, server))
}
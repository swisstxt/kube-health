// kube-health
//
// Configurable health checker with built-in web server
//
// Supports a single argument, the location of the configuration file.
// Will read /etc/kube-health/config.json if run without arguments.

package main

import (
	"log"
	"net/http"
	"os"
)

const (
	defaultConfigName = "/etc/kube-health/config.json"
)

var (
	version  string = ""
	revision string = ""
)

func main() {
	var configname string
	if len(os.Args) > 1 {
		configname = os.Args[1]
	} else {
		configname = defaultConfigName
	}

	vstring := ""
	if version != "" {
		vstring += " version " + version
	}
	if revision != "" {
		vstring += " revision " + revision
	}
	log.Print("kube-health" + vstring)
	log.Print("Copyright © 2017-2018 Gregor Riepl")

	config, err := LoadConfiguration(configname)
	if err != nil {
		log.Fatal("Error loading configuration: ", err.Error())
	}

	server := &Server{Config: config}
	log.Fatal(http.ListenAndServe(config.Listen, server))
}

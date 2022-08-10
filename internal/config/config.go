package config

import (
	"flag"
	"log"
	"os"
	"strconv"
)

// Configuration variables.
var (
	InProduction       bool   // set if running in production or not
	KubeConfigLocation string // location of kubeconfig relative to app
	WebServerPath      string // path webserver will be serving on
	WebServerUrl       string // url (including http(s)://) web server will be hosted on
)

func inProduction() (bool, string, string, string) {
	var appInProduction bool
	var kubeConfig string
	var err error

	a := os.Getenv("IN_PRODUCTION")
	appInProduction = false // default to false
	if len(a) > 0 {
		appInProduction, err = strconv.ParseBool(a) // parse value
		if err != nil {
			panic(err.Error())
		}
	}
	if appInProduction {
		log.Printf("App set to run in PRODUCTION mode")
	} else {
		log.Printf("App set to run in DEVELOPMENT mode")
	}

	kubeConfig = os.Getenv("KUBE_CONFIG_LOCATION")
	if len(kubeConfig) == 0 {
		kubeConfig = "./test-kubeconfig" // set this to local config to test with
	}
	log.Printf("App set to use config located at: %s", kubeConfig)

	webServerPath := os.Getenv("WEB_SERVER_PATH")
	if len(webServerPath) == 0 {
		webServerPath = "/kube-view" // note dockerfile ENV variable set to this path
	}
	log.Printf("App web server set to use path: %s", webServerPath)

	webServerUrl := os.Getenv("WEB_SERVER_URL")
	if len(webServerUrl) == 0 {
		webServerUrl = "http://localhost:8080" // note dockerfile ENV variable set to this path
	}
	log.Printf("App web server hosted on address: %s", webServerUrl)

	return appInProduction, kubeConfig, webServerPath, webServerUrl
}

// Set configuration variables from os.Args
func Set() {
	appInProduction, kubeConfig, webServerPath, webServerUrl := inProduction()

	flag.BoolVar(&InProduction, "inProduction", appInProduction, "Set if app in production environment.")
	flag.StringVar(&KubeConfigLocation, "kubeConfig", kubeConfig, "Location of kube config to import.")
	flag.StringVar(&WebServerPath, "webServerPath", webServerPath, "Path web server with serve on")
	flag.StringVar(&WebServerUrl, "webServerUrl", webServerUrl, "url (including http(s)://) web server will be hosted on")
}

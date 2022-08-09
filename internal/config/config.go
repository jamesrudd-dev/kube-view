package config

import (
	"flag"
	"os"
	"strconv"
)

// Configuration variables.
var (
	InProduction       bool   // set if running in production or not
	KubeConfigLocation string // location of kubeconfig relative to app
)

func inProduction() (bool, string) {
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
		kubeConfig = os.Getenv("KUBE_CONFIG_LOCATION")
	} else {
		kubeConfig = "./test-kubeconfig" // set this to local config to test with
	}

	return appInProduction, kubeConfig
}

// Set configuration variables from os.Args
func Set() {
	appInProduction, kubeConfig := inProduction()

	flag.BoolVar(&InProduction, "inProduction", appInProduction, "Set if app in production environment.")
	flag.StringVar(&KubeConfigLocation, "kubeConfig", kubeConfig, "Location of kube config to import.")
}

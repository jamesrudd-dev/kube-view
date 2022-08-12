package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"k8s.io/client-go/util/homedir"
)

// Configuration variables.
var (
	InProduction       bool   // set if running in production or not
	KubeConfigLocation string // location of kubeconfig relative to app
	WebServerPath      string // path webserver will be serving on
	ImageTagFilter     string // set this to clean up images tags to remove (for example) the AWS prefix, i.e. ###########.dkr.ecr.ap-southeast-2.amazonaws.com/, currently only one filter available
	NamespaceFilter    string // comma seperated list containing namespaces desired to be removed from search, note: strings.Contains used for filter
)

// pullEnvVars will set variables from environmental variables (if exist)
// otherwise set default variables. These are passed to Set() function below.
func pullEnvVars() (bool, string, string, string, string) {
	var inProduction bool
	var err error

	inProduction = false // default to false
	a := os.Getenv("IN_PRODUCTION")
	if len(a) > 0 {
		inProduction, err = strconv.ParseBool(a) // parse value
		if err != nil {
			panic(errors.New("CONFIG - config.inProduction: failed to parse string as boolean"))
		}
		log.Printf("App set to run in PRODUCTION mode")
	}
	if !inProduction {
		log.Printf("App set to run in DEVELOPMENT mode")
	}

	kubeConfigLocation := os.Getenv("KUBE_CONFIG_LOCATION")
	if len(kubeConfigLocation) == 0 {
		kubeConfigLocation = filepath.Join(homedir.HomeDir(), ".kube", "config") // defaults to home directory kube config (for outside container dev)
	}
	log.Printf("App set to use config located at: %s", kubeConfigLocation)

	webServerPath := os.Getenv("WEB_SERVER_PATH")
	if len(webServerPath) == 0 {
		webServerPath = "/kube-view"
	}
	log.Printf("App web server set to use path: %s", webServerPath)

	imageTagFilter := os.Getenv("IMAGE_TAG_FILTER")
	if len(imageTagFilter) == 0 {
		imageTagFilter = ""
	}
	log.Printf("App filtering the following image tags: %s", imageTagFilter)

	namespaceFilter := os.Getenv("NAMESPACE_FILTER")
	if len(namespaceFilter) == 0 {
		namespaceFilter = ""
	}
	log.Printf("App filtering the following namespaces: %s", namespaceFilter)

	return inProduction, kubeConfigLocation, webServerPath, imageTagFilter, namespaceFilter
}

// Set will set the global configuration variables
func Set() {
	inProduction, kubeConfigLocation, webServerPath, imageTagFilter, namespaceFilter := pullEnvVars()

	flag.BoolVar(&InProduction, "inProduction", inProduction, "Set if app in production environment")
	flag.StringVar(&KubeConfigLocation, "kubeConfigLocation", kubeConfigLocation, "Location of kube config to import")
	flag.StringVar(&WebServerPath, "webServerPath", webServerPath, "Path web server with serve on")
	flag.StringVar(&ImageTagFilter, "imageTagFilter", imageTagFilter, "Set this to clean up images tags to remove image prefix")
	flag.StringVar(&NamespaceFilter, "namespaceFilter", namespaceFilter, "Comma seperated list containing namespaces desired to be removed from search")
}

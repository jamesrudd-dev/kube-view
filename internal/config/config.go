package config

type AppConfig struct {
	InCluster          bool   // true if app deployed inside Kubernetes cluster
	KubeConfigLocation string // location of kubeconfig relative to app
	HttpServer         bool   // true if api HTTP server desired to be run
}

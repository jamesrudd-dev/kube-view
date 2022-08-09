package main

import (
	"jamesrudd-dev/kube-view/internal/api"
	"jamesrudd-dev/kube-view/internal/database"
	"jamesrudd-dev/kube-view/internal/handlers"
	"net/http"
	"time"
)

func runHttpServer() {
	// define server handler and runtime options
	srv := &http.Server{
		Handler:      routes(),
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// run server
	srv.ListenAndServe()
}

func main() {

	// set kubeconfig
	clientSet, err := handlers.SetKubeConfig()
	if err != nil {
		panic(err.Error())
	}

	// read kubeconfig for clusterList
	clusterList, err := handlers.ReadConfig("/home/jrudd/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	// connect to redis (base database - epe-kubernetes)
	redisClient, err := database.InitialConnectRedis()
	if err != nil {
		panic(err.Error())
	}

	// pass config and database connection config to api's
	api.SetKubeConfig(clientSet)
	api.SetDatabase(redisClient)
	api.SetClusterList(clusterList)

	// do initial scrape of epe-kubernetes to test config and database connection
	err = handlers.ScrapeKubernetes(clientSet, redisClient)
	if err != nil {
		panic(err.Error())
	}

	runHttpServer()

	defer database.CloseRedis()

}

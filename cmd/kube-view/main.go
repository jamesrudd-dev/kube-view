package main

import (
	"errors"
	"jamesrudd-dev/kube-view/internal/api"
	"jamesrudd-dev/kube-view/internal/config"
	"jamesrudd-dev/kube-view/internal/database"
	"jamesrudd-dev/kube-view/internal/handlers"
	"log"
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

	// Set OS Arguments
	config.Set()

	// read kubeconfig for clusterList
	clusterList, err := handlers.ReadConfig(config.KubeConfigLocation)
	if err != nil {
		err = errors.New("handlers.ReadConfig: failed to read kube config from given location")
		panic(err)
	}

	// set kubeconfig
	clientSet, clusterDatabase, err := handlers.SetKubeConfig(config.KubeConfigLocation, clusterList)
	if err != nil {
		err = errors.New("handlers.SetKubeConfig: failed to set the initial kube config")
		panic(err)
	}

	// connect to redis (base database - epe-kubernetes)
	redisClient, err := database.InitialConnectRedis(clusterDatabase)
	if err != nil {
		err = errors.New("database.InitialConnectRedis: failed initial connection to Redis")
		panic(err)
	}

	// pass config and database connection config to api's
	api.SetKubeConfig(clientSet)
	api.SetDatabase(redisClient)
	api.SetClusterList(clusterList)

	// do initial scrape of epe-kubernetes to test config and database connection
	err = handlers.ScrapeKubernetes(clientSet, redisClient)
	if err != nil {
		err = errors.New("handlers.ScrapeKubernetes: failed initial scrape of kubernetes deployments")
		panic(err)
	}

	log.Println("Starting web server...")
	runHttpServer()

	defer database.CloseRedis()

}

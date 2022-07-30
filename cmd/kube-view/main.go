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

	clientSet, err := handlers.SetKubeConfig()
	if err != nil {
		panic(err.Error())
	}

	redisClient, err := database.ConnectRedis()
	if err != nil {
		panic(err.Error())
	}

	api.SetDatabase(redisClient)

	err = handlers.ScrapeKubernetes(clientSet, redisClient)
	if err != nil {
		panic(err.Error())
	}

	runHttpServer()

	defer database.CloseRedis()

}

package main

import (
	"jamesrudd-dev/kube-view/internal/api"
	"jamesrudd-dev/kube-view/internal/config"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func routes() http.Handler {
	if config.InProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// home page
	router.Use(static.Serve(config.WebServerPath, static.LocalFile("./frontend/build", true)))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.WebServerUrl},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:8080"
		},
		MaxAge: 12 * time.Hour,
	}))

	// create api routes
	router.GET("/deployments/:cluster/:namespace", api.GetDeploymentsfromNamespace)
	router.GET("/cluster/:cluster/namespaces", api.GetClusterNamespaces)
	router.GET("cluster/list", api.GetClusterList)
	router.POST("/cluster/:cluster/refresh", api.PostClusterRefresh)

	return router
}

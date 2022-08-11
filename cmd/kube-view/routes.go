package main

import (
	"jamesrudd-dev/kube-view/internal/api"
	"jamesrudd-dev/kube-view/internal/config"
	"net/http"

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

	// cors
	router.Use(cors.Default())

	// create api routes
	router.GET("/kube-view/deployments/:cluster/:namespace", api.GetDeploymentsfromNamespace)
	router.GET("/kube-view/cluster/:cluster/namespaces", api.GetClusterNamespaces)
	router.GET("/kube-view/cluster/list", api.GetClusterList)
	router.POST("/kube-view/cluster/:cluster/refresh", api.PostClusterRefresh)

	return router
}

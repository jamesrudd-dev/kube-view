package main

import (
	"jamesrudd-dev/kube-view/internal/api"
	"jamesrudd-dev/kube-view/internal/config"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// routes set the http handlers of the web server
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
	router.GET(config.WebServerPath+"/deployments/:cluster/:namespace", api.GetDeploymentsfromNamespace)
	router.GET(config.WebServerPath+"/cluster/:cluster/namespaces", api.GetClusterNamespaces)
	router.GET(config.WebServerPath+"/cluster/list", api.GetClusterList)
	router.POST(config.WebServerPath+"/cluster/:cluster/refresh", api.PostClusterRefresh)

	return router
}

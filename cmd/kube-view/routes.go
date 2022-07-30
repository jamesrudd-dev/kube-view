package main

import (
	"jamesrudd-dev/kube-view/internal/api"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func routes() http.Handler {
	router := gin.Default()

	// home page
	router.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))

	// create api routes
	router.GET("/deployments/:namespace", api.GetDeploymentsfromNamespace)
	router.GET("/cluster/namespaces", api.GetClusterNamespaces)
	router.POST("/cluster/refresh", api.PostClusterRefresh)
	router.POST("/cluster-change/:context", api.PostChangeClusterContext)

	return router
}

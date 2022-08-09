package main

import (
	"jamesrudd-dev/kube-view/internal/api"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func routes() http.Handler {
	router := gin.Default()

	// home page
	// router.Use(static.Serve("/", static.LocalFile("./frontend/build", true)))
	router.Use(cors.Default())

	// create api routes
	router.GET("/deployments/:cluster/:namespace", api.GetDeploymentsfromNamespace)
	router.GET("/cluster/:cluster/namespaces", api.GetClusterNamespaces)
	router.GET("cluster/list", api.GetClusterList)
	router.POST("/cluster/:cluster/refresh", api.PostClusterRefresh)

	return router
}

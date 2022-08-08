package api

import (
	"context"
	"encoding/json"
	"fmt"
	"jamesrudd-dev/kube-view/internal/database"
	"jamesrudd-dev/kube-view/internal/handlers"
	"jamesrudd-dev/kube-view/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var RDB *redis.Client
var KC *kubernetes.Clientset

func SetDatabase(rdb *redis.Client) {
	RDB = rdb
}

func SetKubeConfig(config *kubernetes.Clientset) {
	KC = config
}

func GetDeploymentsfromNamespace(c *gin.Context) {
	var cursor uint64
	var n int
	var clusterDatabase int

	ns := c.Param("namespace")
	ctx := c.Param("cluster")

	if ctx == "epe-kubernetes" {
		clusterDatabase = 0
	} else if ctx == "aus-prod-kubernetes" {
		clusterDatabase = 1
	} else if ctx == "old-aus-prod-kubernetes" {
		clusterDatabase = 2
	} else if ctx == "us-uat-kubernetes" {
		clusterDatabase = 3
	} else if ctx == "us-prod-kubernetes" {
		clusterDatabase = 4
	} else if ctx == "uk-prod-kubernetes" {
		clusterDatabase = 5
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "cluster context doesn't exist",
		})
	}

	RDB, _ = database.ChangeRedisDatabase(RDB, clusterDatabase)

	config, err := handlers.SetKubeContext(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "failed to changed context",
		})
	}

	KC, err = kubernetes.NewForConfig(config)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "failed to changed context",
		})
	}

	handlers.ScrapeKubernetes(KC, RDB)

	for {
		var keys []string
		var err error
		keys, cursor, err = RDB.Scan(cursor, ns+"*", 10).Result()
		if err != nil {
			panic(err)
		}
		n += len(keys)
		if cursor == 0 {
			break
		}
	}

	if n == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "namespace not found",
		})
	}

	test := make([]models.KubernetesDeployment, n)

	for i := 0; i < n; i++ {
		val, err := RDB.Get(ns + "_" + strconv.Itoa(i)).Result()
		if err != nil {
			fmt.Println(err)
		}
		var prettyJSON models.KubernetesDeployment
		err = json.Unmarshal([]byte(val), &prettyJSON)
		if err != nil {
			panic(err.Error())
		}
		test[i] = prettyJSON
	}

	c.IndentedJSON(http.StatusOK, test)
}

func PostClusterRefresh(c *gin.Context) {
	var clusterDatabase int

	ctx := c.Param("cluster")

	if ctx == "epe-kubernetes" {
		clusterDatabase = 0
	} else if ctx == "aus-prod-kubernetes" {
		clusterDatabase = 1
	} else if ctx == "old-aus-prod-kubernetes" {
		clusterDatabase = 2
	} else if ctx == "us-uat-kubernetes" {
		clusterDatabase = 3
	} else if ctx == "us-prod-kubernetes" {
		clusterDatabase = 4
	} else if ctx == "uk-prod-kubernetes" {
		clusterDatabase = 5
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "cluster context doesn't exist",
		})
	}

	RDB, _ = database.ChangeRedisDatabase(RDB, clusterDatabase)

	err := handlers.ScrapeKubernetes(KC, RDB)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "deployment list has faled to refresh",
		})
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"success": "deployment list has successfully reset and refreshed",
	})
}

func GetClusterNamespaces(c *gin.Context) {
	ctx := c.Param("cluster")

	config, err := handlers.SetKubeContext(ctx)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "failed to changed context",
		})
	}

	KC, err = kubernetes.NewForConfig(config)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "failed to changed context",
		})
	}

	// get list of all namespaces
	nsList, err := KC.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"error": "failed to fetch namespaces in cluster",
		})
	}
	var test []string
	for _, n := range nsList.Items {
		if strings.Contains(n.Name, "kube") || n.Name == "nginx-ingress" || n.Name == "verdaccio" || n.Name == "lens-metrics" || n.Name == "monitoring" {
			continue
		}
		test = append(test, n.Name)
	}

	c.IndentedJSON(http.StatusOK, test)
}

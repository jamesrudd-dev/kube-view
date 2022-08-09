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
var CL []models.ClusterList

func SetDatabase(rdb *redis.Client) {
	RDB = rdb
}

func SetKubeConfig(config *kubernetes.Clientset) {
	KC = config
}

func SetClusterList(clusterList []models.ClusterList) {
	CL = clusterList
}

func GetDeploymentsfromNamespace(c *gin.Context) {
	var cursor uint64
	var n int
	var clusterDatabase int

	ns := c.Param("namespace")
	ctx := c.Param("cluster")

	for i := range CL {
		if strings.Contains(CL[i].Cluster, ctx) {
			clusterDatabase = i
		}
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

	deployments := make([]models.KubernetesDeployment, n)

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
		deployments[i] = prettyJSON
	}

	c.IndentedJSON(http.StatusOK, deployments)
}

func PostClusterRefresh(c *gin.Context) {
	var clusterDatabase int

	ctx := c.Param("cluster")

	for i := range CL {
		if strings.Contains(CL[i].Cluster, ctx) {
			clusterDatabase = i
		}
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

	err = handlers.ScrapeKubernetes(KC, RDB)
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
	namespaceList := make([]models.NamespaceList, len(nsList.Items))
	propsID := 0
	for _, n := range nsList.Items {
		if strings.Contains(n.Name, "kube") || n.Name == "nginx-ingress" || n.Name == "verdaccio" || n.Name == "lens-metrics" || n.Name == "monitoring" || n.Name == "linkerd" {
			namespaceList = append(namespaceList[:propsID], namespaceList[propsID+1:]...)
			continue
		}
		namespaceList[propsID].ID = propsID
		namespaceList[propsID].Namespace = n.Name
		propsID++
	}

	c.IndentedJSON(http.StatusOK, namespaceList)
}

func GetClusterList(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, CL)

}

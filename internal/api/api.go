package api

import (
	"context"
	"encoding/json"
	"jamesrudd-dev/kube-view/internal/config"
	"jamesrudd-dev/kube-view/internal/database"
	"jamesrudd-dev/kube-view/internal/handlers"
	"jamesrudd-dev/kube-view/internal/models"
	"log"
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

// Pull database connection from main package
func SetDatabase(redisClient *redis.Client) {
	RDB = redisClient
}

// Pull client set from main package
func SetKubeConfig(kubeConfig *kubernetes.Clientset) {
	KC = kubeConfig
}

// Pull cluster list from main package
func SetClusterList(clusterList []models.ClusterList) {
	CL = clusterList
}

// GetDeploymentsFromNamespace returns JSON object of deployments in current
// cluster context and database.
func GetDeploymentsfromNamespace(c *gin.Context) {
	var cursor uint64
	var n int
	var clusterDatabase int

	ns := c.Param("namespace")
	ctx := c.Param("cluster")

	// set kubeconfig from the provided cluster context
	KC, clusterDatabase, err := handlers.SetKubeConfig(config.KubeConfigLocation, ctx, CL)
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{
			"error": "failed to set kubeconfig to selected cluster context",
		})
	}

	// change database to selected cluster
	RDB, err = database.ChangeRedisDatabase(RDB, clusterDatabase)
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{
			"error": "failed to change redis database to selected cluster",
		})
	}

	// Update the deployment list for the selected cluster in the respective database
	err = handlers.ScrapeKubernetes(KC, RDB)
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{
			"error": "failed to get deployments from selected namespace",
		})
	}

	// scan redis database from keys matching selected namespace
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

	// create dynamic map with length being number of keys found above
	deployments := make([]models.KubernetesDeployment, n)

	// build JSON object containing deployments list
	for i := 0; i < n; i++ {
		val, err := RDB.Get(ns + "_" + strconv.Itoa(i)).Result()
		if err != nil {
			log.Println(err)
		}
		var prettyJSON models.KubernetesDeployment
		err = json.Unmarshal([]byte(val), &prettyJSON)
		if err != nil {
			c.IndentedJSON(http.StatusServiceUnavailable, gin.H{
				"error": "failed to get deployments from selected namespace",
			})
		}
		deployments[i] = prettyJSON
	}

	// return success 200
	c.IndentedJSON(http.StatusOK, deployments)
}

// PostClusterRefresh will flush the redis DB of current deployments
// before re-collecting data calling ScrapeKubernetes function.
func PostClusterRefresh(c *gin.Context) {
	var clusterDatabase int

	ctx := c.Param("cluster")

	// set kubeconfig from the provided cluster context
	KC, clusterDatabase, err := handlers.SetKubeConfig(config.KubeConfigLocation, ctx, CL)
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{
			"error": "failed to set kubeconfig to selected cluster context",
		})
	}

	// change database to selected cluster
	RDB, err = database.ChangeRedisDatabase(RDB, clusterDatabase)
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{
			"error": "failed to change redis database to selected cluster",
		})
	}

	// Update the deployment list for the selected cluster in the respective database
	err = handlers.ScrapeKubernetes(KC, RDB)
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{
			"error": "deployment list has faled to refresh",
		})
	}

	// return success 200
	c.IndentedJSON(http.StatusOK, gin.H{
		"success": "deployment list has successfully reset and refreshed",
	})
}

// GetClusterNamespaces returns a JSON object containing namespaces of
// selected cluster
func GetClusterNamespaces(c *gin.Context) {
	ctx := c.Param("cluster")

	// set kubeconfig from the provided cluster context
	KC, _, err := handlers.SetKubeConfig(config.KubeConfigLocation, ctx, CL)
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{
			"error": "failed to set kubeconfig to selected cluster context",
		})
	}

	// get list of all namespaces
	nsList, err := KC.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{
			"error": "failed to fetch namespaces in cluster",
		})
	}
	namespaceList := make([]models.NamespaceList, len(nsList.Items))
	filteredNamespaces := strings.Fields(strings.Replace(config.NamespaceFilter, ",", " ", -1))

	propsID := 0
	var stringMatch string
	for _, n := range nsList.Items {
		// filter to remove unwanted namespaces ## TODO - clean the logic here, need a better way to compare strings
		for i := range filteredNamespaces {
			if strings.Contains(n.Name, filteredNamespaces[i]) {
				stringMatch = "match"
			}
		}
		if stringMatch == "match" {
			stringMatch = ""
			namespaceList = append(namespaceList[:propsID], namespaceList[propsID+1:]...)
			continue
		}
		namespaceList[propsID].ID = propsID
		namespaceList[propsID].Namespace = n.Name
		propsID++
	}

	c.IndentedJSON(http.StatusOK, namespaceList)
}

// Returns the JSON object of the cluster list obtained from ReadConfig function.
func GetClusterList(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, CL)

}

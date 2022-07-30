package api

import (
	"encoding/json"
	"fmt"
	"jamesrudd-dev/kube-view/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var RDB *redis.Client

func SetDatabase(rdb *redis.Client) {
	RDB = rdb
}

func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

func GetDeploymentsfromNamespace(c *gin.Context) {
	var cursor uint64
	var n int

	ns := c.Param("namespace")

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
			"message": "namespace not found",
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
		test = append(test, prettyJSON)
	}

	c.IndentedJSON(http.StatusOK, test)
}

func PostDeploymentsRefresh(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Deployment list has successfully refreshed",
	})
}

func GetClusterNamespaces(c *gin.Context) {

}

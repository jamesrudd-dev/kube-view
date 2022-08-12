package handlers

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"jamesrudd-dev/kube-view/internal/config"
	"jamesrudd-dev/kube-view/internal/models"
	"os"
	"regexp"
	"strings"

	"github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var Kubeconfig *string

// SetKubeConfig returns the clientset config to be used with k8s api.
// It pulls in context value from above function.
func SetKubeConfig(kubeConfig string, ctx string, clusterList []models.ClusterList) (*kubernetes.Clientset, int, error) {
	var clusterDatabase int

	// no context was set (meaning likely from first call from main)
	// set appropriate flags and context to first context in clusterList
	if ctx == "" {
		Kubeconfig = flag.String("kubeconfig", kubeConfig, "absolute path to the kubeconfig file")
		flag.Parse()
		ctx = clusterList[0].Cluster
	}

	// if context was set, get database number for redis connection
	for i := range clusterList {
		if strings.Contains(clusterList[i].Cluster, ctx) {
			clusterDatabase = i
		}
	}

	// return config from context selection
	config, err := SetKubeContext(ctx)
	if err != nil {
		return nil, -1, err
	}

	// set clientset to use kubernetes config set to selected context
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, -1, err
	}

	return clientSet, clusterDatabase, nil
}

// SetKubeContext overrides the context value within the kube config - called within SetKubeConfig function
func SetKubeContext(context string) (*rest.Config, error) {
	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *Kubeconfig}
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: context}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

// ReadConfig aims to extract the cluster names from the context section
// of a kube config.
func ReadConfig(filename string) ([]models.ClusterList, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := string(b)

	r, _ := regexp.Compile("cluster: (.*-kubernetes)")
	clusters := r.FindAllString(config, -1)

	clusterList := make([]models.ClusterList, len(clusters))
	propsID := 0
	for _, n := range clusters {
		cluster := strings.Split(n, ":")
		clusterList[propsID].ID = propsID
		clusterList[propsID].Cluster = strings.TrimSpace(cluster[1])
		propsID++
	}

	return clusterList, nil
}

// ScrapeKubernetes uses the client set and redis connection to connect to the current
// context and extract deployment information to store in currently set redis database.
func ScrapeKubernetes(clientSet *kubernetes.Clientset, rdb *redis.Client) error {
	// get list of all namespaces
	nsList, err := clientSet.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return err
	}

	// clear existing database for clean read
	rdb.FlushDB()

	// range through all namespaces to get deployments per namespace
	for _, n := range nsList.Items {

		filteredNamespaces := strings.Split(config.NamespaceFilter, ",")
		for i := range filteredNamespaces {
			if strings.Contains(n.Name, filteredNamespaces[i]) {
				continue
			}
		}

		// if strings.Contains(n.Name, "kube") || n.Name == "nginx-ingress" || n.Name == "verdaccio" || n.Name == "lens-metrics" || n.Name == "monitoring" {
		// 	continue
		// }

		// get list of all deployments
		deployments, err := clientSet.AppsV1().Deployments(n.Name).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			return err
		}

		// Marshal indent gives a pretty print of json object
		a, _ := json.MarshalIndent(deployments, "", "    ")

		// make dynamic array from number of deployment
		kubeData := make([]models.KubernetesDeployment, len(deployments.Items))

		deploymentNames := (gjson.GetBytes(a, "items.#.metadata.name")).Array()
		imageNames := (gjson.GetBytes(a, "items.#.spec.template.spec.containers.0.image")).Array()

		z := 0
		for i := 0; i < len(deployments.Items); i++ {
			kubeData[i].ID = z
			kubeData[i].Namespace = n.Name
			kubeData[i].DeploymentName = (deploymentNames[i]).String()

			// image tag filter to remove any prefixes from images ##TODO - enable this to handle multiple filters
			filteredImages := imageNames[i].String()
			if config.ImageTagFilter != "" {
				filteredImages = strings.ReplaceAll((imageNames[i]).String(), config.ImageTagFilter, "")
			}
			if strings.Contains(filteredImages, ":") {
				name := (strings.Split(filteredImages, ":"))[0]
				tag := (strings.Split(filteredImages, ":"))[1]
				kubeData[i].ImageName = name
				kubeData[i].ImageTag = tag
			} else {
				kubeData[i].ImageName = filteredImages
				kubeData[i].ImageTag = "latest"
			}

			// Encode into JSON
			marshalledData, err := json.Marshal(kubeData[i])
			if err != nil {
				return err
			}

			// Commit the data to database
			err = rdb.Set(fmt.Sprintf("%s_%d", n.Name, z), marshalledData, 0).Err()
			if err != nil {
				return err
			}
			z++
		}
	}
	return nil
}

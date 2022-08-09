package handlers

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"jamesrudd-dev/kube-view/internal/models"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var Kubeconfig *string

func SetKubeContext(context string) (*rest.Config, error) {
	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *Kubeconfig}
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: context}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func SetKubeConfig() (*kubernetes.Clientset, error) {
	// pull in kubeconfig (if running outside cluster)
	if home := homedir.HomeDir(); home != "" {
		Kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		Kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := SetKubeContext("epe-kubernetes")
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

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
		clusterList[propsID].ID = propsID
		clusterList[propsID].Cluster = n
		propsID++
	}

	return clusterList, nil
}

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

		if strings.Contains(n.Name, "kube") || n.Name == "nginx-ingress" || n.Name == "verdaccio" || n.Name == "lens-metrics" || n.Name == "monitoring" {
			continue
		}

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

			s := strings.ReplaceAll((imageNames[i]).String(), "635705773620.dkr.ecr.ap-southeast-2.amazonaws.com/", "")
			if strings.Contains(s, ":") {
				name := (strings.Split(s, ":"))[0]
				tag := (strings.Split(s, ":"))[1]
				kubeData[i].ImageName = name
				kubeData[i].ImageTag = tag
			} else {
				kubeData[i].ImageName = s
				kubeData[i].ImageTag = "latest"
			}

			marshalledData, err := json.Marshal(kubeData[i])
			if err != nil {
				return err
			}

			//err = rdb.Set(fmt.Sprintf("%s_%s", n.Name, kubeData[i].DeploymentName), marshalledData, 0).Err()
			err = rdb.Set(fmt.Sprintf("%s_%d", n.Name, z), marshalledData, 0).Err()
			if err != nil {
				return err
			}
			z++
		}
	}
	return nil
}

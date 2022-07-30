package handlers

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"jamesrudd-dev/kube-view/internal/models"
	"path/filepath"
	"strings"

	"github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func SetKubeConfig() (*kubernetes.Clientset, error) {
	// pull in kubeconfig (if running outside cluster)
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// creates the in-cluster config (if deployed in cluster)
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	panic(err.Error())
	// }

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}

func ScrapeKubernetes(clientSet *kubernetes.Clientset, rdb *redis.Client) error {
	// get list of all namespaces
	nsList, err := clientSet.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return err
	}

	// clear existing database for clean read
	rdb.FlushAll()

	// range through all namespaces to get deployments per namespace
	for _, n := range nsList.Items {

		if strings.Contains(n.Name, "kube") || n.Name == "nginx-ingress" || n.Name == "verdaccio" {
			continue
		}

		// get list of all deployments
		deployments, err := clientSet.AppsV1().Deployments(n.Name).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			return err
		}

		// Marshal indent gives a pretty print of json object
		a, _ := json.MarshalIndent(deployments, "", "    ")

		// println(string(a))
		// result := gjson.Parse(string(a))
		// println(result.String())

		// make dynamic array from number of deployment
		kubeData := make([]models.KubernetesDeployment, len(deployments.Items))

		deploymentNames := (gjson.GetBytes(a, "items.#.metadata.name")).Array()
		imageNames := (gjson.GetBytes(a, "items.#.spec.template.spec.containers.0.image")).Array()

		// println(deploymentNames[0].String())
		// println(imageNames[0].String())
		var key = 0
		for i := 0; i < len(deployments.Items); i++ {
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

			//kubeData[i].ImageName = (imageNames[i]).String()

			b, err := json.Marshal(kubeData[i])
			if err != nil {
				return err
			}

			//fmt.Println(string(b))

			err = rdb.Set(fmt.Sprintf("%s_%d", n.Name, key), b, 0).Err()
			key++
			if err != nil {
				return err
			}
		}
	}
	return nil
}

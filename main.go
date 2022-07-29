package main

import (
	"context"
	"encoding/json"
	"flag"
	"path/filepath"

	"github.com/tidwall/gjson"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type kubernetesDeployment struct {
	deploymentName   string
	imageName        string
	creationTime     string
	deploymentStatus string
}

func main() {

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
		panic(err.Error())
	}

	clientSet, _ := kubernetes.NewForConfig(config)

	namespace := "default"

	deployments, err := clientSet.AppsV1().Deployments(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	//a, _ := json.Marshal(deployments)
	a, _ := json.MarshalIndent(deployments, "", "    ")

	//println(string(a))
	//result := gjson.Parse(string(a))
	//println(result.String())

	result := gjson.GetBytes(a, "items.#.spec.template.spec.containers.0.image")
	println(result.String())

}

package models

type KubernetesDeployment struct {
	ID             int    `json:"id"`
	Namespace      string `json:"namespace"`
	DeploymentName string `json:"deploymentName"`
	ImageName      string `json:"imageName"`
	ImageTag       string `json:"imageTag"`
}

type NamespaceList struct {
	ID        int    `json:"id"`
	Namespace string `json:"namespace"`
}

type ClusterList struct {
	ID      int    `json:"id"`
	Cluster string `json:"cluster"`
}

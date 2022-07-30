package models

type KubernetesDeployment struct {
	Namespace      string `json:"namespace"`
	DeploymentName string `json:"deploymentName"`
	ImageName      string `json:"imageName"`
	ImageTag       string `json:"imageTag"`
}

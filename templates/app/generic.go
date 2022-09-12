package app

import (
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Enviroment struct {
	Name  string
	Value string
}

type GenericApp struct {
	Name         string
	Image        string
	Port         int
	ResourceName string
	Environments *[]Enviroment
}

func NewGenericApp(app GenericApp) ([]*unstructured.Unstructured, error) {
	deploymentStr := `
	apiVersion: apps/v1
	kind: Deployment
	metadata:
		name: nginx-deployment
		labels:
			app: nginx
	spec:
		replicas: 3
		selector:
			matchLabels:
				app: nginx
		template:
			metadata:
				labels:
					app: nginx
			spec:
				containers:
				- name: nginx
					image: nginx:1.14.2
					ports:
					- containerPort: 80
	---
	apiVersion: v1
	kind: Service
	metadata:
		name: my-service
	spec:
		selector:
			app.kubernetes.io/name: MyApp
		ports:
			- protocol: TCP
				port: 80
				targetPort: 9376
	`

	items, err := kube.SplitYAML([]byte(deploymentStr))
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		item.SetName(app.Name)
	}

	return items, nil
}

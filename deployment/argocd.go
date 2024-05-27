package deployment

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type ArgoCDDeploymentController struct {
	config    ArgoCDDeploymentConfig
	clientset *kubernetes.Clientset
	client    *dynamic.DynamicClient
}

type ArgoCDDeploymentConfig struct {
	Id                   string
	ArgocdNamespace      string
	AppName              string
	DestServer           string
	SourceRepoURL        string
	SourceTargetRevision string
	SourcePath           string
}

func NewArgoCDDeploymentController(config ArgoCDDeploymentConfig, client *dynamic.DynamicClient) *ArgoCDDeploymentController {
	return &ArgoCDDeploymentController{
		config: config,
		client: client,
	}
}

func (c ArgoCDDeploymentController) Deploy() {
	application := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "argoproj.io/v1alpha1",
			"kind":       "Application",
			"metadata": map[string]interface{}{
				"generateName": c.config.AppName,
				"namespace":    c.config.ArgocdNamespace,
			},
			"spec": map[string]interface{}{
				"project": "default",
				"destination": map[string]interface{}{
					"server":    c.config.DestServer,
					"namespace": "kbot-test",
				},
				"source": map[string]interface{}{
					"repoURL":        c.config.SourceRepoURL,
					"targetRevision": c.config.SourceTargetRevision,
					"path":           c.config.SourcePath,
				},
				"syncPolicy": map[string]interface{}{
					"automated": make(map[string]interface{}),
					"syncOptions": []string{
						"CreateNamespace=true",
					},
				},
			},
		},
	}

	application.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "argoproj.io",
		Version: "v1alpha1",
		Kind:    "Application",
	})

	list, err := c.client.Resource(schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}).Namespace(c.config.ArgocdNamespace).List(context.TODO(), metav1.ListOptions{})
	// list, err := c.client.Resource(schema.GroupVersionResource{
	// 	Group:    "apps",
	// 	Version:  "v1",
	// 	Resource: "deployments",
	// }).Namespace(c.config.ArgocdNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, dep := range list.Items {
		fmt.Print("\t")
		fmt.Println(dep.GetName())
	}

	// application := &unstructured.Unstructured{
	// 	Object: map[string]interface{}{
	// 		"apiVersion": "apps/v1",
	// 		"kind":       "Deployment",
	// 		"metadata": map[string]interface{}{
	// 			"name": "demo-deployment",
	// 		},
	// }

	appRes := schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "applications"}
	// appRes := schema.GroupVersionResource{Group: "argoproj.io", Version: "argoproj.io/v1alpha1", Resource: "applications.argoproj.io"}
	// client, err := dynamic.NewForConfig(config)
	// Create Deployment
	// ObjectMeta.GenerateName --- prefix
	fmt.Println("Creating deployment...")
	result, err := c.client.Resource(appRes).Namespace(c.config.ArgocdNamespace).Create(context.TODO(), application, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Print(result)

	jsonResult, _ := json.Marshal(result)
	fmt.Println(string(jsonResult))

	fmt.Println()
	fmt.Println("Resource created: ", result.GetName())
}

package deployment

import (
	"context"
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
	id                   string
	argocdNamespace      string
	appName              string
	destServer           string
	sourceRepoURL        string
	sourceTargetRevision string
	sourcePath           string
}

func NewArgoCDDeploymentController(config ArgoCDDeploymentConfig) *ArgoCDDeploymentController {
	return &ArgoCDDeploymentController{
		config: config,
	}
}

func (c ArgoCDDeploymentController) Deploy() {
	application := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "argoproj.io/v1alpha1",
			"kind":       "Application",
			"metadata": map[string]interface{}{
				"name":      c.config.appName,
				"namespace": c.config.argocdNamespace,
			},
			"spec": map[string]interface{}{
				"project": "default",
				"destination": map[string]interface{}{
					"server": c.config.destServer,
				},
				"source": map[string]interface{}{
					"repoURL":        c.config.sourceRepoURL,
					"targetRevision": c.config.sourceTargetRevision,
					"path":           c.config.sourcePath,
				},
				"syncPolicy": map[string]interface{}{
					"syncOptions": []string{
						"CreateNamespace=true",
					},
				},
			},
		},
	}

	// application := &unstructured.Unstructured{
	// 	Object: map[string]interface{}{
	// 		"apiVersion": "apps/v1",
	// 		"kind":       "Deployment",
	// 		"metadata": map[string]interface{}{
	// 			"name": "demo-deployment",
	// 		},
	// }

	appRes := schema.GroupVersionResource{Group: "argoproj.io", Version: "argoproj.io/v1alpha1", Resource: "applications"}
	// client, err := dynamic.NewForConfig(config)
	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := c.client.Resource(appRes).Namespace(c.config.argocdNamespace).Create(context.TODO(), application, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Print(result)
}

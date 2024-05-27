package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	apiv1 "k8s.io/api/core/v1"
	// v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func NewFluxDeploymentController(config FluxDeploymentConfig, clientset *kubernetes.Clientset, context context.Context) *FluxDeploymentController {
	return &FluxDeploymentController{
		config:    config,
		clientset: clientset,
		context:   context,
	}
}

type FluxDeploymentController struct {
	config    FluxDeploymentConfig
	clientset *kubernetes.Clientset
	context   context.Context
	// client    *kubeclient.Client
	// client *dynamic.DynamicClient
}

type FluxDeploymentConfig struct {
	Id        string
	Namespace string
	// AppName              string
	// DestServer           string
	// SourceRepoURL        string
	// SourceTargetRevision string
	// SourcePath           string
}

func (c FluxDeploymentController) Deploy() {
	clientset := c.clientset
	podName := "flux-init"

	req := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:    "flux-cli",
					Image:   "fluxcd/flux-cli:v2.3.0",
					Command: []string{"check"},
					Args: []string{
						"--server", "https://192.168.1.119:16443",
					},
				},
			},
		},
	}

	// deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	pod, err := clientset.CoreV1().Pods(c.config.Namespace).Create(c.context, req, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	c.waitForPod(pod)

	fmt.Println("Pod created: ", pod.GetName())

	jsonResult, _ := json.Marshal(pod)
	fmt.Println(string(jsonResult))
}

func (c FluxDeploymentController) waitForPod(pod *apiv1.Pod) error {
	// c.clientset.CoreV1().Pods(c.config.Namespace).Watch(c.context)

	w, err := c.clientset.CoreV1().Pods(c.config.Namespace).Watch(c.context, metav1.ListOptions{
		Watch:           true,
		ResourceVersion: pod.ResourceVersion,
		FieldSelector:   fields.Set{"metadata.name": pod.GetName()}.String(),
		LabelSelector:   labels.Everything().String(),
	})
	if err != nil {
		return err
	}

	status := pod.Status
	func() {
		for {
			select {
			case events, ok := <-w.ResultChan():
				if !ok {
					return
				}
				resp := events.Object.(*apiv1.Pod)
				fmt.Println("Pod status:", resp.Status.Phase)
				status = resp.Status
				if resp.Status.Phase != apiv1.PodPending {
					w.Stop()
				}
			case <-time.After(10 * time.Second):
				fmt.Println("timeout to wait for pod active")
				w.Stop()
			}
		}
	}()

	if status.Phase != apiv1.PodRunning {
		return fmt.Errorf("pod is unavailable: %v", status.Phase)
	}

	return nil
}

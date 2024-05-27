/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"path/filepath"
	"yevhenhrytsai/k8s-client-test/deployment"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

type ControllerType string

const (
	ArgocdControllerType ControllerType = "argocd"
	FluxControllerType   ControllerType = "flux"
)

const (
	controller = FluxControllerType
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	switch controller {
	case ArgocdControllerType:
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err)
		}
		client, err := dynamic.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		deployment := deployment.NewArgoCDDeploymentController(deployment.ArgoCDDeploymentConfig{
			Id:                   "abc",
			AppName:              "kbot-",
			ArgocdNamespace:      "argocd",
			DestServer:           "https://192.168.1.119:16443",
			SourceRepoURL:        "https://github.com/yevgen-grytsay/kbot",
			SourceTargetRevision: "argocd",
			SourcePath:           "helm",
		}, client)

		deployment.Deploy()
	case FluxControllerType:
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		deployment := deployment.NewFluxDeploymentController(
			deployment.FluxDeploymentConfig{
				Namespace: "default",
			},
			clientset,
			context.Background(),
		)

		deployment.Deploy()
	}

	/* // use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	} */

	/* // create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		namespace := "default"
		pod := "example-xxxxx"
		_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}

		// time.Sleep(10 * time.Second)

		secretsList, err := clientset.CoreV1().Secrets("argocd").List(context.TODO(), metav1.ListOptions{})
		// fmt.Println(secretsList)
		for _, sec := range secretsList.Items {
			fmt.Println(sec.ObjectMeta.Name)
		}

		os.Exit(0)
	} */
}

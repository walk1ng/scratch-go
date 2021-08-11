package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/user"
	"path"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

/*
reference: https://kubernetes.io/zh/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/
*/

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	kubeconfig := path.Join(user.HomeDir, ".kube", "config")
	c, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		c, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	clients, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	dc, err := dynamic.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	// deploy name
	name := "patch-demo"

	// get
	_, err = clients.AppsV1().Deployments("default").Get(context.Background(), name, v1.GetOptions{})
	if err != nil {
		fmt.Printf("failed to get deployment %q: %v\n", name, err)
		return
	}

	// yaml patch
	patch, err := ioutil.ReadFile("./application/patch-file-replicas.yaml")
	if err != nil {
		fmt.Printf("failed to read patch file: %v\n", err)
		return
	}

	fmt.Printf("%s\n", string(patch))

	// json patch
	// patch must be json format to submit to k8s api-server
	patch, err = yaml.ToJSON(patch)
	if err != nil {
		fmt.Printf("failed to convert patch yaml to patch json: %v\n", err)
		return
	}

	fmt.Printf("%s\n", string(patch))

	_, err = clients.AppsV1().Deployments("default").Patch(context.Background(),
		name,
		types.MergePatchType,
		patch,
		v1.PatchOptions{})

	if err != nil {
		fmt.Printf("failed to patch deployment %q: %v\n", name, err)
		return
	}

	// unstructured
	r := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}
	obj, err := dc.Resource(r).Namespace("default").Get(context.Background(), name, v1.GetOptions{})
	if err != nil {
		fmt.Printf("failed to get unstructured: %v\n", err)
		return
	}

	fmt.Println(obj.GetName())

	// patch
	patch, err = ioutil.ReadFile("./application/patch-file-tolerations.yaml")
	if err != nil {
		fmt.Printf("failed to read patch file: %v\n", err)
		return
	}

	patch, err = yaml.ToJSON(patch)
	if err != nil {
		fmt.Printf("failed to convert patch yaml to patch json: %v\n", err)
		return
	}

	_, err = dc.Resource(r).Namespace("default").Patch(context.Background(), name,
		types.StrategicMergePatchType,
		patch,
		v1.PatchOptions{})
	if err != nil {
		fmt.Printf("failed to patch unstructured: %v\n", err)
		return
	}
}

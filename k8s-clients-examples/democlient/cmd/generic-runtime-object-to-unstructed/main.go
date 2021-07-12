package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"democlient/pkg/utils"
)

func main() {
	// build the config
	c, err := clientcmd.BuildConfigFromFlags("", utils.UserConfig())
	utils.ExitErr("failed to load client config", err)

	// create client or panic
	client := kubernetes.NewForConfigOrDie(c)

	// create dynamic client or panic
	// dc := dynamic.NewForConfigOrDie(c)

	/*
		1. get deployment
		2. convert the deployment to runtime.object
		3. convert the runtime.object to unstructured.Unstructured
		4. convert the unstructured.Unstructured to the deployment
	*/
	dep, err := client.AppsV1().Deployments("testns").Get("myapp", v1.GetOptions{})
	utils.ExitErr("failed to get deployment with client", err)

	// pod to runtime.object
	depObj := dep.DeepCopyObject()

	// convert runtime.object to unstructured.Unstructured
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(depObj)
	utils.ExitErr("failed to convert runtime.object to unstructured.Unstructured", err)

	// convert unstructured.Unstructured to the deployment
	myDeployment := new(appsv1.Deployment)
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj, myDeployment)
	utils.ExitErr("failed to convert unstructured.Unstructured to deployment", err)

	fmt.Printf("%v\n", myDeployment)

}

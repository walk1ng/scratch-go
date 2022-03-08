package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	client, err := runtimeclient.New(config, runtimeclient.Options{
		// scheme for kubernetes resources
		Scheme: scheme.Scheme,
	})
	if err != nil {
		panic(err)
	}

	// list pods
	podList := &corev1.PodList{}
	err = client.List(context.Background(), podList, runtimeclient.InNamespace("default"))
	if err != nil {
		panic(err)
	}
	for _, pod := range podList.Items {
		println(pod.Name)
	}

	println("-------------------------")

	// get deployment
	deploy := appsv1.Deployment{}
	err = client.Get(context.Background(), types.NamespacedName{
		Namespace: "default",
		Name:      "loki-kube-state-metrics",
	}, &deploy)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", deploy.Status)

	println("-------------------------")

	// create pod
	pod := corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name:      "myapp",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "myapp",
					Image:           "nginx",
					ImagePullPolicy: corev1.PullIfNotPresent,
				},
			},
		},
	}
	err = client.Create(context.Background(), &pod, &runtimeclient.CreateOptions{})
	if err != nil {
		panic(err)
	}
}

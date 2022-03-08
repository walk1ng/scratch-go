package main

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	kudov1beta1 "github.com/kudobuilder/kudo/pkg/apis/kudo/v1beta1"
	redisv1 "github.com/spotahome/redis-operator/api/redisfailover/v1"
)

const (
	deploymentName   = "kube-ops-view"
	defaultNamespace = corev1.NamespaceDefault
	podName          = "myapp"
	containerName    = "myapp"
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
	err = client.List(context.Background(), podList, runtimeclient.InNamespace(defaultNamespace))
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
		Namespace: defaultNamespace,
		Name:      deploymentName,
	}, &deploy)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", deploy.Status)

	println("-------------------------")

	// create pod
	err = client.Get(context.Background(), types.NamespacedName{Namespace: defaultNamespace, Name: podName}, &corev1.Pod{})
	if err == nil {
		println("skip to creating pod")
	}
	if err != nil {
		if errors.IsNotFound(err) {
			println("pod is not existed, creating")
			pod := corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Namespace: defaultNamespace,
					Name:      podName,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            containerName,
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
		} else {
			fmt.Printf("failed to get pod %s/%s\n", podName, defaultNamespace)
		}
	}

	println("-------------------------")

	// custom resources
	crScheme := runtime.NewScheme()
	// scheme to redisfailover
	redisv1.AddToScheme(crScheme)
	// scheme to kudo
	kudov1beta1.AddToScheme(crScheme)

	cl, err := runtimeclient.New(config, runtimeclient.Options{
		Scheme: crScheme,
	})
	if err != nil {
		panic(err)
	}

	redisList := redisv1.RedisFailoverList{}
	err = cl.List(context.Background(), &redisList, runtimeclient.InNamespace(defaultNamespace))
	if err != nil {
		panic(err)
	}
	println(len(redisList.Items), "redises")
	for _, redis := range redisList.Items {
		println(redis.Name, "/", redis.Namespace)
	}

	instanceList := kudov1beta1.InstanceList{}
	err = cl.List(context.Background(), &instanceList)
	if err != nil {
		panic(err)
	}
	println(len(instanceList.Items), "kudo instances")
	for _, instance := range instanceList.Items {
		println(instance.Name, "/", instance.Namespace)
	}

	println("-------------------------")
	/*
		use controller-runtime client to create resource
		which CR type implement client.Object interface
	*/

}

package main

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

const (
	podName   = "kube-apiserver-xx-xx-xx"
	nodeName  = "xx-xx-xx"
	namespace = "kube-system"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	client, err := metricsv.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	nometrics, err := client.MetricsV1beta1().NodeMetricses().Get(context.Background(), nodeName, v1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println(nometrics.Name, ":")
	for k, v := range nometrics.Usage {
		fmt.Printf("%-10s:%-20s\n", strings.ToUpper(k.String()), &v)
	}

	poMetrics, err := client.MetricsV1beta1().PodMetricses(namespace).Get(context.Background(), podName, v1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(podName, ":")
	for _, container := range poMetrics.Containers {
		fmt.Println(container.Name, ":")
		for k, v := range container.Usage {
			if k == corev1.ResourceCPU {
				fmt.Printf("%-10s:%-20s(%dm)\n", strings.ToUpper((k.String())), &v, v.MilliValue())
			}
			if k == corev1.ResourceMemory {
				fmt.Printf("%-10s:%-20s(%dMi)\n", strings.ToUpper((k.String())), &v, v.Value()/1024/1024)
			}
		}
	}
}

package k8sunittestdemo

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type API struct {
	Client kubernetes.Interface
}

func (api API) NewPodWithMeta(namespace, name string) error {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: corev1.PodSpec{
			NodeName: "test_host",
		},
	}
	_, err := api.Client.CoreV1().Pods(namespace).Create(context.Background(), pod, metav1.CreateOptions{})
	return err
}

type Cache struct {
	Client kubernetes.Interface
	Pods   map[string]*corev1.Pod
}

func (c Cache) AddPodInCache(p *corev1.Pod) error {
	c.Pods[p.Name] = p
	return nil
}

package main

import (
	"context"
	"os/user"
	"path"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var logger *zap.Logger

var theten int64 = 10

func main() {
	initLogger()

	config, err := getK8sConfig()
	if err != nil {
		panic(err)
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Panic("new k8s client", zap.Error(err))
	}

	corev1Client := c.CoreV1()

	newPod := generatePod()
	logger.Info("create pod", zap.String("name", newPod.ObjectMeta.Name), zap.String("namespace", newPod.ObjectMeta.Namespace))

	pod, err := corev1Client.Pods("default").Create(context.Background(), newPod, v1.CreateOptions{})
	if err != nil {
		logger.Error("create pod", zap.String("name", pod.ObjectMeta.Name), zap.String("namespace", pod.ObjectMeta.Namespace), zap.Error(err))
		return
	}

	watcher, err := corev1Client.Pods("default").Watch(context.Background(), v1.SingleObject(pod.ObjectMeta))
	if err != nil {
		logger.Error("watch pod", zap.String("name", pod.ObjectMeta.Name), zap.String("namespace", pod.ObjectMeta.Namespace), zap.Error(err))
		return
	}

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Modified:
			pod := event.Object.(*corev1.Pod)
			for _, cond := range pod.Status.Conditions {
				if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
					logger.Info("stop watch pod", zap.String("name", pod.ObjectMeta.Name), zap.String("namespace", pod.ObjectMeta.Namespace), zap.String("condition", string(cond.Type)), zap.String("status", string(cond.Status)))
					watcher.Stop()

					logger.Info("delete pod", zap.String("name", pod.ObjectMeta.Name), zap.String("namespace", pod.ObjectMeta.Namespace))
					err := corev1Client.Pods("default").Delete(context.Background(), pod.Name, v1.DeleteOptions{GracePeriodSeconds: &theten})
					if err != nil {
						logger.Info("delete pod", zap.Error(err), zap.String("name", pod.ObjectMeta.Name), zap.String("namespace", pod.ObjectMeta.Namespace))
						return
					}
					break
				}
				logger.Info("still watch pod", zap.String("name", pod.ObjectMeta.Name), zap.String("namespace", pod.ObjectMeta.Namespace), zap.String("condition", string(cond.Type)), zap.String("status", string(cond.Status)))
			}
		default:
			logger.Panic("unexpected event", zap.String("type", string(event.Type)))
		}
	}
}

func getK8sConfig() (*rest.Config, error) {
	u, err := user.Current()
	if err != nil {
		logger.Error("get k8s config", zap.Error(err))
		return nil, err
	}
	configPath := path.Join(u.HomeDir, ".kube", "config")
	c, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		logger.Error("get k8s config", zap.Error(err))
		return nil, err
	}
	return c, nil
}

func generatePod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name:      "worker",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "worker",
					Image:           "busybox",
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command:         []string{"cat"},
					Stdin:           true,
				},
			},
			TerminationGracePeriodSeconds: &theten,
		},
	}
}

func initLogger() {
	logger, _ = zap.NewProduction()
}

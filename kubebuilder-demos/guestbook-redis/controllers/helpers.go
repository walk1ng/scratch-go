package controllers

import (
	"context"
	"fmt"
	webappv1 "guestbook-redis/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

type redisRole string

const (
	leader   redisRole = "leader"
	follower redisRole = "follower"
)

// helper functions for Redis

func (r *RedisReconciler) leaderDeployment(redis webappv1.Redis) (*appsv1.Deployment, error) {
	defOne := int32(1)
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", redis.Name, leader),
			Namespace: redis.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &defOne,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"redis": redis.Name,
					"role":  string(leader),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"redis": redis.Name,
						"role":  string(leader),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  string(leader),
							Image: "k8s.gcr.io/redis:e2e",
							Ports: []corev1.ContainerPort{
								{
									Name:          "redis",
									ContainerPort: 6379,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewMilliQuantity(100000, resource.BinarySI),
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}

	// set the ownerReference
	if err := ctrl.SetControllerReference(&redis, deploy, r.Scheme); err != nil {
		return nil, err
	}

	return deploy, nil
}

func (r *RedisReconciler) followerDeployment(redis webappv1.Redis) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", redis.Name, follower),
			Namespace: redis.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: redis.Spec.FollowerReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"redis": redis.Name,
					"role":  string(follower),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"redis": redis.Name,
						"role":  string(follower),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  string(follower),
							Image: "gcr.io/google_samples/gb-redisslave:v3",
							Ports: []corev1.ContainerPort{
								{
									Name:          "redis",
									ContainerPort: 6379,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
									corev1.ResourceMemory: *resource.NewMilliQuantity(100000, resource.BinarySI),
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}

	// set the ownerReference
	if err := ctrl.SetControllerReference(&redis, deploy, r.Scheme); err != nil {
		return nil, err
	}

	return deploy, nil
}

func (r *RedisReconciler) desiredService(redis webappv1.Redis, role redisRole) (*corev1.Service, error) {
	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      redis.Name + string(role),
			Namespace: redis.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{},
			Selector: map[string]string{
				"redis": redis.Name,
				"role":  string(role),
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	// set the ownerReference
	if err := ctrl.SetControllerReference(&redis, svc, r.Scheme); err != nil {
		return nil, err
	}

	return svc, nil
}

// helper functions for GuestBook

func (r *GuestBookReconciler) desiredDeployment(book webappv1.GuestBook, redis webappv1.Redis) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      book.Name,
			Namespace: book.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: book.Spec.Frontend.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"guestbook": book.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"guestbook": book.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "frontend",
							Image:           "gcr.io/google-samples/gb-frontend:v4",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Env: []corev1.EnvVar{
								{
									Name:  "GET_HOSTS_FROM",
									Value: "env",
								},
								{
									Name:  "REDIS_MASTER_SERVICE_HOST",
									Value: redis.Status.LeaderService,
								},
								{
									Name:  "REDIS_SLAVE_SERVICE_HOST",
									Value: redis.Status.FollowerService,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 80,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Resources: *book.Spec.Frontend.Resources.DeepCopy(),
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}

	// set the onwerReference
	if err := ctrl.SetControllerReference(&book, deploy, r.Scheme); err != nil {
		return nil, err
	}

	return deploy, nil
}

func (r *GuestBookReconciler) desiredService(book webappv1.GuestBook) (*corev1.Service, error) {
	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      book.Name,
			Namespace: book.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.FromString("http"),
				},
			},
			Selector: map[string]string{
				"guestbook": book.Name,
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	// set the ownerReference
	if err := ctrl.SetControllerReference(&book, svc, r.Scheme); err != nil {
		return nil, err
	}

	return svc, nil
}

func (r *GuestBookReconciler) booksUsingRedis(obj handler.MapObject) []ctrl.Request {
	listOptions := []client.ListOption{
		client.MatchingField(".spec.redisName", obj.Meta.GetName()),
		client.InNamespace(obj.Meta.GetNamespace()),
	}

	var bookList webappv1.GuestBookList
	if err := r.List(context.Background(), &bookList, listOptions...); err != nil {
		return nil
	}

	res := make([]ctrl.Request, len(bookList.Items))
	for i, book := range bookList.Items {
		res[i].Name = book.GetName()
		res[i].Namespace = book.GetNamespace()
	}
	return res
}

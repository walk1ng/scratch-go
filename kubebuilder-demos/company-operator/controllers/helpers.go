package controllers

import (
	webappv1 "company-operator/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	defOne int32 = 1
)

func (r *EmployeeReconciler) desiredDeployment(em webappv1.Employee) (appsv1.Deployment, error) {
	dep := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: metav1.SchemeGroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "worker-" + em.Name,
			Namespace: em.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&em, em.GroupVersionKind()),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &defOne,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"role":     string(em.Spec.Role),
					"employee": em.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{

					Containers: []corev1.Container{{
						Name:            "worker",
						Image:           "busybox",
						Command:         []string{"echo", "start", "sleep", "600", "&&", "echo", "done"},
						ImagePullPolicy: corev1.PullIfNotPresent,
						Resources:       em.Spec.Resources,
					}},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	return dep, nil
}

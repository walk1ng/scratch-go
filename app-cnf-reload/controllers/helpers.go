package controllers

import (
	automatev1 "app-cnf-reload/api/v1"
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ApplicationReconciler) desiredDeployment(app automatev1.Application, conf automatev1.Configuration) (appsv1.Deployment, error) {
	deploy := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			Labels: map[string]string{
				"managedBy":    app.Name,
				"configuredBy": conf.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&app, app.GroupVersionKind()),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: conf.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": app.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": app.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "core",
							Image: "nginx:1.9.11",
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("50m"),
									corev1.ResourceMemory: resource.MustParse("100Mi"),
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("50m"),
									corev1.ResourceMemory: resource.MustParse("100Mi"),
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}
	return deploy, nil
}

func (r *ApplicationReconciler) applicationUsingConfiguration(o client.Object) []reconcile.Request {
	r.Log.Info("watch object changed", "namespace", o.GetNamespace(), "name", o.GetName())

	listOptions := []client.ListOption{
		client.InNamespace(o.GetNamespace()),
		client.MatchingFields{".spec.configurationName": o.GetName()},
	}

	var appList automatev1.ApplicationList
	if err := r.List(context.Background(), &appList, listOptions...); err != nil {
		r.Log.Error(err, "failed to list applications effected by the change", "listoptions", listOptions)
		return nil
	}

	requests := make([]reconcile.Request, len(appList.Items))

	for i, app := range appList.Items {
		requests[i] = reconcile.Request{NamespacedName: types.NamespacedName{
			Namespace: app.Namespace,
			Name:      app.Name,
		}}
		r.Log.Info("enqueue request", "namespace", app.Namespace, "name", app.Name)
	}
	return requests
}

func (r *ConfigurationReconciler) configurationsUsedByApplication(o client.Object) []reconcile.Request {
	r.Log.Info("watch object changed", "namespace", o.GetNamespace(), "name", o.GetName())

	confKey := types.NamespacedName{
		Namespace: o.GetNamespace(),
		Name:      o.(*automatev1.Application).Spec.ConfigurationName,
	}

	var conf automatev1.Configuration
	if err := r.Get(context.Background(), confKey, &conf); err != nil {
		if client.IgnoreNotFound(err) != nil {
			r.Log.Error(err, "failed to get configuration effected by the change", "namespace", confKey.Namespace, "name", confKey.Name)
		}
		return nil
	}

	r.Log.Info("enqueue request", "namespace", confKey.Namespace, "name", confKey.Name)
	return []reconcile.Request{{NamespacedName: confKey}}

}

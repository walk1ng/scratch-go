package controllers

import (
	webappv1 "company-operator/api/v1"
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	defOne int32 = 1
)

func (r *EmployeeReconciler) desiredDeployment(em webappv1.Employee) (appsv1.Deployment, error) {
	dep := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "worker-" + em.Name,
			Namespace: em.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&em, em.GroupVersionKind()),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &em.Spec.DesiredWorkers,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"role":     string(em.Spec.Role),
					"employee": em.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"role":     string(em.Spec.Role),
						"employee": em.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:            "worker",
						Image:           "busybox",
						Command:         []string{"sleep", "600"},
						ImagePullPolicy: corev1.PullIfNotPresent,
						Resources:       em.Spec.Resources,
					}},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}

	return dep, nil
}

func (r *EmployeeReconciler) getEmployeeWorkState(ctx context.Context, em *webappv1.Employee) error {
	var depl appsv1.Deployment
	if err := r.Get(ctx, types.NamespacedName{Namespace: em.Namespace, Name: "worker-" + em.Name}, &depl); err != nil {
		return err
	}

	if depl.Status.AvailableReplicas == 0 {
		em.Status.WorkState = webappv1.NotStart
	} else if depl.Status.AvailableReplicas < depl.Status.Replicas {
		em.Status.WorkState = webappv1.Prepare
	} else {
		em.Status.WorkState = webappv1.Working
	}

	em.Status.ActualWorkers = depl.Status.ReadyReplicas

	return nil
}

func (r *CompanyReconciler) companyForEmployees(o client.Object) []reconcile.Request {
	empl := o.(*webappv1.Employee)

	var comp webappv1.Company
	err := r.Get(context.Background(), types.NamespacedName{Namespace: o.GetNamespace(), Name: empl.Spec.Company}, &comp)
	if err != nil {
		r.Log.Error(err, "fetch company failed")
		return nil
	}
	return []reconcile.Request{reconcile.Request{NamespacedName: types.NamespacedName{
		Namespace: comp.Namespace,
		Name:      comp.Name,
	}}}
}

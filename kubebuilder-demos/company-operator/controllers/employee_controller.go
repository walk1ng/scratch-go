/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	webappv1 "company-operator/api/v1"
)

// EmployeeReconciler reconciles a Employee object
type EmployeeReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=employees,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=employees/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=employees/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Employee object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *EmployeeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("employee", req.NamespacedName)
	log.Info("reconciling employee")

	var emp webappv1.Employee
	err := r.Get(ctx, req.NamespacedName, &emp)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploy, err := r.desiredDeployment(emp)
	if err != nil {
		return ctrl.Result{}, err
	}
	applyOpts := []client.PatchOption{
		client.ForceOwnership,
		client.FieldOwner("employee-controller"),
	}
	err = r.Patch(ctx, &deploy, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.getEmployeeWorkState(ctx, &emp)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Status().Update(ctx, &emp)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled employee")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *EmployeeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Employee{}).
		/*
			Employee can watch the desired deployment by two ways:
			*1 - own the deployment
				Owns(&appsv1.Deployment{}).
			*2 - watch the deployment and as owner will be enqueued while events occured on deployments
				Watches(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{OwnerType: &webappv1.Employee{}, IsController: true}).
			these two ways can watch the only deploymnents owner by employee rather than all deployments.
		*/
		// Owns(&appsv1.Deployment{}).
		Watches(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{OwnerType: &webappv1.Employee{}, IsController: true}).
		Complete(r)
}

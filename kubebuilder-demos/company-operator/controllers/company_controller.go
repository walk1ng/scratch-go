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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	webappv1 "company-operator/api/v1"
)

// CompanyReconciler reconciles a Company object
type CompanyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=companies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=companies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=companies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Company object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *CompanyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("company", req.NamespacedName)
	log.Info("reconciling company")

	// your logic here
	var comp webappv1.Company
	err := r.Get(ctx, req.NamespacedName, &comp)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var empls webappv1.EmployeeList
	listOpts := []client.ListOption{
		client.InNamespace(comp.Namespace),
		client.MatchingFields{".spec.company": comp.Name},
	}
	err = r.List(ctx, &empls, listOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	comp.Status.Employees = len(empls.Items)
	err = r.Status().Update(ctx, &comp)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled company")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CompanyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Company{}).
		Watches(&source.Kind{Type: &webappv1.Employee{}}, handler.EnqueueRequestsFromMapFunc(r.companyForEmployees)).
		Complete(r)
}

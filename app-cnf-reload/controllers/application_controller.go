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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	automatev1 "app-cnf-reload/api/v1"
)

// ApplicationReconciler reconciles a Application object
type ApplicationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=automate.walk1ng.dev,resources=applications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=automate.walk1ng.dev,resources=applications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=automate.walk1ng.dev,resources=applications/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=list;watch;get;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Application object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("application", req.NamespacedName)

	log.Info("reconciling Application")
	// your logic here
	var app automatev1.Application
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var conf automatev1.Configuration
	if err := r.Get(ctx, types.NamespacedName{
		Namespace: app.Namespace,
		Name:      app.Spec.ConfigurationName,
	}, &conf); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploy, err := r.desiredDeployment(app, conf)
	if err != nil {
		return ctrl.Result{}, err
	}
	applyOptions := []client.PatchOption{client.ForceOwnership, client.FieldOwner("application-controller")}
	err = r.Patch(ctx, &deploy, client.Apply, applyOptions...)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled Application")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetFieldIndexer().IndexField(context.Background(), &automatev1.Application{},
		".spec.configurationName", func(o client.Object) []string {
			configurationName := o.(*automatev1.Application).Spec.ConfigurationName
			if configurationName == "" {
				return nil
			}
			return []string{configurationName}
		})
	return ctrl.NewControllerManagedBy(mgr).
		For(&automatev1.Application{}).
		Owns(&appsv1.Deployment{}).
		Watches(
			&source.Kind{Type: &automatev1.Configuration{}},
			handler.EnqueueRequestsFromMapFunc(r.applicationUsingConfiguration)).
		Complete(r)
}

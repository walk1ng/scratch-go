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
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	automatev1 "app-cnf-reload/api/v1"
)

// ConfigurationReconciler reconciles a Configuration object
type ConfigurationReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=automate.walk1ng.dev,resources=configurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=automate.walk1ng.dev,resources=configurations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=automate.walk1ng.dev,resources=configurations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Configuration object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *ConfigurationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("configuration", req.NamespacedName)

	log.Info("reconciling Configuration")
	// your logic here
	var conf automatev1.Configuration
	if err := r.Get(ctx, req.NamespacedName, &conf); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// fetch applications which refer the configuration
	var appList automatev1.ApplicationList
	listOpts := []client.ListOption{client.InNamespace(conf.Namespace), client.MatchingFields{".spec.configurationName": conf.Name}}
	err := r.List(ctx, &appList, listOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	var apps []string

	if len(appList.Items) == 0 {
		apps = []string{}
	} else {
		for _, app := range appList.Items {
			apps = append(apps, app.Name)
		}
	}

	conf.Status.Applications = strings.Join(apps, ",")
	err = r.Status().Update(ctx, &conf)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled Configuration")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigurationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&automatev1.Configuration{}).
		Watches(&source.Kind{Type: &automatev1.Application{}},
			handler.EnqueueRequestsFromMapFunc(r.configurationsUsedByApplication)).
		Complete(r)
}

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
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webappv1 "kubebuilder-advanced-demo/api/v1"
)

// NicePodReconciler reconciles a NicePod object
type NicePodReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	Eventer record.EventRecorder
}

//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=nicepods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=nicepods/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=nicepods/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NicePod object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *NicePodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("nicepod", req.NamespacedName)

	log.Info("reconciling NicePod")
	// your logic here
	var pod webappv1.NicePod
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// record event
	// Type of this event (Normal, Warning)
	r.Eventer.Eventf(&pod, "Normal", "Reconciling", "Reconciling NicePod %s/%s", req.Namespace, req.Name)

	log.Info("reconciled NicePod")
	// record event
	// Type of this event (Normal, Warning)
	r.Eventer.Eventf(&pod, "Normal", "Reconciled", "Reconciled NicePod %s/%s", req.Namespace, req.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NicePodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.NicePod{}).
		Complete(r)
}

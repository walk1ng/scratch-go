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
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	webappv1 "kubebuilder-advanced-demo/api/v1"
)

// NiceServiceReconciler reconciles a NiceService object
type NiceServiceReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	Eventer record.EventRecorder
}

//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=niceservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=niceservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=niceservices/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NiceService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *NiceServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("niceservice", req.NamespacedName)

	log.Info("reconciling NiceService")

	// your logic here
	var service webappv1.NiceService
	if err := r.Get(ctx, req.NamespacedName, &service); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// record event
	// Type of this event (Normal, Warning)
	r.Eventer.Eventf(&service, "Normal", "Reconciling", "Reconciling NiceService %s/%s", req.Namespace, req.Name)

	var podList webappv1.NicePodList
	listOpts := []client.ListOption{
		client.MatchingLabelsSelector{Selector: labels.SelectorFromSet(service.Spec.Selector)},
		client.InNamespace(service.GetNamespace()),
	}
	err := r.List(ctx, &podList, listOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	// matching no one nicepods
	if len(podList.Items) == 0 {
		return ctrl.Result{}, nil
	}

	for i, pod := range podList.Items {
		service.Status.EndPoints[i] = pod.Name
	}

	err = r.Status().Update(ctx, &service)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled NiceService")
	// record event
	// Type of this event (Normal, Warning)
	r.Eventer.Eventf(&service, "Normal", "Reconciled", "Reconciled NiceService %s/%s", req.Namespace, req.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NiceServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.NiceService{}).
		Watches(&source.Kind{Type: &webappv1.NicePod{}}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}

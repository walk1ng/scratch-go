/*


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

	webappv1 "guestbook-redis/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// GuestBookReconciler reconciles a GuestBook object
type GuestBookReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=guestbooks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=guestbooks/status,verbs=get;update;patch

func (r *GuestBookReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("guestbook", req.NamespacedName)

	log.Info("reconciling guestbook")
	// your logic here
	var book webappv1.GuestBook
	if err := r.Get(ctx, req.NamespacedName, &book); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	var redis webappv1.Redis
	redisKey := client.ObjectKey{
		Namespace: book.Namespace,
		Name:      book.Spec.RedisName,
	}
	if err := r.Get(ctx, redisKey, &redis); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// deployment of guestbook
	deploy, err := r.desiredDeployment(book, redis)
	if err != nil {
		return ctrl.Result{}, err
	}
	// service of guestbook
	svc, err := r.desiredService(book)
	if err != nil {
		return ctrl.Result{}, err
	}

	// server-side apply
	// TODO: understand the patch pattern
	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner("guestbook-controller")}
	err = r.Patch(ctx, deploy, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Patch(ctx, svc, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	// update guestbook status
	book.Status.URL = "TODO"
	err = r.Status().Update(ctx, &book)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled guestbook")
	return ctrl.Result{}, nil
}

func (r *GuestBookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetFieldIndexer().IndexField(&webappv1.GuestBook{},
		".spec.RedisName", func(o runtime.Object) []string {
			redisName := o.(*webappv1.GuestBook).Spec.RedisName
			if redisName == "" {
				return nil
			}
			return []string{redisName}
		})

	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.GuestBook{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Watches(&source.Kind{
			Type: &webappv1.Redis{},
		}, &handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(r.booksUsingRedis),
		}).
		Complete(r)
}
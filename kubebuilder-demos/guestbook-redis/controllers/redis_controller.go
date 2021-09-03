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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webappv1 "guestbook-redis/api/v1"
)

// RedisReconciler reconciles a Redis object
type RedisReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=redis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.walk1ng.dev,resources=redis/status,verbs=get;update;patch

func (r *RedisReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("redis", req.NamespacedName)

	// your logic here
	log.Info("reconciling redis")

	var redis webappv1.Redis
	if err := r.Get(ctx, req.NamespacedName, &redis); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// leader deployment
	deplLeader, err := r.leaderDeployment(redis)
	if err != nil {
		return ctrl.Result{}, err
	}
	// leader service
	svcLeader, err := r.desiredService(redis, leader)
	if err != nil {
		return ctrl.Result{}, err
	}

	// follower deployment
	deplFollower, err := r.followerDeployment(redis)
	if err != nil {
		return ctrl.Result{}, err
	}
	// follower service
	svcFollower, err := r.desiredService(redis, follower)
	if err != nil {
		return ctrl.Result{}, err
	}

	// TODO: understand the patch pattern
	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner("redis-controller")}
	err = r.Patch(ctx, deplLeader, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Patch(ctx, svcLeader, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Patch(ctx, deplFollower, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Patch(ctx, svcFollower, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	// update the status
	redis.Status.LeaderService = svcLeader.Name
	redis.Status.FollowerService = svcFollower.Name
	if err := r.Status().Update(ctx, &redis); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled redis")
	return ctrl.Result{}, nil
}

func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Redis{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}

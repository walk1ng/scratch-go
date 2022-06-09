/*
Copyright 2022.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	demov1 "watch-external-res/api/v1"
	"watch-external-res/util"
)

const (
	configMapField                = ".spec.configMap"
	configDeploymentFinalizer     = "configdeployment.demo.walk1ng.dev/finalizer"
	referConfigMapResourceVersion = "refer-configmap-rev"
	volumeName                    = "config"
	mountPath                     = "/etc/app"
)

var (
	theOne int32 = 1
)

// ConfigDeploymentReconciler reconciles a ConfigDeployment object
type ConfigDeploymentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=demo.walk1ng.dev,resources=configdeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demo.walk1ng.dev,resources=configdeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demo.walk1ng.dev,resources=configdeployments/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *ConfigDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log = log.FromContext(ctx).WithValues("configDeployment", req.NamespacedName)

	// TODO(user): your logic here
	var cfd demov1.ConfigDeployment
	if err := r.Get(ctx, req.NamespacedName, &cfd); err != nil {
		if err := client.IgnoreNotFound(err); err == nil {
			r.Log.Info("The configDeployment resource was gone, ignore the request")
			return ctrl.Result{}, nil
		} else {
			r.Log.Error(err, "Unable to fetch configDeployment")
			return ctrl.Result{}, err
		}
	}

	if cfd.DeletionTimestamp.IsZero() {
		r.Log.Info("Enter apply logic for configDeployment CR")
		// apply phase
		if !util.SliceContainsString(cfd.Finalizers, configDeploymentFinalizer) {
			cfd.Finalizers = append(cfd.Finalizers, configDeploymentFinalizer)
			if err := r.Update(ctx, &cfd); err != nil {
				return ctrl.Result{}, err
			}
		}

		// apply the deployment here
		var configmapResVersion string
		if cfd.Spec.ConfigMap != "" {
			configMapName := cfd.Spec.ConfigMap
			foundConfigMap := &corev1.ConfigMap{}
			key := types.NamespacedName{
				Namespace: cfd.Namespace,
				Name:      configMapName,
			}
			if err := r.Get(ctx, key, foundConfigMap); err != nil {
				r.Log.Error(err, "Unable to fetch specified configmap", "namespacedname", key)
				return ctrl.Result{}, err
			}
			configmapResVersion = foundConfigMap.ResourceVersion
		}

		foundDeployment := &appsv1.Deployment{}
		key := types.NamespacedName{
			Namespace: cfd.Namespace,
			Name:      cfd.Name + "-deploy",
		}
		if err := r.Get(ctx, key, foundDeployment); err != nil {
			if client.IgnoreNotFound(err) == nil {
				// create deployment
				r.Log.Info("Set the resourceversion of configmap in pod template")

				deploy, err := r.desiredDeployment(&cfd)
				if err != nil {
					return ctrl.Result{}, err
				}

				deploy.Spec.Template.ObjectMeta.Annotations[referConfigMapResourceVersion] = configmapResVersion
				if err := r.Create(ctx, deploy); err != nil {
					r.Log.Error(err, "Failed to create the deployment")
					return ctrl.Result{}, err
				}
				r.Log.Info("Create the deployment successfully")
				return ctrl.Result{}, nil
			} else {
				r.Log.Error(err, "Unable to fetch the deployment")
				return ctrl.Result{}, err
			}
		}

		// found the deployment, update if required
		if foundDeployment.Spec.Template.ObjectMeta.Annotations[referConfigMapResourceVersion] != configmapResVersion {
			r.Log.Info("Refresh the resourceversion of configmap in pod template")
			foundDeployment.Spec.Template.ObjectMeta.Annotations[referConfigMapResourceVersion] = configmapResVersion

			r.Log.Info("Refresh the volumes configurations in pod template")
			r.handleVolumesForDeployment(&cfd, foundDeployment)
			if err := r.Update(ctx, foundDeployment); err != nil {
				r.Log.Error(err, "Failed to update the deployment")
				return ctrl.Result{}, err
			}
		}
	} else {
		// cfd is marked as deleted
		r.Log.Info("Enter delete logic for configDeployment CR")
		if util.SliceContainsString(cfd.Finalizers, configDeploymentFinalizer) {
			// clear deployment
			r.Log.Info("pre-remove the deployment before the finalizer gone")
			key := types.NamespacedName{
				Namespace: cfd.Namespace,
				Name:      cfd.Name + "-deploy",
			}
			foundDeployment := &appsv1.Deployment{}
			if err := r.Get(ctx, key, foundDeployment); err != nil {
				if client.IgnoreNotFound(err) != nil {
					return ctrl.Result{}, err
				}
				r.Log.Info("Deployment was gone")
			} else {
				if err := r.Delete(ctx, foundDeployment); err != nil {
					r.Log.Error(err, "Failed to delete deployment")
					return ctrl.Result{}, err
				}
				r.Log.Info("Delete deployment successfully")
			}
		}

		r.Log.Info("Remove configDeployment finalizer, the configDeployment will be GC by kubernetes soon")
		cfd.Finalizers = util.RemoveStringFromSlice(cfd.Finalizers, configDeploymentFinalizer)
		if err := r.Update(ctx, &cfd); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &demov1.ConfigDeployment{}, configMapField,
		func(o client.Object) []string {
			cfd := o.(*demov1.ConfigDeployment)
			if cfd.Spec.ConfigMap == "" {
				return nil
			}
			return []string{cfd.Spec.ConfigMap}
		}); err != nil {
		r.Log.Error(err, "failed to index field for configDeployment")
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.ConfigDeployment{}).
		Owns(&appsv1.Deployment{}).
		Watches(&source.Kind{
			Type: &corev1.ConfigMap{},
		}, handler.EnqueueRequestsFromMapFunc(r.findObjectForConfigMap),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}

func (r *ConfigDeploymentReconciler) findObjectForConfigMap(configMap client.Object) []reconcile.Request {
	affectedConfigDeploymentList := &demov1.ConfigDeploymentList{}

	err := r.List(context.Background(), affectedConfigDeploymentList, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(configMapField, configMap.GetName()),
		Namespace:     configMap.GetNamespace(),
	})
	if err != nil {
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, len(affectedConfigDeploymentList.Items))
	for i, cfd := range affectedConfigDeploymentList.Items {
		requests[i] = reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: cfd.Namespace,
				Name:      cfd.Name,
			},
		}
	}

	return requests
}

func (r *ConfigDeploymentReconciler) desiredDeployment(cfd *demov1.ConfigDeployment) (*appsv1.Deployment, error) {
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfd.Name + "-deploy",
			Namespace: cfd.Namespace,
			Labels:    map[string]string{},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &theOne,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"configDeployment": cfd.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"configDeployment": cfd.Name,
					},
					Annotations: make(map[string]string),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "app",
							Image:           "busybox",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command: []string{
								"cat",
								mountPath + "/*",
								"&&",
								"sleep",
								"3600",
							},
						},
					},
				},
			},
		},
	}

	// handle volumes
	r.handleVolumesForDeployment(cfd, deploy)

	// onwer reference
	if err := ctrl.SetControllerReference(cfd, deploy, r.Scheme); err != nil {
		return nil, err
	}

	return deploy, nil
}

func (r *ConfigDeploymentReconciler) handleVolumesForDeployment(cfd *demov1.ConfigDeployment, in *appsv1.Deployment) *appsv1.Deployment {
	if cfd.Spec.ConfigMap != "" {
		// volumes
		volumes := make([]corev1.Volume, 1)
		volumes[0] = corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: cfd.Spec.ConfigMap},
				},
			},
		}
		in.Spec.Template.Spec.Volumes = volumes
		r.Log.Info("Configure the volumes for the deployment")

		// volume mounts
		volumeMounts := make([]corev1.VolumeMount, 1)
		volumeMounts[0] = corev1.VolumeMount{
			Name:      volumeName,
			ReadOnly:  true,
			MountPath: mountPath,
		}
		in.Spec.Template.Spec.Containers[0].VolumeMounts = volumeMounts
		r.Log.Info("Configure the volumeMounts for the deployment")
	} else {
		if len(in.Spec.Template.Spec.Volumes) != 0 {
			in.Spec.Template.Spec.Volumes = nil
			r.Log.Info("Unset the volumes for the deployment")
		}

		if len(in.Spec.Template.Spec.Containers[0].VolumeMounts) != 0 {
			in.Spec.Template.Spec.Containers[0].VolumeMounts = nil
			r.Log.Info("Unset the volumeMounts for the deployment")
		}
	}

	// do nothing
	return in
}

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

package v1

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var configurationlog = logf.Log.WithName("configuration-resource")

func (r *Configuration) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-automate-walk1ng-dev-v1-configuration,mutating=true,failurePolicy=fail,sideEffects=None,groups=automate.walk1ng.dev,resources=configurations,verbs=create;update,versions=v1,name=mconfiguration.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &Configuration{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Configuration) Default() {
	configurationlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	if len(r.Labels) == 0 {
		r.Labels["configurations.automate.walk1ng.dev"] = "true"
		r.Labels["producedBy"] = "walk1ng.dev"
	}

}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-automate-walk1ng-dev-v1-configuration,mutating=false,failurePolicy=fail,sideEffects=None,groups=automate.walk1ng.dev,resources=configurations,verbs=create;update,versions=v1,name=vconfiguration.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Configuration{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Configuration) ValidateCreate() error {
	configurationlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Configuration) ValidateUpdate(old runtime.Object) error {
	configurationlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Configuration) ValidateDelete() error {
	configurationlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Configuration) validate() error {
	errReplicasLt10 := errors.Errorf("replicas of configurations %s/%s is larger than 10", r.Namespace, r.Name)

	if *r.Spec.Replicas > 10 {
		return errReplicasLt10
	}

	return nil
}

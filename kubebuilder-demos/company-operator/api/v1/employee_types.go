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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Role string

const (
	Engineer Role = "engineer"
	Sales    Role = "sales"
	Manager  Role = "manager"
)

type WorkState string

const (
	NotStart WorkState = "notStart"
	Prepare  WorkState = "prepare"
	Working  WorkState = "working"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EmployeeSpec defines the desired state of Employee
type EmployeeSpec struct {
	Role      Role                        `json:"role,omitempty"`
	Company   string                      `json:"company,omitempty"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// +kubebuilder:validation:Minimum=1
	DesiredWorkers int32 `json:"desiredWorkers,omitempty"`
}

// EmployeeStatus defines the observed state of Employee
type EmployeeStatus struct {
	ActualWorkers int32     `json:"actualWorkers,omitempty"`
	WorkState     WorkState `json:"workState,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".spec.role",name="Role",type="string"
// +kubebuilder:printcolumn:JSONPath=".spec.company",name="Company",type="string"
// +kubebuilder:printcolumn:JSONPath=".spec.desiredWorkers",name="DesiredWorkers",type="integer"
// +kubebuilder:printcolumn:JSONPath=".status.actualWorkers",name="ActualWorkers",type="integer"
// +kubebuilder:printcolumn:JSONPath=".status.workState",name="WorkState",type="string"

// Employee is the Schema for the employees API
type Employee struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmployeeSpec   `json:"spec,omitempty"`
	Status EmployeeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EmployeeList contains a list of Employee
type EmployeeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Employee `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Employee{}, &EmployeeList{})
}

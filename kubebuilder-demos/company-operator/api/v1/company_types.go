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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CompanySpec defines the desired state of Company
type CompanySpec struct {
	City    string `json:"city,omitempty"`
	Address string `json:"address,omitempty"`
}

// CompanyStatus defines the observed state of Company
type CompanyStatus struct {
	Employees int `json:"employees,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".spec.city",name="City",type="string"
// +kubebuilder:printcolumn:JSONPath=".spec.address",name="Addr",type="string"
// +kubebuilder:printcolumn:JSONPath=".status.employees",name="Employees",type="integer"

// Company is the Schema for the companies API
type Company struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CompanySpec   `json:"spec,omitempty"`
	Status CompanyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CompanyList contains a list of Company
type CompanyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Company `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Company{}, &CompanyList{})
}

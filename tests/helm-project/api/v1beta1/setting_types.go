/*
Copyright 2023.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SettingSpec defines the desired state of Setting
type SettingSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Setting. Edit setting_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

type SettingPhase string

// These are the valid phases of node.
const (
	SettingPending SettingPhase = "Pending"
	SettingUnknown SettingPhase = "Unknown"
	SettingActive  SettingPhase = "Active"
)

// SettingStatus defines the observed state of Setting
type SettingStatus struct {
	// Phase represents the current phase of Setting.
	//+kubebuilder:default:=Unknown
	Phase SettingPhase `json:"phase,omitempty"`
	// Represents the observations of a Setting's current state.
	// Setting.status.conditions.type are: "Available", "Progressing", and "Degraded"
	// Setting.status.conditions.status are one of True, False, Unknown.
	// Setting.status.conditions.reason the value should be a CamelCase string and producers of specific
	// condition types may define expected values and meanings for this field, and whether the values
	// are considered a guaranteed API.
	// Setting.status.conditions.Message is a human readable message indicating details about the transition.
	// For further information see: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Setting is the Schema for the settings API
type Setting struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SettingSpec   `json:"spec,omitempty"`
	Status SettingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SettingList contains a list of Setting
type SettingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Setting `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Setting{}, &SettingList{})
}

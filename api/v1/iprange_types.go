/*
Copyright 2020.

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

// IPRangeSpec defines the desired state of IPRange
type IPRangeSpec struct {
	// Range represent the IP range in CIDR format
	// i.e. 10.0.0.0/16 or 2001:db2::/64
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:MinLength=8
	Range string `json:"range,omitempty"`

	// +optional
	// Addresses represent the IP addresses of the range and its status.
	// Each address may be associated to one kubernetes object (i.e. Services)
	// +listType=set
	Addresses []string `json:"addresses,omitempty"`
}

// IPRangeStatus defines the observed state of IPRange
type IPRangeStatus struct {
	// Free represent the number of IP addresses that are not allocated in the Range
	// +optional
	Free int64 `json:"free,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// IPRange is the Schema for the ipranges API
type IPRange struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IPRangeSpec   `json:"spec,omitempty"`
	Status IPRangeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IPRangeList contains a list of IPRange
type IPRangeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPRange `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IPRange{}, &IPRangeList{})
}

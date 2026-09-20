/*
Copyright 2026.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AutolockSpec defines the desired state of Autolock.
type AutolockSpec struct {
    Webapp WebappSpec `json:"webapp,omitempty"`
    Main   MainSpec   `json:"main,omitempty"`
    MQTT   MQTTSpec   `json:"mqtt,omitempty"`
}

type MainSpec struct {
    Config MainConfig `json:"config,omitempty"`
	Replicas int32  `json:"replicas,omitempty"`
}

type MainConfig struct {
    AutoLock bool               `json:"auto_lock"`
    TimeoutSeq  int             `json:"timeout_seq"`
	IgnoreClsw bool             `json:"ignore_clsw"`
	RotateDirection string      `json:"rotate_direction"`
	AuthorizeInternalUsers bool `json:"authorize_internal_users"`
	AuthorizeExternalUsers bool `json:"authorize_external_users"`
}

type WebappSpec struct{
	Replicas int32  `json:"replicas,omitempty"`
}

type MQTTSpec struct{
	Replicas int32  `json:"replicas,omitempty"`
}

// AutolockStatus defines the observed state of Autolock.
type AutolockStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Autolock is the Schema for the autolocks API.
type Autolock struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutolockSpec   `json:"spec,omitempty"`
	Status AutolockStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AutolockList contains a list of Autolock.
type AutolockList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Autolock `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Autolock{}, &AutolockList{})
}

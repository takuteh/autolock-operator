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
	Image    string `json:"image"`
	Config   Config `json:"config,omitempty"`
	Replicas int32  `json:"replicas,omitempty"`
}

type Config struct {
	MainConfig MainConfig `json:"main,omitempty"`
	MQTTConfig MQTTConfig `json:"mqtt,omitempty"`
}

type MainConfig struct {
	AutoLock               bool   `json:"autolock"`
	TimeoutSeq             int    `json:"timeout_seq"`
	IgnoreClsw             bool   `json:"ignore_clsw"`
	RotateDirection        string `json:"rotate_direction"`
	AuthorizeInternalUsers bool   `json:"authorize_internal_users"`
	AuthorizeExternalUsers bool   `json:"authorize_external_users"`
}

type MQTTConfig struct {
	BrokerAddress string        `json:"broker_address"`
	MQTTPort      int           `json:"mqtt_port"`
	Publish       MQTTPublish   `json:"publish"`
	Subscribe     MQTTSubscribe `json:"subscribe"`
}

type MQTTPublish struct {
	Boot MQTTMessage `json:"boot"`
}

type MQTTSubscribe struct {
	ChangeConfig MQTTTopicMessage `json:"change_config"`
	Close        MQTTMessage      `json:"close"`
	Open         MQTTMessage      `json:"open"`
	RelayOff     MQTTMessage      `json:"relay_off"`
	RelayOn      MQTTMessage      `json:"relay_on"`
}

type MQTTMessage struct {
	Message string `json:"message"`
	Topic   string `json:"topic"`
}

type MQTTTopicMessage struct {
	Topic string `json:"topic"`
}

type WebappSpec struct {
	Image    string      `json:"image"`
	Replicas int32       `json:"replicas,omitempty"`
	Ingress  IngressSpec `json:"ingress,omitempty"`
}

type IngressSpec struct {
	Enabled   bool   `json:"enabled"`
	ClassName string `json:"className,omitempty"`
	Host      string `json:"host"`
	Path      string `json:"path,omitempty"`
	PathType  string `json:"pathType,omitempty"`
}

type MQTTSpec struct {
	Image    string `json:"image"`
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

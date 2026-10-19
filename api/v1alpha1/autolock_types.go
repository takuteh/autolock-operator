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

// AutolockSpec defines the desired configuration of an Autolock instance.
type AutolockSpec struct {
	// Webapp configures the web application used to manage Autolock settings and authorized users.
	Webapp WebappSpec `json:"webapp,omitempty"`

	// Main configures the main application that controls the Autolock mechanism.
	Main MainSpec `json:"main,omitempty"`

	// MQTT configures the MQTT broker used for communication between Autolock components.
	MQTT MQTTSpec `json:"mqtt,omitempty"`

	// AuthDB configures the database that stores user authentication information and expiration data.
	AuthDB AuthDBSpec `json:"authdb,omitempty"`

	// Webhook configures the application that receives requests from external services such as LINE and Slack
	// and translates them into API requests for Autolock.
	Webhook WebhookSpec `json:"webhook,omitempty"`
}

// MainSpec defines the configuration for the Autolock main application.
type MainSpec struct {
	// Image is the container image used by the main application.
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// Config defines the operational settings and MQTT messaging configuration for the main application.
	Config Config `json:"config"`

	// Replicas is the desired number of main application replicas.
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas,omitempty"`

	// PigpioHost is the hostname or IP address of the pigpio daemon used to control GPIO.
	// +kubebuilder:validation:MinLength=1
	PigpioHost string `json:"pigpioHost"`
}

// Config defines the configuration used by the main application.
type Config struct {
	// MainConfig defines the operational settings of the Autolock mechanism.
	MainConfig MainConfig `json:"main"`

	// MQTTConfig defines the MQTT broker connection and messaging settings.
	MQTTConfig MQTTConfig `json:"mqtt"`
}

// MainConfig defines the operational settings of the Autolock mechanism.
type MainConfig struct {
	// AutoLock determines whether the reed switch is used to detect door opening and closing
	// and enable automatic locking.
	AutoLock bool `json:"autolock"`

	// TimeoutSeq specifies the time in seconds to wait before locking the door
	// again when the lock has been released but the door has not been opened or closed.
	// Set to "no" to disable the timeout.
	TimeoutSeq string `json:"timeout_seq"`

	// IgnoreClsw determines whether close-button input is ignored while the door is open.
	IgnoreClsw bool `json:"ignore_clsw"`

	// RotateDirection specifies the thumb-turn rotation direction.
	// +kubebuilder:validation:Enum=right;left
	RotateDirection string `json:"rotate_direction"`

	// AuthorizeInternalUsers determines whether user authentication is required
	// for direct API requests received through MQTT from the same Layer 2 network.
	AuthorizeInternalUsers bool `json:"authorize_internal_users"`

	// AuthorizeExternalUsers determines whether user authentication is required
	// for requests received through external webhooks, such as LINE or Slack.
	AuthorizeExternalUsers bool `json:"authorize_external_users"`
}

// MQTTConfig defines the MQTT broker connection and messaging settings.
type MQTTConfig struct {
	// BrokerAddress is the hostname or IP address of the MQTT broker.
	// +kubebuilder:validation:MinLength=1
	BrokerAddress string `json:"broker_address"`

	// MQTTPort is the TCP port used to connect to the MQTT broker.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	MQTTPort int `json:"mqtt_port"`

	// Publish defines the MQTT messages published by the main application.
	Publish MQTTPublish `json:"publish"`

	// Subscribe defines the MQTT topics and messages consumed by the main application.
	Subscribe MQTTSubscribe `json:"subscribe"`
}

// MQTTPublish defines the MQTT publish message configuration.
type MQTTPublish struct {
	// Boot defines the MQTT message published when the Autolock main application starts.
	Boot MQTTMessage `json:"boot"`
}

// MQTTSubscribe defines the MQTT subscription configuration for Autolock operations.
type MQTTSubscribe struct {
	// ChangeConfig defines the topic used to receive configuration change notifications.
	ChangeConfig MQTTTopicMessage `json:"change_config"`

	// Close defines the MQTT topic and message associated with the close operation.
	Close MQTTMessage `json:"close"`

	// Open defines the MQTT topic and message associated with the open operation.
	Open MQTTMessage `json:"open"`

	// RelayOff defines the MQTT topic and message associated with turning the relay off.
	RelayOff MQTTMessage `json:"relay_off"`

	// RelayOn defines the MQTT topic and message associated with turning the relay on.
	RelayOn MQTTMessage `json:"relay_on"`
}

// MQTTMessage defines an MQTT topic and its message payload.
type MQTTMessage struct {
	// Message is the payload published to the MQTT topic.
	Message string `json:"message"`

	// Topic is the MQTT topic used to publish or receive the message.
	Topic string `json:"topic"`
}

// MQTTTopicMessage defines an MQTT topic used to receive notifications without a configured message payload.
type MQTTTopicMessage struct {
	// Topic is the MQTT topic used to receive notifications.
	Topic string `json:"topic"`
}

// WebappSpec defines the configuration of the Autolock web application.
type WebappSpec struct {
	// Image is the container image used by the web application.
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// Replicas is the desired number of web application replicas.
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas,omitempty"`

	// Ingress defines the external access configuration for the web application.
	Ingress IngressSpec `json:"ingress,omitempty"`
}

// WebhookSpec defines the configuration of the Autolock webhook application.
type WebhookSpec struct {
	// Image is the container image used by the webhook application.
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// Replicas is the desired number of webhook application replicas.
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas,omitempty"`

	// Ingress defines the external access configuration for the webhook application.
	Ingress IngressSpec `json:"ingress,omitempty"`
}

// IngressSpec defines the Ingress configuration for exposing an application.
type IngressSpec struct {
	// Enabled determines whether an Ingress resource is enabled for the application.
	Enabled bool `json:"enabled"`

	// ClassName specifies the IngressClass used to handle the Ingress.
	ClassName string `json:"className,omitempty"`

	// Host is the hostname used to route requests to the application.
	// +kubebuilder:validation:MinLength=1
	Host string `json:"host"`

	// Path is the URL path used to route requests to the application.
	Path string `json:"path,omitempty"`

	// PathType specifies how the Ingress path is matched.
	// +kubebuilder:validation:Enum=Exact;Prefix;ImplementationSpecific
	PathType string `json:"pathType,omitempty"`
}

// MQTTSpec defines the configuration of the MQTT broker workload.
type MQTTSpec struct {
	// Image is the container image used by the MQTT broker.
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// Replicas is the desired number of MQTT broker replicas.
	// +kubebuilder:validation:Minimum=0
	Replicas int32 `json:"replicas,omitempty"`
}

// AuthDBSpec defines the configuration of the authentication database.
type AuthDBSpec struct {
	// Image is the container image used by the authentication database.
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// SecretName is the name of the Kubernetes Secret containing database credentials.
	// +kubebuilder:validation:MinLength=1
	SecretName string `json:"secretName"`
}

// AutolockStatus defines the observed state of an Autolock instance.
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

	// Spec defines the desired state of the Autolock instance.
	Spec AutolockSpec `json:"spec,omitempty"`

	// Status reports the observed state of the Autolock instance.
	Status AutolockStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AutolockList contains a list of Autolock resources.
type AutolockList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items contains the Autolock resources in this list.
	Items []Autolock `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Autolock{}, &AutolockList{})
}
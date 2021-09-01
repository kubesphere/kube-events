/*
Copyright 2020 The KubeSphere Authors.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RulerSpec defines the desired state of Ruler
type RulerSpec struct {
	// Number of desired pods. Defaults to 1.
	Replicas *int32 `json:"replicas,omitempty"`
	// Docker image of events ruler
	Image string `json:"image"`
	// Image pull policy. One of Always, Never, IfNotPresent.
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// If specified, the pod's scheduling constraints.
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's tolerations.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// Resources defines resources requests and limits for single Pod.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Namespaces to be selected for Rules discovery.
	// If unspecified, discover Rule instances from all namespaces.
	RuleNamespaceSelector *metav1.LabelSelector `json:"ruleNamespaceSelector,omitempty"`
	// A selector to select Rules instances.
	RuleSelector *metav1.LabelSelector `json:"ruleSelector,omitempty"`
	// Sinks defines sinks detail of this ruler
	Sinks *RulerSinks `json:"sinks,omitempty"`
}

// RulerStatus defines the observed state of Ruler
type RulerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=erl

// Ruler is the Schema for the ruler API
type Ruler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specification of the desired behavior of the Ruler.
	Spec RulerSpec `json:"spec"`

	Status RulerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RulerList contains a list of Ruler
type RulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of Rulers
	Items []Ruler `json:"items"`
}

// RulerSinks defines a set of sinks for Events Ruler
type RulerSinks struct {
	// Alertmanager is an alertmanager sink to which only alerts can sink.
	Alertmanager *RulerAlertmanagerSink `json:"alertmanager,omitempty"`
	// Webhooks is a list of RulerWebhookSink to which notifications or alerts can sink
	Webhooks []*RulerWebhookSink `json:"webhooks,omitempty"`
	// Stdout can config write notifications or alerts to stdout; do nothing when no configuration
	Stdout *RulerStdoutSink `json:"stdout,omitempty"`
}

// RulerAlertmanagerSink is a sink to alertmanager service on k8s
type RulerAlertmanagerSink struct {
	// `namespace` is the namespace of the alertmanager service.
	Namespace string `json:"namespace"`
	// `name` is the name of the alertmanager service.
	Name string `json:"name"`
	// `port` is the port on the alertmanager service. Default to 9093.
	// `port` should be a valid port number (1-65535, inclusive).
	Port *int `json:"port,omitempty"`
	// TargetPort is the port to access on the backend instances targeted by the alertmanager service.
	// If this is not specified, the value of the 'port' field is used.
	TargetPort *int `json:"targetPort,omitempty"`
}

// RulerWebhookSink defines parameters for webhook sink of Events Ruler.
type RulerWebhookSink struct {
	// Type represents that the sink is for notification or alert.
	// Available values are `notification` and `alert`
	Type RulerSinkType `json:"type"`
	// `url` gives the location of the webhook, in standard URL form (`scheme://host:port/path`).
	// Exactly one of `url` or `service` must be specified.
	Url string `json:"url,omitempty"`
	// `service` is a reference to the service for this webhook. Either
	// `service` or `url` must be specified.
	// If the webhook is running within the cluster, then you should use `service`.
	Service *ServiceReference `json:"service,omitempty"`
}

// RulerStdoutSink defines parameters for stdout sink of Events Ruler.
type RulerStdoutSink struct {
	// Type represents that the sink is for notification or alert.
	// Available values are `notification` and `alert`
	Type RulerSinkType `json:"type"`
}

// ServiceReference holds a reference to k8s Service
type ServiceReference struct {
	// `namespace` is the namespace of the service.
	Namespace string `json:"namespace"`
	// `name` is the name of the service.
	Name string `json:"name"`
	// `port` is the port on the service and should be a valid port number (1-65535, inclusive).
	Port *int `json:"port,omitempty"`
	// `path` is an optional URL path which will be sent in any request to this service.
	Path string `json:"path,omitempty"`
}

type RulerSinkType string

const (
	// RulerSinkTypeNotification represents event notifications sink.
	RulerSinkTypeNotification = "notification"
	// RulerSinkTypeAlert represents alert messages sink.
	RulerSinkTypeAlert = "alert"
)

func init() {
	SchemeBuilder.Register(&Ruler{}, &RulerList{})
}

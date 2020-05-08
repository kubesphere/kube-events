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

// KubeEventsRulerSpec defines the desired state of KubeEventsRuler
type KubeEventsRulerSpec struct {
	Replicas        *int32            `json:"replicas,omitempty"`
	Image           string            `json:"image,omitempty"`
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Resources defines resources requests and limits for single Pod.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Namespaces to be selected for KubeEventsRules discovery.
	// If unspecified, discover KubeEventsRule instances from all namespaces.
	RuleNamespaceSelector *metav1.LabelSelector `json:"ruleNamespaceSelector,omitempty"`
	// A selector to select KubeEventsRules instances.
	RuleSelector *metav1.LabelSelector `json:"ruleSelector,omitempty"`
	// Sinks defines sinks detail of this ruler
	Sinks *RulerSinks `json:"sinks,omitempty"`
}

// KubeEventsRulerStatus defines the observed state of KubeEventsRuler
type KubeEventsRulerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// KubeEventsRuler is the Schema for the kubeeventsrulers API
type KubeEventsRuler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeEventsRulerSpec   `json:"spec,omitempty"`
	Status KubeEventsRulerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubeEventsRulerList contains a list of KubeEventsRuler
type KubeEventsRulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeEventsRuler `json:"items"`
}

type RulerSinks struct {
	// Alertmanager is for sinking alerts
	Alertmanager *RulerAlertmanagerSink `json:"alertmanager,omitempty"`
	// Webhooks is a list of RulerWebhookSink to which notifications or alerts can sink
	Webhooks []*RulerWebhookSink `json:"webhooks,omitempty"`
	Stdout   *RulerStdoutSink    `json:"stdout,omitempty"`
}

// RulerAlertmanagerSink is a sink to alertmanager service on k8s
type RulerAlertmanagerSink struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	Port      *int   `json:"port,omitempty"`
	// TargetPort is the port to access on the backend instances targeted by the service.
	// If this is not specified, the value of the 'port' field is used.
	TargetPort *int `json:"targetPort,omitempty"`
}

// RulerWebhookSink defines parameters for webhook sink of Events Ruler.
type RulerWebhookSink struct {
	// Type represents that the sink is for notification or alert.
	// Available values are `notification` and `alert`
	Type    RulerSinkType     `json:"type,omitempty"`
	Url     string            `json:"namespace,omitempty"`
	Service *ServiceReference `json:"service,omitempty"`
}

// RulerStdoutSink defines parameters for stdout sink of Events Ruler.
type RulerStdoutSink struct {
	Type RulerSinkType `json:"type,omitempty"`
}

// ServiceReference holds a reference to k8s Service
type ServiceReference struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
	Port      *int   `json:"port,omitempty"`
	Path      string `json:"path,omitempty"`
}

type RulerSinkType string

const (
	// RulerSinkTypeNotification represents event notifications sink.
	RulerSinkTypeNotification = "notification"
	// RulerSinkTypeAlert represents alert messages sink.
	RulerSinkTypeAlert = "alert"
)

func init() {
	SchemeBuilder.Register(&KubeEventsRuler{}, &KubeEventsRulerList{})
}

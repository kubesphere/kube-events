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

// ExporterSpec defines the desired state of Exporter
type ExporterSpec struct {
	// Docker image of events exporter
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
	// Sinks defines details of events sinks
	Sinks *ExporterSinks `json:"sinks,omitempty"`
}

// ExporterStatus defines the observed state of Exporter
type ExporterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=ex

// Exporter is the Schema for the exporter API
type Exporter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the specification of the desired behavior of the Exporter.
	Spec ExporterSpec `json:"spec"`

	Status ExporterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ExporterList contains a list of Exporter
type ExporterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// List of Exporters
	Items []Exporter `json:"items"`
}

// ExporterSinks defines a set of sinks for Events Exporter
type ExporterSinks struct {
	// Webhooks is a list of ExporterWebhookSink
	Webhooks []*ExporterWebhookSink `json:"webhooks,omitempty"`
	// Stdout represents whether to write events to stdout.
	// Output when configure an empty struct `{}`, but do nothing when no configuration
	Stdout *ExporterStdoutSink `json:"stdout,omitempty"`
}

// ExporterStdoutSink defines parameters for stdout sink of Events Exporter.
type ExporterStdoutSink struct {
}

// ExporterWebhookSink defines parameters for webhook sink of Events Exporter.
type ExporterWebhookSink struct {
	// `url` gives the location of the webhook, in standard URL form (`scheme://host:port/path`).
	// Exactly one of `url` or `service` must be specified.
	Url string `json:"url,omitempty"`
	// `service` is a reference to the service for this webhook. Either
	// `service` or `url` must be specified.
	// If the webhook is running within the cluster, then you should use `service`.
	Service *ServiceReference `json:"service,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Exporter{}, &ExporterList{})
}

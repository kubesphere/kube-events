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

// KubeEventsExporterSpec defines the desired state of KubeEventsExporter
type KubeEventsExporterSpec struct {
	Image string `json:"image,omitempty"`
	// Resources defines resources requests and limits for single Pod.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Sinks defines details of events sinks
	Sinks *ExporterSinks `json:"sinks,omitempty"`
}

// KubeEventsExporterStatus defines the observed state of KubeEventsExporter
type KubeEventsExporterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// KubeEventsExporter is the Schema for the kubeeventsexporters API
type KubeEventsExporter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeEventsExporterSpec   `json:"spec,omitempty"`
	Status KubeEventsExporterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubeEventsExporterList contains a list of KubeEventsExporter
type KubeEventsExporterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeEventsExporter `json:"items"`
}

// ExporterSinks defines a set of sinks for Events Exporter
type ExporterSinks struct {
	Webhooks []*ExporterWebhookSink `json:"webhooks,omitempty"`
	Stdout  *ExporterStdoutSink  `json:"stdout,omitempty"`
}

// ExporterStdoutSink defines parameters for stdout sink of Events Exporter.
type ExporterStdoutSink struct {
}

// ExporterWebhookSink defines parameters for webhook sink of Events Exporter.
type ExporterWebhookSink struct {
	Url string `json:"url,omitempty"`
	Service *ServiceReference `json:"service,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KubeEventsExporter{}, &KubeEventsExporterList{})
}

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RuleSpec defines the desired state of Rule
type RuleSpec struct {
	Rules []EventRule `json:"rules,omitempty"`
}

// RuleStatus defines the observed state of Rule
type RuleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=er

// Rule is the Schema for the Rule API
type Rule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RuleSpec   `json:"spec"`
	Status RuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RuleList contains a list of Rule
type RuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rule `json:"items"`
}

// EventRule describes a notification or alert rule
type EventRule struct {
	// Name is simple name of rule
	Name string `json:"name,omitempty"`
	// Condition is a string similar with the where part of sql.
	// See supported grammar details on https://github.com/kubesphere/event-rule-engine#supported-grammer .
	// For example: `type="Warning" and involvedObject.kind="Pod" and reason="FailedMount"`
	Condition string `json:"condition,omitempty"`
	// Labels
	Labels map[string]string `json:"labels,omitempty"`
	// Values of Annotations can use format string with the fields of the event.
	// For example: `{"message": "%message"}`
	Annotations map[string]string `json:"annotations,omitempty"`
	// Enable is whether to enable the rule, default to `false`
	Enable bool `json:"enable,omitempty"`
	// Type represents that the rule is for notification or alert.
	// Available values are `notification` and `alert`
	Type RuleType `json:"type,omitempty"`
}

type RuleType string

const (
	// RuleTypeNotification represents that the rule will used to generate notifications
	// based on the original event objects.
	RuleTypeNotification = "notification"
	// RuleTypeAlert represents that the rule will be used to generate alert messages
	// that conform to the alertmanager protocol.
	RuleTypeAlert = "alert"
)

func init() {
	SchemeBuilder.Register(&Rule{}, &RuleList{})
}

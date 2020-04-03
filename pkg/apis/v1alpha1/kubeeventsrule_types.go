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

// KubeEventsRuleSpec defines the desired state of KubeEventsRule
type KubeEventsRuleSpec struct {
	Rules []Rule `json:"rules,omitempty"`
}

// KubeEventsRuleStatus defines the observed state of KubeEventsRule
type KubeEventsRuleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// KubeEventsRule is the Schema for the kubeeventsrules API
type KubeEventsRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeEventsRuleSpec   `json:"spec,omitempty"`
	Status KubeEventsRuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KubeEventsRuleList contains a list of KubeEventsRule
type KubeEventsRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeEventsRule `json:"items"`
}

// Rule describes a notification or alerting rule
type Rule struct {
	Name      string `json:"name,omitempty"`
	Summary   string `json:"summary,omitempty"`
	SummaryCn string `json:"summaryCn,omitempty"`
	// Condition is a string similar with the where part of sql.
	// The usage is as follows: event.type="Warning" and event.involvedObject.kind="Pod" and event.reason="FailedMount"
	Condition string   `json:"condition,omitempty"`
	Message   string   `json:"message,omitempty"`
	Priority  string   `json:"priority,omitempty"`
	Source    string   `json:"source,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	Enable    bool     `json:"enable,omitempty"`
	Type      RuleType `json:"type,omitempty"`
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
	SchemeBuilder.Register(&KubeEventsRule{}, &KubeEventsRuleList{})
}

<br>
# API Docs

This Document documents the types introduced by the Kube-Events to be consumed by users.

> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents
* [ExporterSinks](#exportersinks)
* [ExporterWebhookSink](#exporterwebhooksink)
* [KubeEventsExporter](#kubeeventsexporter)
* [KubeEventsExporterList](#kubeeventsexporterlist)
* [KubeEventsExporterSpec](#kubeeventsexporterspec)
* [KubeEventsRuler](#kubeeventsruler)
* [KubeEventsRulerList](#kubeeventsrulerlist)
* [KubeEventsRulerSpec](#kubeeventsrulerspec)
* [RulerAlertmanagerSink](#ruleralertmanagersink)
* [RulerSinks](#rulersinks)
* [RulerStdoutSink](#rulerstdoutsink)
* [RulerWebhookSink](#rulerwebhooksink)
* [ServiceReference](#servicereference)
* [KubeEventsRule](#kubeeventsrule)
* [KubeEventsRuleList](#kubeeventsrulelist)
* [KubeEventsRuleSpec](#kubeeventsrulespec)
* [Rule](#rule)

## ExporterSinks

ExporterSinks defines a set of sinks for Events Exporter

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| webhooks | Webhooks is a list of ExporterWebhookSink | []*[ExporterWebhookSink](#exporterwebhooksink) | false |
| stdout | Stdout represents whether to write events to stdout | *[ExporterStdoutSink](#exporterstdoutsink) | false |

[Back to TOC](#table-of-contents)

## ExporterWebhookSink

ExporterWebhookSink defines parameters for webhook sink of Events Exporter.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| url |  | string | false |
| service |  | *[ServiceReference](#servicereference) | false |

[Back to TOC](#table-of-contents)

## KubeEventsExporter

KubeEventsExporter is the Schema for the kubeeventsexporters API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [KubeEventsExporterSpec](#kubeeventsexporterspec) | false |
| status |  | [KubeEventsExporterStatus](#kubeeventsexporterstatus) | false |

[Back to TOC](#table-of-contents)

## KubeEventsExporterList

KubeEventsExporterList contains a list of KubeEventsExporter

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][KubeEventsExporter](#kubeeventsexporter) | true |

[Back to TOC](#table-of-contents)

## KubeEventsExporterSpec

KubeEventsExporterSpec defines the desired state of KubeEventsExporter

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image |  | string | false |
| imagePullPolicy |  | corev1.PullPolicy | false |
| resources | Resources defines resources requests and limits for single Pod. | corev1.ResourceRequirements | false |
| sinks | Sinks defines details of events sinks | *[ExporterSinks](#exportersinks) | false |

[Back to TOC](#table-of-contents)

## KubeEventsRuler

KubeEventsRuler is the Schema for the kubeeventsrulers API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [KubeEventsRulerSpec](#kubeeventsrulerspec) | false |
| status |  | [KubeEventsRulerStatus](#kubeeventsrulerstatus) | false |

[Back to TOC](#table-of-contents)

## KubeEventsRulerList

KubeEventsRulerList contains a list of KubeEventsRuler

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][KubeEventsRuler](#kubeeventsruler) | true |

[Back to TOC](#table-of-contents)

## KubeEventsRulerSpec

KubeEventsRulerSpec defines the desired state of KubeEventsRuler

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| replicas |  | *int32 | false |
| image |  | string | false |
| imagePullPolicy |  | corev1.PullPolicy | false |
| resources | Resources defines resources requests and limits for single Pod. | corev1.ResourceRequirements | false |
| ruleNamespaceSelector | Namespaces to be selected for KubeEventsRules discovery. If unspecified, discover KubeEventsRule instances from all namespaces. | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) | false |
| ruleSelector | A selector to select KubeEventsRules instances. | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) | false |
| sinks | Sinks defines sinks detail of this ruler | *[RulerSinks](#rulersinks) | false |

[Back to TOC](#table-of-contents)

## RulerAlertmanagerSink

RulerAlertmanagerSink is a sink to alertmanager service on k8s

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| namespace |  | string | false |
| name |  | string | false |
| port |  | *int | false |
| targetPort | TargetPort is the port to access on the backend instances targeted by the service. If this is not specified, the value of the 'port' field is used. | *int | false |

[Back to TOC](#table-of-contents)

## RulerSinks



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| alertmanager | Alertmanager is for sinking alerts | *[RulerAlertmanagerSink](#ruleralertmanagersink) | false |
| webhooks | Webhooks is a list of RulerWebhookSink to which notifications or alerts can sink | []*[RulerWebhookSink](#rulerwebhooksink) | false |
| stdout |  | *[RulerStdoutSink](#rulerstdoutsink) | false |

[Back to TOC](#table-of-contents)

## RulerStdoutSink

RulerStdoutSink defines parameters for stdout sink of Events Ruler.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type |  | RulerSinkType | false |

[Back to TOC](#table-of-contents)

## RulerWebhookSink

RulerWebhookSink defines parameters for webhook sink of Events Ruler.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type | Type represents that the sink is for notification or alert. Available values are `notification` and `alert` | RulerSinkType | false |
| namespace |  | string | false |
| service |  | *[ServiceReference](#servicereference) | false |

[Back to TOC](#table-of-contents)

## ServiceReference

ServiceReference holds a reference to k8s Service

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| namespace |  | string | false |
| name |  | string | false |
| port |  | *int | false |
| path |  | string | false |

[Back to TOC](#table-of-contents)

## KubeEventsRule

KubeEventsRule is the Schema for the kubeeventsrules API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [KubeEventsRuleSpec](#kubeeventsrulespec) | false |
| status |  | [KubeEventsRuleStatus](#kubeeventsrulestatus) | false |

[Back to TOC](#table-of-contents)

## KubeEventsRuleList

KubeEventsRuleList contains a list of KubeEventsRule

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][KubeEventsRule](#kubeeventsrule) | true |

[Back to TOC](#table-of-contents)

## KubeEventsRuleSpec

KubeEventsRuleSpec defines the desired state of KubeEventsRule

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| rules |  | [][Rule](#rule) | false |

[Back to TOC](#table-of-contents)

## Rule

Rule describes a notification or alert rule

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name |  | string | false |
| summary |  | string | false |
| summaryCn |  | string | false |
| condition | Condition is a string similar with the where part of sql. The usage is as follows: event.type=\"Warning\" and event.involvedObject.kind=\"Pod\" and event.reason=\"FailedMount\" | string | false |
| message |  | string | false |
| priority |  | string | false |
| source |  | string | false |
| tags |  | []string | false |
| enable | Enable is whether to enable the rule | bool | false |
| type | Type represents that the rule is for notification or alert. Available values are `notification` and `alert` | RuleType | false |

[Back to TOC](#table-of-contents)

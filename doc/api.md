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
| stdout | Stdout represents whether to write events to stdout. Output when configure an empty struct `{}`, but do nothing when no configuration | *[ExporterStdoutSink](#exporterstdoutsink) | false |

[Back to TOC](#table-of-contents)

## ExporterWebhookSink

ExporterWebhookSink defines parameters for webhook sink of Events Exporter.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| url | `url` gives the location of the webhook, in standard URL form (`scheme://host:port/path`). Exactly one of `url` or `service` must be specified. | string | false |
| service | `service` is a reference to the service for this webhook. Either `service` or `url` must be specified. If the webhook is running within the cluster, then you should use `service`. | *[ServiceReference](#servicereference) | false |

[Back to TOC](#table-of-contents)

## KubeEventsExporter

KubeEventsExporter is the Schema for the kubeeventsexporters API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec | Spec defines the specification of the desired behavior of the KubeEventsExporter. | [KubeEventsExporterSpec](#kubeeventsexporterspec) | true |
| status |  | [KubeEventsExporterStatus](#kubeeventsexporterstatus) | false |

[Back to TOC](#table-of-contents)

## KubeEventsExporterList

KubeEventsExporterList contains a list of KubeEventsExporter

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items | List of KubeEventsExporters | [][KubeEventsExporter](#kubeeventsexporter) | true |

[Back to TOC](#table-of-contents)

## KubeEventsExporterSpec

KubeEventsExporterSpec defines the desired state of KubeEventsExporter

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Docker image of kube-events-exporter | string | true |
| imagePullPolicy | Image pull policy. One of Always, Never, IfNotPresent. | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#container-v1-core) | false |
| resources | Resources defines resources requests and limits for single Pod. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core) | false |
| sinks | Sinks defines details of events sinks | *[ExporterSinks](#exportersinks) | false |

[Back to TOC](#table-of-contents)

## KubeEventsRuler

KubeEventsRuler is the Schema for the kubeeventsrulers API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec | Spec defines the specification of the desired behavior of the KubeEventsRuler. | [KubeEventsRulerSpec](#kubeeventsrulerspec) | true |
| status |  | [KubeEventsRulerStatus](#kubeeventsrulerstatus) | false |

[Back to TOC](#table-of-contents)

## KubeEventsRulerList

KubeEventsRulerList contains a list of KubeEventsRuler

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items | List of KubeEventsRulers | [][KubeEventsRuler](#kubeeventsruler) | true |

[Back to TOC](#table-of-contents)

## KubeEventsRulerSpec

KubeEventsRulerSpec defines the desired state of KubeEventsRuler

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| replicas | Number of desired pods. Defaults to 1. | *int32 | false |
| image | Docker image of kube-events-exporter | string | true |
| imagePullPolicy | Image pull policy. One of Always, Never, IfNotPresent. | [corev1.PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#container-v1-core) | false |
| resources | Resources defines resources requests and limits for single Pod. | [corev1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core) | false |
| ruleNamespaceSelector | Namespaces to be selected for KubeEventsRules discovery. If unspecified, discover KubeEventsRule instances from all namespaces. | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) | false |
| ruleSelector | A selector to select KubeEventsRules instances. | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#labelselector-v1-meta) | false |
| sinks | Sinks defines sinks detail of this ruler | *[RulerSinks](#rulersinks) | false |

[Back to TOC](#table-of-contents)

## RulerAlertmanagerSink

RulerAlertmanagerSink is a sink to alertmanager service on k8s

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| namespace | `namespace` is the namespace of the alertmanager service. | string | true |
| name | `name` is the name of the alertmanager service. | string | true |
| port | `port` is the port on the alertmanager service. Default to 9093. `port` should be a valid port number (1-65535, inclusive). | *int | false |
| targetPort | TargetPort is the port to access on the backend instances targeted by the alertmanager service. If this is not specified, the value of the 'port' field is used. | *int | false |

[Back to TOC](#table-of-contents)

## RulerSinks

RulerSinks defines a set of sinks for Events Ruler

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| alertmanager | Alertmanager is an alertmanager sink to which only alerts can sink. | *[RulerAlertmanagerSink](#ruleralertmanagersink) | false |
| webhooks | Webhooks is a list of RulerWebhookSink to which notifications or alerts can sink | []*[RulerWebhookSink](#rulerwebhooksink) | false |
| stdout | Stdout can config write notifications or alerts to stdout; do nothing when no configuration | *[RulerStdoutSink](#rulerstdoutsink) | false |

[Back to TOC](#table-of-contents)

## RulerStdoutSink

RulerStdoutSink defines parameters for stdout sink of Events Ruler.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type | Type represents that the sink is for notification or alert. Available values are `notification` and `alert` | RulerSinkType | true |

[Back to TOC](#table-of-contents)

## RulerWebhookSink

RulerWebhookSink defines parameters for webhook sink of Events Ruler.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type | Type represents that the sink is for notification or alert. Available values are `notification` and `alert` | RulerSinkType | true |
| url | `url` gives the location of the webhook, in standard URL form (`scheme://host:port/path`). Exactly one of `url` or `service` must be specified. | string | false |
| service | `service` is a reference to the service for this webhook. Either `service` or `url` must be specified. If the webhook is running within the cluster, then you should use `service`. | *[ServiceReference](#servicereference) | false |

[Back to TOC](#table-of-contents)

## ServiceReference

ServiceReference holds a reference to k8s Service

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| namespace | `namespace` is the namespace of the service. | string | true |
| name | `name` is the name of the service. | string | true |
| port | `port` is the port on the service and should be a valid port number (1-65535, inclusive). | *int | false |
| path | `path` is an optional URL path which will be sent in any request to this service. | string | false |

[Back to TOC](#table-of-contents)

## KubeEventsRule

KubeEventsRule is the Schema for the kubeeventsrules API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [KubeEventsRuleSpec](#kubeeventsrulespec) | true |
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
| name | Name is simple name of rule | string | false |
| condition | Condition is a string similar with the where part of sql (please use double quotation to mark a string). For example: `event.type="Warning" and event.involvedObject.kind="Pod" and event.reason="FailedMount"` | string | false |
| labels | Labels | map[string]string | false |
| annotations | Values of Annotations can use format string with the fields of the event. For example: `{"message": "%event.message"}` | map[string]string | false |
| enable | Enable is whether to enable the rule | bool | false |
| type | Type represents that the rule is for notification or alert. Available values are `notification` and `alert` | RuleType | false |

[Back to TOC](#table-of-contents)

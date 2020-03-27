# Kube-Events

Kube-Events observes events, exports them, filters them by related rules, and then notifies remaining events through webhooks, or processing them as alerts which will eventually be sent to alertmanager or webhooks. It contains a complete set of components: KubeEvents Exporter (events watching and exporting), KubeEvents Ruler (events filtering, notification and alerting). They depend on CRDs (`KubeEventsExporter`, `KubeEventsRuler`, `KubeEventsRule`) and an integrated Operator. 

# CustomResourceDefinitions

The Operator acts on the following CRDs:

- **`KubeEventsExporter`**, which defines a desired events Exporter deployment. The Operator ensures at all times that a deployment matching the resource definition is running.
- **`KubeEventsRuler`**, which defines a desired events Ruler deployment. The Operator ensures at all times that a deployment matching the resource definition is running.
- **`KubeEventsRules`**, which defines a desired events rule set, which can be used to filter events by the Ruler. 

# QuickStart

You can install the operator in any kubernetes cluster with following commands:

```shell
kubectl apply -f  ./deploy/operator.yaml
```

Then install default rules which have cluster scope:

```shell
kubectl apply -f  ./deploy/crs-rules-default.yaml
```

And install exporter and ruler:

```shell
kubectl apply -f  ./deploy/crs.yaml
```
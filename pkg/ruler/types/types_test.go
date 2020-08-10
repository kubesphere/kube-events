package types

import (
	"encoding/json"
	"github.com/kubesphere/kube-events/pkg/apis/v1alpha1"
	"k8s.io/api/core/v1"
	"testing"
)

func TestEvalByRule(t *testing.T) {
	evtString := `
{
    "metadata": {
        "name": "kube-events-ruler-64c959d7ff-m4gmd.16292df7c28120f6",
        "namespace": "default",
        "selfLink": "/api/v1/namespaces/default/events/kube-events-ruler-64c959d7ff-m4gmd.16292df7c28120f6",
        "uid": "46a523cb-62dd-4cc0-92fd-d2ac879deb0c",
        "resourceVersion": "2536141",
        "creationTimestamp": "2020-08-08T03:41:25Z"
    },
    "involvedObject": {
        "kind": "Pod",
        "namespace": "default",
        "name": "kube-events-ruler-64c959d7ff-m4gmd",
        "uid": "d3543e69-bdf3-4777-9b4b-97dc4058f18c",
        "apiVersion": "v1",
        "resourceVersion": "2535341",
        "fieldPath": "spec.containers{config-reloader}"
    },
    "reason": "Killing",
    "message": "Stopping container config-reloader",
    "source": {
        "component": "kubelet",
        "host": "i-ywi8iwim"
    },
    "firstTimestamp": "2020-08-08T03:41:25Z",
    "lastTimestamp": "2020-08-08T03:41:25Z",
    "count": 1,
    "type": "Normal",
    "eventTime": null,
    "reportingComponent": "",
    "reportingInstance": ""
}
`
	var evt v1.Event
	err := json.Unmarshal([]byte(evtString), &evt)
	if err != nil {
		t.Error(err)
	}
	wEvt := Event{Event: &evt}
	ok, err := wEvt.EvalByRule(&v1alpha1.EventRule{
		Name:      "Test",
		Condition: `metadata.name regex ".*-events-.*"`,
		Enable:    true,
		Type:      "alert",
	})
	if err != nil {
		t.Error(err)
	}
	if !ok {
		t.Fatalf("expected: %v, but got %v", true, ok)
	}
}

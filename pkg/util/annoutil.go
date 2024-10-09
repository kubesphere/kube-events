package util

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

const (
	EventVersionOld = "corev1"
	EventVersionNew = "eventsv1"
)

var cluster string
var EventVersion string

func SetClusterName(client *kubernetes.Clientset) {
	setCluster(client)
	t := time.NewTicker(60 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if cluster != "" {
				return
			}
			setCluster(client)
			klog.Infof("current cluster is [%s]", GetCluster())
		}
	}

}

func setCluster(client *kubernetes.Clientset) {

	ns, err := client.CoreV1().Namespaces().Get(context.Background(), "kubesphere-system", metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get namespace kubesphere-system error: %s", err)
		return
	}

	if ns.Annotations != nil {
		cluster = ns.Annotations["cluster.kubesphere.io/name"]
	}

}

func GetCluster() string {
	return cluster
}

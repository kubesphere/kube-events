package util

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

var cluster string
var NewEventType bool

func SetClusterName(client *kubernetes.Clientset) {
	setCluster(client)
	c := time.Tick(60 * time.Second)
	for {
		select {
		case <-c:
			setCluster(client)
			klog.Infof("current cluster is [%s]", GetCluster())
		}
	}

}

func setCluster(client *kubernetes.Clientset) {

	ns, err := client.CoreV1().Namespaces().Get(context.Background(), "kubesphere-system", metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get namespace kubesphere-system error: %s", err)
	}

	if ns.Annotations != nil {
		cluster = ns.Annotations["cluster.kubesphere.io/name"]
	}

}

func GetCluster() string {
	return cluster
}

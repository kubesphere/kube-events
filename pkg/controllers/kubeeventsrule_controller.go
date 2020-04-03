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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	loggingv1alpha1 "github.com/kubesphere/kube-events/pkg/apis/v1alpha1"
)

// KubeEventsRuleReconciler reconciles a KubeEventsRule object
type KubeEventsRuleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=logging.kubesphere.io,resources=kubeeventsrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=logging.kubesphere.io,resources=kubeeventsrules/status,verbs=get;update;patch

func (r *KubeEventsRuleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("kubeeventsrule", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *KubeEventsRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loggingv1alpha1.KubeEventsRule{}).
		Complete(r)
}

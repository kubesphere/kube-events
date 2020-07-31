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
	"fmt"
	"path"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/kubesphere/kube-events/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/yaml"

	eventsv1alpha1 "github.com/kubesphere/kube-events/pkg/apis/v1alpha1"
)

const (
	configDirEventsRuler      = "/etc/events-ruler/config"
	configFileNameEventsRuler = "events-ruler.yaml"

	labelKeyEventsRuler          = "events-ruler"
	labelKeyEventsRulerNamespace = "events-ruler-ns"

	finalizerNameEventsRuler = "rulers.finalizer.events.kubesphere.io"
)

// RulerReconciler reconciles a Ruler object
type RulerReconciler struct {
	Conf *config.OperatorConfig
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=events.kubesphere.io,resources=rulers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=events.kubesphere.io,resources=rulers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=services/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch

func (r *RulerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("ruler", req.NamespacedName)

	ker := &eventsv1alpha1.Ruler{}
	e := r.Get(ctx, req.NamespacedName, ker)

	if e != nil {
		if apierrs.IsNotFound(e) {
			return ctrl.Result{}, nil
		}
		log.Error(e, "unable to fetch Ruler")
		return ctrl.Result{}, e
	}

	if ker.ObjectMeta.DeletionTimestamp.IsZero() {
		if !hasFinalizer(&ker.ObjectMeta, finalizerNameEventsRuler) {
			controllerutil.AddFinalizer(&ker.ObjectMeta, finalizerNameEventsRuler)
			if e = r.Update(ctx, ker); e != nil {
				return ctrl.Result{}, e
			}
		}
	} else {
		if hasFinalizer(&ker.ObjectMeta, finalizerNameEventsRuler) {
			crb := &rbacv1.ClusterRoleBinding{}
			crb.Name = fmt.Sprintf("%s-%s", ker.Namespace, ker.Name)
			if e = r.Delete(ctx, crb); e != nil && !apierrs.IsNotFound(e) {
				return ctrl.Result{}, e
			}
			cr := &rbacv1.ClusterRole{}
			cr.Name = fmt.Sprintf("%s-%s", ker.Namespace, ker.Name)
			if e = r.Delete(ctx, cr); e != nil && !apierrs.IsNotFound(e) {
				return ctrl.Result{}, e
			}
			controllerutil.RemoveFinalizer(&ker.ObjectMeta, finalizerNameEventsRuler)
			if e = r.Update(ctx, ker); e != nil {
				return ctrl.Result{}, e
			}
		}
		return ctrl.Result{}, nil
	}

	sa := &corev1.ServiceAccount{}
	sa.Namespace = ker.Namespace
	sa.Name = ker.Name
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, sa, r.serviceAccountMutate(sa, ker)); e != nil {
		return ctrl.Result{}, e
	}
	cr := &rbacv1.ClusterRole{}
	cr.Name = fmt.Sprintf("%s-%s", ker.Namespace, ker.Name)
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, cr, r.clusterRoleMutate(cr, ker)); e != nil {
		return ctrl.Result{}, e
	}
	crb := &rbacv1.ClusterRoleBinding{}
	crb.Name = fmt.Sprintf("%s-%s", ker.Namespace, ker.Name)
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, crb, r.clusterRoleBindingMutate(crb, cr, sa, ker)); e != nil {
		return ctrl.Result{}, e
	}
	cm := &corev1.ConfigMap{}
	cm.Namespace = ker.Namespace
	cm.Name = ker.Name
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, cm, r.configMutate(cm, ker)); e != nil {
		return ctrl.Result{}, e
	}
	deploy := &appsv1.Deployment{}
	deploy.Name = ker.Name
	deploy.Namespace = ker.Namespace
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, deploy, r.deployMutate(deploy, cm, sa, ker)); e != nil {
		return ctrl.Result{}, e
	}
	svc := &corev1.Service{}
	svc.Name = ker.Name
	svc.Namespace = ker.Namespace
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, svc, r.serviceMutate(svc, deploy, ker)); e != nil {
		return ctrl.Result{}, e
	}

	return ctrl.Result{}, nil
}

func (r *RulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	meets := func(meta metav1.Object, obj runtime.Object) bool {
		if meta == nil || obj == nil {
			return false
		}
		switch obj.(type) {
		case *eventsv1alpha1.Ruler:
			return true
		case *appsv1.Deployment, *corev1.Service, *corev1.ConfigMap, *corev1.ServiceAccount:
			if ls := meta.GetLabels(); ls != nil {
				_, ok := ls[labelKeyEventsRuler]
				return ok
			}
			return false
		case *rbacv1.ClusterRole, *rbacv1.ClusterRoleBinding:
			if ls := meta.GetLabels(); ls != nil {
				_, ok1 := ls[labelKeyEventsRuler]
				_, ok2 := ls[labelKeyEventsRulerNamespace]
				return ok1 && ok2
			}
			return false
		}
		return false
	}
	preficateFuncs := predicate.Funcs{
		CreateFunc: func(event event.CreateEvent) bool {
			return meets(event.Meta, event.Object)
		},
		UpdateFunc: func(event event.UpdateEvent) bool {
			if meets(event.MetaOld, event.ObjectOld) {
				if event.MetaOld != nil && event.MetaNew != nil {
					return event.MetaOld.GetResourceVersion() != event.MetaNew.GetResourceVersion()
				}
			}
			return false
		},
		DeleteFunc: func(event event.DeleteEvent) bool {
			return meets(event.Meta, event.Object)
		},
	}
	enq := func(meta metav1.Object, q workqueue.RateLimitingInterface) {
		if ls := meta.GetLabels(); ls != nil {
			name, ok1 := ls[labelKeyEventsRuler]
			ns, ok2 := ls[labelKeyEventsRulerNamespace]
			if ok1 && ok2 {
				q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
					Name:      name,
					Namespace: ns,
				}})
			}
		}
	}
	chandler := handler.Funcs{
		CreateFunc: func(createEvent event.CreateEvent, limitingInterface workqueue.RateLimitingInterface) {
			if meta := createEvent.Meta; meta != nil {
				enq(meta, limitingInterface)
			}
		},
		UpdateFunc: func(updateEvent event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
			if meta := updateEvent.MetaOld; meta != nil {
				enq(meta, limitingInterface)
			}
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&eventsv1alpha1.Ruler{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.ServiceAccount{}).
		Watches(&source.Kind{Type: &rbacv1.ClusterRole{}}, chandler).
		Watches(&source.Kind{Type: &rbacv1.ClusterRoleBinding{}}, chandler).
		WithEventFilter(&preficateFuncs).
		Complete(r)
}

func (r *RulerReconciler) serviceAccountMutate(sa *corev1.ServiceAccount,
	ker *eventsv1alpha1.Ruler) controllerutil.MutateFn {
	return func() error {
		sa.Labels = r.relativeResourcesShareLabels(ker)
		sa.SetOwnerReferences(nil)
		return controllerutil.SetControllerReference(ker, sa, r.Scheme)
	}
}

func (r *RulerReconciler) clusterRoleMutate(cr *rbacv1.ClusterRole,
	ker *eventsv1alpha1.Ruler) controllerutil.MutateFn {
	return func() error {
		cr.Labels = r.relativeResourcesShareLabels(ker)
		cr.Labels[labelKeyEventsRulerNamespace] = ker.Namespace
		cr.Rules = []rbacv1.PolicyRule{{
			APIGroups: []string{"events.kubesphere.io"},
			Resources: []string{"rules"},
			Verbs:     []string{"get", "list", "watch"},
		}, {
			APIGroups: []string{""},
			Resources: []string{"namespaces"},
			Verbs:     []string{"get", "list", "watch"},
		}}
		return nil
	}
}

func (r *RulerReconciler) clusterRoleBindingMutate(crb *rbacv1.ClusterRoleBinding,
	cr *rbacv1.ClusterRole, sa *corev1.ServiceAccount,
	ker *eventsv1alpha1.Ruler) controllerutil.MutateFn {
	return func() error {
		crb.Labels = r.relativeResourcesShareLabels(ker)
		crb.Labels[labelKeyEventsRulerNamespace] = ker.Namespace
		crb.RoleRef = rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     cr.Name,
		}
		crb.Subjects = []rbacv1.Subject{{
			Kind:      rbacv1.ServiceAccountKind,
			Name:      sa.Name,
			Namespace: sa.Namespace,
		}}
		return nil
	}
}

func (r *RulerReconciler) configMutate(cm *corev1.ConfigMap,
	ker *eventsv1alpha1.Ruler) controllerutil.MutateFn {
	return func() error {
		cm.Labels = r.relativeResourcesShareLabels(ker)

		c := config.RulerConfig{
			Sinks: &config.RulerSinks{},
		}
		if ss := ker.Spec.Sinks; ss != nil {
			if ss.Alertmanager != nil {
				c.Sinks.Alertmanager = &config.RulerAlertmanagerSink{
					Namespace:  ss.Alertmanager.Namespace,
					Name:       ss.Alertmanager.Name,
					Port:       ss.Alertmanager.Port,
					TargetPort: ss.Alertmanager.TargetPort,
				}
			}
			if ss.Stdout != nil {
				c.Sinks.Stdout = &config.RulerStdoutSink{
					Type: config.RulerSinkType(ss.Stdout.Type),
				}
			}
			for _, webhook := range ss.Webhooks {
				cw := &config.RulerWebhookSink{
					Type: config.RulerSinkType(string(webhook.Type)),
				}
				if webhook.Url != "" {
					cw.Url = webhook.Url
					c.Sinks.Webhooks = append(c.Sinks.Webhooks, cw)
				} else if webhook.Service != nil {
					cw.Service = &config.ServiceReference{
						Namespace: webhook.Service.Namespace,
						Name:      webhook.Service.Name,
						Port:      webhook.Service.Port,
						Path:      webhook.Service.Path,
					}
					c.Sinks.Webhooks = append(c.Sinks.Webhooks, cw)
				}
			}
		}

		bs, e := yaml.Marshal(c)
		if e != nil {
			return e
		}
		cm.Data = map[string]string{
			configFileNameEventsRuler: string(bs),
		}
		cm.SetOwnerReferences(nil)
		return controllerutil.SetControllerReference(ker, cm, r.Scheme)
	}
}

func (r *RulerReconciler) deployMutate(deploy *appsv1.Deployment,
	cm *corev1.ConfigMap, sa *corev1.ServiceAccount,
	ker *eventsv1alpha1.Ruler) controllerutil.MutateFn {
	return func() error {
		deploy.Labels = r.relativeResourcesShareLabels(ker)

		replicas := int32(1)
		if ker.Spec.Replicas != nil {
			replicas = *ker.Spec.Replicas
		}
		deploy.Spec.Replicas = &replicas

		podLabels := r.relativeResourcesShareLabels(ker)
		podLabels["app"] = deploy.Name
		deploy.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: podLabels,
		}
		deploy.Spec.Template.Labels = podLabels
		deploy.Spec.Template.Spec.ServiceAccountName = sa.Name

		expcConfV := corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cm.Name,
					},
				},
			},
		}
		var confV corev1.Volume
		for _, v := range deploy.Spec.Template.Spec.Volumes {
			if v.Name == expcConfV.Name {
				if v.ConfigMap != nil {
					confV.Name = v.Name
					confV.ConfigMap = &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: v.ConfigMap.Name,
						},
					}
				}
			}
		}
		if !reflect.DeepEqual(expcConfV, confV) {
			deploy.Spec.Template.Spec.Volumes = []corev1.Volume{expcConfV}
		}

		reloaderRes := corev1.ResourceRequirements{Limits: corev1.ResourceList{}}
		if r.Conf.ConfigReloaderCPU != "0" {
			reloaderRes.Limits[corev1.ResourceCPU] = resource.MustParse(r.Conf.ConfigReloaderCPU)
		}
		if r.Conf.ConfigReloaderMemory != "0" {
			reloaderRes.Limits[corev1.ResourceMemory] = resource.MustParse(r.Conf.ConfigReloaderMemory)
		}
		expcRulerC := corev1.Container{
			Name: "events-ruler",
			Args: []string{
				fmt.Sprintf("--config.file=%s", path.Join(configDirEventsRuler, configFileNameEventsRuler)),
			},
			Image:           ker.Spec.Image,
			ImagePullPolicy: ker.Spec.ImagePullPolicy,
			Resources:       ker.Spec.Resources,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      expcConfV.Name,
					MountPath: configDirEventsRuler,
				},
			},
		}
		expcReloaderC := corev1.Container{
			Name:      "config-reloader",
			Image:     r.Conf.ConfigReloaderImage,
			Resources: reloaderRes,
			Args: []string{
				fmt.Sprintf("--volume-dir=%s", configDirEventsRuler),
				"--webhook-url=http://127.0.0.1:8443/-/reload",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      expcConfV.Name,
					MountPath: configDirEventsRuler,
				},
			},
		}
		var rulerC, reloaderC corev1.Container
		for _, c := range deploy.Spec.Template.Spec.Containers {
			simplec := corev1.Container{
				Name:         c.Name,
				Image:        c.Image,
				Resources:    c.Resources,
				Args:         c.Args,
				VolumeMounts: c.VolumeMounts,
			}
			if simplec.Name == expcRulerC.Name {
				rulerC = simplec
			} else if simplec.Name == expcReloaderC.Name {
				reloaderC = simplec
			}
		}
		if !reflect.DeepEqual(expcRulerC, rulerC) ||
			!reflect.DeepEqual(expcReloaderC, reloaderC) {
			deploy.Spec.Template.Spec.Containers = []corev1.Container{expcRulerC, expcReloaderC}
		}
		deploy.SetOwnerReferences(nil)
		return controllerutil.SetControllerReference(ker, deploy, r.Scheme)
	}
}

func (r *RulerReconciler) serviceMutate(svc *corev1.Service, deploy *appsv1.Deployment,
	ker *eventsv1alpha1.Ruler) controllerutil.MutateFn {
	return func() error {
		svc.Labels = r.relativeResourcesShareLabels(ker)

		podLabels := r.relativeResourcesShareLabels(ker)
		podLabels["app"] = deploy.Name
		svc.Spec.ClusterIP = "None"
		svc.Spec.Ports = []corev1.ServicePort{
			{
				Port:       8443,
				TargetPort: intstr.FromInt(8443),
				Protocol:   corev1.ProtocolTCP,
			},
		}
		svc.Spec.Selector = podLabels
		svc.SetOwnerReferences(nil)
		return controllerutil.SetControllerReference(ker, svc, r.Scheme)
	}
}

func (r *RulerReconciler) relativeResourcesShareLabels(ker *eventsv1alpha1.Ruler) map[string]string {
	ls := make(map[string]string)
	ls[labelKeyEventsRuler] = ker.Name
	return ls
}

func hasFinalizer(o metav1.Object, finalizer string) bool {
	f := o.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return true
		}
	}
	return false
}

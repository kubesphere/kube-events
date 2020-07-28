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

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/yaml"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	eventsv1alpha1 "github.com/kubesphere/kube-events/pkg/apis/v1alpha1"
	"github.com/kubesphere/kube-events/pkg/config"
)

const (
	configDirEventsExporter      = "/etc/events-exporter/config"
	configFileNameEventsExporter = "events-exporter.yaml"

	labelKeyEventsExporter          = "events-exporter"
	labelKeyEventsExporterNamespace = "events-exporter-ns"

	finalizerNameEventsExporter = "exporters.finalizer.events.kubesphere.io"
)

// ExporterReconciler reconciles a Exporter object
type ExporterReconciler struct {
	Conf *config.OperatorConfig
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=events.kubesphere.io,resources=exporters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=events.kubesphere.io,resources=exporters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=clusterrolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch

func (r *ExporterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("exporter", req.NamespacedName)

	kee := &eventsv1alpha1.Exporter{}
	e := r.Get(ctx, req.NamespacedName, kee)
	if e != nil {
		if apierrs.IsNotFound(e) {
			return ctrl.Result{}, nil
		}
		log.Error(e, "unable to fetch Exporter")
		return ctrl.Result{}, e
	}

	if kee.ObjectMeta.DeletionTimestamp.IsZero() {
		if !hasFinalizer(&kee.ObjectMeta, finalizerNameEventsExporter) {
			controllerutil.AddFinalizer(&kee.ObjectMeta, finalizerNameEventsExporter)
			if e = r.Update(ctx, kee); e != nil {
				return ctrl.Result{}, e
			}
		}
	} else {
		if hasFinalizer(&kee.ObjectMeta, finalizerNameEventsExporter) {
			crb := &rbacv1.ClusterRoleBinding{}
			crb.Name = fmt.Sprintf("%s-%s", kee.Namespace, kee.Name)
			if e = r.Delete(ctx, crb); e != nil && !apierrs.IsNotFound(e) {
				return ctrl.Result{}, e
			}
			cr := &rbacv1.ClusterRole{}
			cr.Name = fmt.Sprintf("%s-%s", kee.Namespace, kee.Name)
			if e = r.Delete(ctx, cr); e != nil && !apierrs.IsNotFound(e) {
				return ctrl.Result{}, e
			}
			controllerutil.RemoveFinalizer(&kee.ObjectMeta, finalizerNameEventsExporter)
			if e = r.Update(ctx, kee); e != nil {
				return ctrl.Result{}, e
			}
		}
		return ctrl.Result{}, nil
	}

	sa := &corev1.ServiceAccount{}
	sa.Namespace = kee.Namespace
	sa.Name = kee.Name
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, sa, r.serviceAccountMutate(sa, kee)); e != nil {
		return ctrl.Result{}, e
	}
	cr := &rbacv1.ClusterRole{}
	cr.Name = fmt.Sprintf("%s-%s", kee.Namespace, kee.Name)
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, cr, r.clusterRoleMutate(cr, kee)); e != nil {
		return ctrl.Result{}, e
	}
	crb := &rbacv1.ClusterRoleBinding{}
	crb.Name = fmt.Sprintf("%s-%s", kee.Namespace, kee.Name)
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, crb, r.clusterRoleBindingMutate(crb, cr, sa, kee)); e != nil {
		return ctrl.Result{}, e
	}
	cm := &corev1.ConfigMap{}
	cm.Namespace = kee.Namespace
	cm.Name = kee.Name
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, cm, r.configMutate(cm, kee)); e != nil {
		return ctrl.Result{}, e
	}
	deploy := &appsv1.Deployment{}
	deploy.Name = kee.Name
	deploy.Namespace = kee.Namespace
	if _, e = controllerutil.CreateOrUpdate(ctx, r.Client, deploy, r.deployMutate(deploy, cm, sa, kee)); e != nil {
		return ctrl.Result{}, e
	}

	return ctrl.Result{}, nil
}

func (r *ExporterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	meets := func(meta metav1.Object, obj runtime.Object) bool {
		if meta == nil || obj == nil {
			return false
		}
		switch obj.(type) {
		case *eventsv1alpha1.Exporter:
			return true
		case *appsv1.Deployment, *corev1.Service, *corev1.ConfigMap, *corev1.ServiceAccount:
			if ls := meta.GetLabels(); ls != nil {
				_, ok := ls[labelKeyEventsExporter]
				return ok
			}
			return false
		case *rbacv1.ClusterRole, *rbacv1.ClusterRoleBinding:
			if ls := meta.GetLabels(); ls != nil {
				_, ok1 := ls[labelKeyEventsExporter]
				_, ok2 := ls[labelKeyEventsExporterNamespace]
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
		For(&eventsv1alpha1.Exporter{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.ServiceAccount{}).
		Watches(&source.Kind{Type: &rbacv1.ClusterRole{}}, chandler).
		Watches(&source.Kind{Type: &rbacv1.ClusterRoleBinding{}}, chandler).
		WithEventFilter(&preficateFuncs).
		Complete(r)
}

func (r *ExporterReconciler) serviceAccountMutate(sa *corev1.ServiceAccount,
	kee *eventsv1alpha1.Exporter) controllerutil.MutateFn {
	return func() error {
		sa.Labels = r.relativeResourcesShareLabels(kee)
		sa.SetOwnerReferences(nil)
		return controllerutil.SetControllerReference(kee, sa, r.Scheme)
	}
}

func (r *ExporterReconciler) clusterRoleMutate(cr *rbacv1.ClusterRole,
	kee *eventsv1alpha1.Exporter) controllerutil.MutateFn {
	return func() error {
		cr.Labels = r.relativeResourcesShareLabels(kee)
		cr.Labels[labelKeyEventsExporterNamespace] = kee.Namespace
		cr.Rules = []rbacv1.PolicyRule{{
			APIGroups: []string{""},
			Resources: []string{"events"},
			Verbs:     []string{"get", "list", "watch"},
		}}
		return nil
	}
}

func (r *ExporterReconciler) clusterRoleBindingMutate(crb *rbacv1.ClusterRoleBinding,
	cr *rbacv1.ClusterRole, sa *corev1.ServiceAccount,
	kee *eventsv1alpha1.Exporter) controllerutil.MutateFn {
	return func() error {
		crb.Labels = r.relativeResourcesShareLabels(kee)
		crb.Labels[labelKeyEventsExporterNamespace] = kee.Namespace
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

func (r *ExporterReconciler) configMutate(cm *corev1.ConfigMap,
	kee *eventsv1alpha1.Exporter) controllerutil.MutateFn {
	return func() error {
		cm.Labels = r.relativeResourcesShareLabels(kee)

		c := config.ExporterConfig{
			Sinks: &config.ExporterSinks{},
		}
		if ss := kee.Spec.Sinks; ss != nil {
			if ss.Stdout != nil {
				c.Sinks.Stdout = &config.ExporterSinkStdout{}
			}
			for _, w := range ss.Webhooks {
				cwebhook := &config.ExporterSinkWebhook{}
				if w.Url != "" {
					cwebhook.Url = w.Url
					c.Sinks.Webhooks = append(c.Sinks.Webhooks, cwebhook)
				} else if w.Service != nil {
					cwebhook.Service = &config.ServiceReference{
						Namespace: w.Service.Namespace,
						Name:      w.Service.Name,
						Port:      w.Service.Port,
						Path:      w.Service.Path,
					}
					c.Sinks.Webhooks = append(c.Sinks.Webhooks, cwebhook)
				}
			}
		}

		bs, e := yaml.Marshal(c)
		if e != nil {
			return e
		}
		cm.Data = map[string]string{
			configFileNameEventsExporter: string(bs),
		}
		cm.SetOwnerReferences(nil)
		return controllerutil.SetControllerReference(kee, cm, r.Scheme)
	}
}

func (r *ExporterReconciler) deployMutate(deploy *appsv1.Deployment,
	cm *corev1.ConfigMap, sa *corev1.ServiceAccount,
	kee *eventsv1alpha1.Exporter) controllerutil.MutateFn {
	return func() error {
		deploy.Labels = r.relativeResourcesShareLabels(kee)

		replicas := int32(1)
		deploy.Spec.Replicas = &replicas

		podLabels := r.relativeResourcesShareLabels(kee)
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
		hostTimeV := corev1.Volume{
			Name: "host-time",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/etc/localtime",
				},
			},
		}
		if !reflect.DeepEqual(expcConfV, confV) {
			deploy.Spec.Template.Spec.Volumes = []corev1.Volume{expcConfV, hostTimeV}
		}

		reloaderRes := corev1.ResourceRequirements{Limits: corev1.ResourceList{}}
		if r.Conf.ConfigReloaderCPU != "0" {
			reloaderRes.Limits[corev1.ResourceCPU] = resource.MustParse(r.Conf.ConfigReloaderCPU)
		}
		if r.Conf.ConfigReloaderMemory != "0" {
			reloaderRes.Limits[corev1.ResourceMemory] = resource.MustParse(r.Conf.ConfigReloaderMemory)
		}
		expcExporterC := corev1.Container{
			Name: "events-exporter",
			Args: []string{
				fmt.Sprintf("--config.file=%s", path.Join(configDirEventsExporter, configFileNameEventsExporter)),
			},
			Image:           kee.Spec.Image,
			ImagePullPolicy: kee.Spec.ImagePullPolicy,
			Resources:       kee.Spec.Resources,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      expcConfV.Name,
					MountPath: configDirEventsExporter,
				},
				{
					Name:      hostTimeV.Name,
					MountPath: hostTimeV.HostPath.Path,
				},
			},
		}
		expcReloaderC := corev1.Container{
			Name:      "config-reloader",
			Image:     r.Conf.ConfigReloaderImage,
			Resources: reloaderRes,
			Args: []string{
				fmt.Sprintf("--volume-dir=%s", configDirEventsExporter),
				"--webhook-url=http://127.0.0.1:8443/-/reload",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      expcConfV.Name,
					MountPath: configDirEventsExporter,
				},
				{
					Name:      hostTimeV.Name,
					MountPath: hostTimeV.HostPath.Path,
				},
			},
		}
		var exporterC, reloaderC corev1.Container
		for _, c := range deploy.Spec.Template.Spec.Containers {
			simplec := corev1.Container{
				Name:         c.Name,
				Image:        c.Image,
				Resources:    c.Resources,
				Args:         c.Args,
				VolumeMounts: c.VolumeMounts,
			}
			if simplec.Name == expcExporterC.Name {
				exporterC = simplec
			} else if simplec.Name == expcReloaderC.Name {
				reloaderC = simplec
			}
		}
		if !reflect.DeepEqual(expcExporterC, exporterC) ||
			!reflect.DeepEqual(expcReloaderC, reloaderC) {
			deploy.Spec.Template.Spec.Containers = []corev1.Container{expcExporterC, expcReloaderC}
		}
		deploy.SetOwnerReferences(nil)
		return controllerutil.SetControllerReference(kee, deploy, r.Scheme)
	}
}

func (r *ExporterReconciler) relativeResourcesShareLabels(kee *eventsv1alpha1.Exporter) map[string]string {
	ls := make(map[string]string)
	ls[labelKeyEventsExporter] = kee.Name
	return ls
}

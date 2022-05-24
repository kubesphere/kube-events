package ruler

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/kubesphere/kube-events/pkg/apis/v1alpha1"
	"github.com/kubesphere/kube-events/pkg/config"
	"github.com/kubesphere/kube-events/pkg/ruler/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	kruntime "k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	toolscache "k8s.io/client-go/tools/cache"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

type RuleCache struct {
	clusterRules   sync.Map // {namespace/name: Rule}
	workspaceRules sync.Map // {workspace/namespace/name: Rule}
	namespaceRules sync.Map // {namespace: {namespace+name: Rule}}
	namespaces     sync.Map // {namespace: Namespace}

	ruleSelector          labels.Selector
	ruleNamespaceSelector labels.Selector

	resCache cache.Cache
}

func (c *RuleCache) Run(ctx context.Context) error {
	go c.resCache.Start(ctx)

	nsInf, e := c.resCache.GetInformer(ctx, &corev1.Namespace{})
	if e != nil {
		return e
	}
	var nsHandler toolscache.ResourceEventHandler
	nsHandler = &toolscache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if ns, ok := obj.(*corev1.Namespace); ok {
				c.namespaces.Store(ns.Name, ns)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			nsHandler.OnAdd(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			if ns, ok := obj.(*corev1.Namespace); ok {
				c.namespaces.Delete(ns.Name)
			}
		},
	}
	nsInf.AddEventHandler(nsHandler)
	if ok := toolscache.WaitForCacheSync(ctx.Done(), nsInf.HasSynced); !ok {
		return fmt.Errorf("namespace cache failed")
	}

	ruleInf, e := c.resCache.GetInformer(ctx, &v1alpha1.Rule{})
	if e != nil {
		return e
	}
	ruleInf.AddEventHandler(toolscache.ResourceEventHandlerFuncs{
		AddFunc: c.ruleAdd,
		UpdateFunc: func(oldObj, newObj interface{}) {
			c.ruleDelete(oldObj)
			c.ruleAdd(newObj)
		},
		DeleteFunc: c.ruleDelete,
	})
	if ok := toolscache.WaitForCacheSync(ctx.Done(), ruleInf.HasSynced); !ok {
		return fmt.Errorf("rule cache failed")
	}
	klog.Info("rule cache succeeded")

	return ctx.Err()
}

func (c *RuleCache) GetRules(ctx context.Context, evt *types.Event) (rules []*v1alpha1.Rule) {
	createRangeFunc := func(f func(key interface{}) bool) func(key, value interface{}) bool {
		return func(key, value interface{}) bool {
			if f(key) {
				rules = append(rules, value.(*v1alpha1.Rule))
			}
			return true
		}
	}
	rangeFunc := createRangeFunc(func(key interface{}) bool {
		return true
	})

	c.clusterRules.Range(rangeFunc)

	var ns string
	if evt != nil && evt.Event != nil {
		ns = evt.Event.InvolvedObject.Namespace
	}
	if ns == "" {
		return
	}
	if v, ok := c.namespaces.Load(ns); ok {
		if ns, ok := v.(*corev1.Namespace); ok {
			keyPrefix := c.workspaceFromNamespace(ns) + "/"
			c.workspaceRules.Range(createRangeFunc(func(key interface{}) bool {
				return strings.HasPrefix(key.(string), keyPrefix)
			}))
		}
	}

	if rs, ok := c.namespaceRules.Load(ns); ok {
		rs.(*sync.Map).Range(rangeFunc)
	}

	return
}

func (c *RuleCache) ruleAdd(obj interface{}) {
	rsc := v1alpha1.GetRuleScopeConfig()
	if rule, ok := obj.(*v1alpha1.Rule); ok && c.selectorMatchesRule(rule) {
		switch rule.Labels[rsc.ScopeLabelKey] {
		case rsc.ScopeLabelValueCluster:
			c.clusterRules.Store(c.ruleNameFromClusterRule(rule), rule)
		case rsc.ScopeLabelValueWorkspace:
			c.workspaceRules.Store(c.ruleNameFromWorkspaceRule(rule), rule)
		case rsc.ScopeLabelValueNamespace:
			m, _ := c.namespaceRules.LoadOrStore(rule.Namespace, &sync.Map{})
			m.(*sync.Map).Store(c.ruleNameFromNamespaceRule(rule), rule)
		}
	}
}

func (c *RuleCache) ruleDelete(obj interface{}) {
	rsc := v1alpha1.GetRuleScopeConfig()
	if rule, ok := obj.(*v1alpha1.Rule); ok {
		switch rule.Labels[rsc.ScopeLabelKey] {
		case rsc.ScopeLabelValueCluster:
			c.clusterRules.Delete(c.ruleNameFromClusterRule(rule))
		case rsc.ScopeLabelValueWorkspace:
			c.workspaceRules.Delete(c.ruleNameFromWorkspaceRule(rule))
		case rsc.ScopeLabelValueNamespace:
			if m, ok := c.namespaceRules.Load(rule.Namespace); ok {
				m.(*sync.Map).Delete(c.ruleNameFromNamespaceRule(rule))
			}
		}
	}
}

func (c *RuleCache) selectorMatchesRule(rule *v1alpha1.Rule) bool {
	if !c.ruleSelector.Matches(labels.Set(rule.Labels)) {
		return false
	}
	if ns, ok := c.namespaces.Load(rule.Namespace); ok &&
		c.ruleNamespaceSelector.Matches(labels.Set(ns.(*corev1.Namespace).Labels)) {
		return true
	}
	return false
}

func (c *RuleCache) ruleNameFromClusterRule(rule *v1alpha1.Rule) string {
	return rule.Namespace + "/" + rule.Name
}

func (c *RuleCache) ruleNameFromWorkspaceRule(rule *v1alpha1.Rule) string {
	return c.workspaceFromRule(rule) + "/" + rule.Namespace + "/" + rule.Name
}

func (c *RuleCache) ruleNameFromNamespaceRule(rule *v1alpha1.Rule) string {
	return rule.Namespace + "/" + rule.Name
}

func (c *RuleCache) workspaceFromNamespace(ns *corev1.Namespace) string {
	return ns.Labels[v1alpha1.GetRuleScopeConfig().ScopeWorkspaceLabelKey]
}

func (c *RuleCache) workspaceFromRule(rule *v1alpha1.Rule) string {
	return rule.Labels[v1alpha1.GetRuleScopeConfig().ScopeWorkspaceLabelKey]
}

func NewRuleCache(kcfg *rest.Config, rcfg *config.RulerConfig) (*RuleCache, error) {
	scheme := kruntime.NewScheme()
	_ = v1alpha1.AddToScheme(scheme)
	_ = clientgoscheme.AddToScheme(scheme)
	resCache, e := cache.New(kcfg, cache.Options{
		Scheme: scheme,
	})
	if e != nil {
		return nil, e
	}

	rc := &RuleCache{
		resCache:              resCache,
		ruleSelector:          labels.Everything(),
		ruleNamespaceSelector: labels.Everything(),
	}
	if rcfg.RuleSelector != nil {
		rc.ruleSelector = rcfg.RuleSelector
	}
	if rcfg.RuleNamespaceSelector != nil {
		rc.ruleNamespaceSelector = rcfg.RuleNamespaceSelector
	}

	return rc, nil
}

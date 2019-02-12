// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/isurulucky/k8sday/03-advanced-kubernetes/pkg/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Hellos returns a HelloInformer.
	Hellos() HelloInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Hellos returns a HelloInformer.
func (v *version) Hellos() HelloInformer {
	return &helloInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

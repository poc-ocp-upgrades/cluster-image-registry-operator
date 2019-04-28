package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	version		= "v1"
	groupName	= "imageregistry.operator.openshift.io"
)

var (
	scheme			= runtime.NewScheme()
	SchemeBuilder		= runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme		= SchemeBuilder.AddToScheme
	SchemeGroupVersion	= schema.GroupVersion{Group: groupName, Version: version}
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	AddToScheme(scheme)
}
func addKnownTypes(scheme *runtime.Scheme) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	scheme.AddKnownTypes(SchemeGroupVersion, &Config{}, &ConfigList{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
func Resource(resource string) schema.GroupResource {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

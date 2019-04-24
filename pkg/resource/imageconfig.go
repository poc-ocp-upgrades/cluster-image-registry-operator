package resource

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	kcorelisters "k8s.io/client-go/listers/core/v1"
	configapi "github.com/openshift/api/config/v1"
	configset "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	configlisters "github.com/openshift/client-go/config/listers/config/v1"
	routelisters "github.com/openshift/client-go/route/listers/route/v1"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	"github.com/openshift/cluster-image-registry-operator/pkg/parameters"
)

var _ Mutator = &generatorImageConfig{}

type generatorImageConfig struct {
	configLister	configlisters.ImageLister
	routeLister	routelisters.RouteNamespaceLister
	serviceLister	kcorelisters.ServiceNamespaceLister
	configClient	configset.ConfigV1Interface
	name		string
	namespace	string
	serviceName	string
}

func newGeneratorImageConfig(configLister configlisters.ImageLister, routeLister routelisters.RouteNamespaceLister, serviceLister kcorelisters.ServiceNamespaceLister, configClient configset.ConfigV1Interface, params *parameters.Globals) *generatorImageConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &generatorImageConfig{configLister: configLister, routeLister: routeLister, serviceLister: serviceLister, configClient: configClient, name: params.ImageConfig.Name, namespace: params.Deployment.Namespace, serviceName: params.Service.Name}
}
func (gic *generatorImageConfig) Type() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &configapi.Image{}
}
func (gic *generatorImageConfig) GetGroup() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return configapi.GroupName
}
func (gic *generatorImageConfig) GetResource() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "images"
}
func (gic *generatorImageConfig) GetNamespace() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ""
}
func (gic *generatorImageConfig) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gic.name
}
func (gic *generatorImageConfig) Get() (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gic.configLister.Get(gic.GetName())
}
func (gic *generatorImageConfig) objectMeta() metav1.ObjectMeta {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return metav1.ObjectMeta{Name: gic.GetName()}
}
func (gic *generatorImageConfig) Create() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ic := &configapi.Image{ObjectMeta: gic.objectMeta()}
	ic, err := gic.configClient.Images().Create(ic)
	if err != nil {
		return err
	}
	externalHostnames, err := gic.getRouteHostnames()
	if err != nil {
		return err
	}
	ic.Status.ExternalRegistryHostnames = externalHostnames
	internalHostname, err := getServiceHostname(gic.serviceLister, gic.serviceName)
	if err != nil {
		return err
	}
	ic.Status.InternalRegistryHostname = internalHostname
	_, err = gic.configClient.Images().UpdateStatus(ic)
	return err
}
func (gic *generatorImageConfig) Update(o runtime.Object) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ic := o.(*configapi.Image)
	externalHostnames, err := gic.getRouteHostnames()
	if err != nil {
		return false, err
	}
	modified := false
	if !reflect.DeepEqual(externalHostnames, ic.Status.ExternalRegistryHostnames) {
		ic.Status.ExternalRegistryHostnames = externalHostnames
		modified = true
	}
	internalHostname, err := getServiceHostname(gic.serviceLister, gic.serviceName)
	if err != nil {
		return false, err
	}
	if ic.Status.InternalRegistryHostname != internalHostname {
		ic.Status.InternalRegistryHostname = internalHostname
		modified = true
	}
	if !modified {
		return false, nil
	}
	_, err = gic.configClient.Images().UpdateStatus(ic)
	return err == nil, err
}
func (gic *generatorImageConfig) Delete(opts *metav1.DeleteOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gic.configClient.Images().Delete(gic.GetName(), opts)
}
func (g *generatorImageConfig) Owned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (gic *generatorImageConfig) getRouteHostnames() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var externalHostnames []string
	routes, err := gic.routeLister.List(labels.Everything())
	if err != nil {
		return []string{}, err
	}
	defaultHost := ""
	for _, route := range routes {
		if !RouteIsCreatedByOperator(route) {
			continue
		}
		for _, ingress := range route.Status.Ingress {
			hostname := ingress.Host
			if len(hostname) == 0 {
				continue
			}
			if strings.HasPrefix(hostname, imageregistryv1.DefaultRouteName+"-"+gic.namespace) {
				defaultHost = hostname
			} else {
				externalHostnames = append(externalHostnames, hostname)
			}
		}
	}
	sort.Strings(externalHostnames)
	if len(defaultHost) > 0 {
		externalHostnames = append([]string{defaultHost}, externalHostnames...)
	}
	return externalHostnames, nil
}
func getServiceHostname(serviceLister kcorelisters.ServiceNamespaceLister, serviceName string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	svc, err := serviceLister.Get(serviceName)
	if errors.IsNotFound(err) {
		return "", nil
	}
	if svc == nil || err != nil {
		return "", err
	}
	svcHostname := fmt.Sprintf("%s.%s.svc:%d", svc.Name, svc.Namespace, svc.Spec.Ports[0].Port)
	return svcHostname, nil
}

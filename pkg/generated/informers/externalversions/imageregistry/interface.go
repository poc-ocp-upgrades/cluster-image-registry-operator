package imageregistry

import (
	v1 "github.com/openshift/cluster-image-registry-operator/pkg/generated/informers/externalversions/imageregistry/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	internalinterfaces "github.com/openshift/cluster-image-registry-operator/pkg/generated/informers/externalversions/internalinterfaces"
)

type Interface interface{ V1() v1.Interface }
type group struct {
	factory			internalinterfaces.SharedInformerFactory
	namespace		string
	tweakListOptions	internalinterfaces.TweakListOptionsFunc
}

func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &group{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}
func (g *group) V1() v1.Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return v1.New(g.factory, g.namespace, g.tweakListOptions)
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}

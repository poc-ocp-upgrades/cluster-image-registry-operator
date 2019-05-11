package resource

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	coreset "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	"github.com/openshift/cluster-image-registry-operator/pkg/parameters"
	"github.com/openshift/cluster-image-registry-operator/pkg/resource/strategy"
)

var _ Mutator = &generatorService{}

type generatorService struct {
	lister		corelisters.ServiceNamespaceLister
	client		coreset.CoreV1Interface
	name		string
	namespace	string
	labels		map[string]string
	port		int
	secretName	string
}

func newGeneratorService(lister corelisters.ServiceNamespaceLister, client coreset.CoreV1Interface, params *parameters.Globals, cr *imageregistryv1.Config) *generatorService {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &generatorService{lister: lister, client: client, name: params.Service.Name, namespace: params.Deployment.Namespace, labels: params.Deployment.Labels, port: params.Container.Port, secretName: imageregistryv1.ImageRegistryName + "-tls"}
}
func (gs *generatorService) Type() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &corev1.Service{}
}
func (gs *generatorService) GetGroup() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return corev1.GroupName
}
func (gs *generatorService) GetResource() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "services"
}
func (gs *generatorService) GetNamespace() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gs.namespace
}
func (gs *generatorService) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gs.name
}
func (gs *generatorService) expected() *corev1.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	svc := &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: gs.GetName(), Namespace: gs.GetNamespace(), Labels: gs.labels}, Spec: corev1.ServiceSpec{Selector: gs.labels, Ports: []corev1.ServicePort{{Name: fmt.Sprintf("%d-tcp", gs.port), Port: int32(gs.port), Protocol: "TCP", TargetPort: intstr.FromInt(gs.port)}}}}
	svc.ObjectMeta.Annotations = map[string]string{"service.alpha.openshift.io/serving-cert-secret-name": gs.secretName}
	return svc
}
func (gs *generatorService) Get() (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gs.lister.Get(gs.GetName())
}
func (gs *generatorService) Create() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	svc := &corev1.Service{}
	n := gs.expected()
	_, err := strategy.Service(svc, n)
	if err != nil {
		return err
	}
	_, err = gs.client.Services(gs.GetNamespace()).Create(svc)
	return err
}
func (gs *generatorService) Update(o runtime.Object) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	svc := o.(*corev1.Service)
	n := gs.expected()
	updated, err := strategy.Service(svc, n)
	if !updated || err != nil {
		return false, err
	}
	_, err = gs.client.Services(gs.GetNamespace()).Update(svc)
	return true, err
}
func (gs *generatorService) Delete(opts *metav1.DeleteOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gs.client.Services(gs.GetNamespace()).Delete(gs.GetName(), opts)
}
func (g *generatorService) Owned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}

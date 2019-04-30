package resource

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	coreset "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	"github.com/openshift/cluster-image-registry-operator/pkg/parameters"
)

var _ Mutator = &generatorServiceAccount{}

type generatorServiceAccount struct {
	lister		corelisters.ServiceAccountNamespaceLister
	client		coreset.CoreV1Interface
	name		string
	namespace	string
}

func newGeneratorServiceAccount(lister corelisters.ServiceAccountNamespaceLister, client coreset.CoreV1Interface, params *parameters.Globals, cr *imageregistryv1.Config) *generatorServiceAccount {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &generatorServiceAccount{lister: lister, client: client, name: params.Pod.ServiceAccount, namespace: params.Deployment.Namespace}
}
func (gsa *generatorServiceAccount) Type() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &corev1.ServiceAccount{}
}
func (gsa *generatorServiceAccount) GetGroup() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return corev1.GroupName
}
func (gsa *generatorServiceAccount) GetResource() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "serviceaccounts"
}
func (gsa *generatorServiceAccount) GetNamespace() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gsa.namespace
}
func (gsa *generatorServiceAccount) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gsa.name
}
func (gsa *generatorServiceAccount) expected() (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sa := &corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: gsa.GetName(), Namespace: gsa.GetNamespace()}}
	return sa, nil
}
func (gsa *generatorServiceAccount) Get() (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gsa.lister.Get(gsa.GetName())
}
func (gsa *generatorServiceAccount) Create() (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return commonCreate(gsa, func(obj runtime.Object) (runtime.Object, error) {
		return gsa.client.ServiceAccounts(gsa.GetNamespace()).Create(obj.(*corev1.ServiceAccount))
	})
}
func (gsa *generatorServiceAccount) Update(o runtime.Object) (runtime.Object, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return commonUpdate(gsa, o, func(obj runtime.Object) (runtime.Object, error) {
		return gsa.client.ServiceAccounts(gsa.GetNamespace()).Update(obj.(*corev1.ServiceAccount))
	})
}
func (gsa *generatorServiceAccount) Delete(opts *metav1.DeleteOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gsa.client.ServiceAccounts(gsa.GetNamespace()).Delete(gsa.GetName(), opts)
}
func (g *generatorServiceAccount) Owned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true
}

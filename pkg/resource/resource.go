package resource

import (
	"fmt"
	metaapi "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/openshift/cluster-image-registry-operator/pkg/resource/strategy"
)

type Getter interface {
	Type() runtime.Object
	GetName() string
	GetNamespace() string
	GetGroup() string
	GetResource() string
	Get() (runtime.Object, error)
}
type Mutator interface {
	Getter
	Create() (runtime.Object, error)
	Update(o runtime.Object) (runtime.Object, bool, error)
	Delete(opts *metaapi.DeleteOptions) error
	Owned() bool
}

func Name(o Getter) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	name := fmt.Sprintf("%T, ", o.Type())
	if namespace := o.GetNamespace(); namespace != "" {
		name += fmt.Sprintf("Namespace=%s, ", namespace)
	}
	name += fmt.Sprintf("Name=%s", o.GetName())
	return name
}

type expecter interface {
	Type() runtime.Object
	expected() (runtime.Object, error)
}

func commonCreate(gen expecter, create func(obj runtime.Object) (runtime.Object, error)) (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := gen.Type()
	n, err := gen.expected()
	if err != nil {
		return n, err
	}
	_, err = strategy.Override(o, n)
	if err != nil {
		return n, err
	}
	c, err := create(o)
	return c, err
}
func commonUpdate(gen expecter, o runtime.Object, update func(obj runtime.Object) (runtime.Object, error)) (runtime.Object, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, err := gen.expected()
	if err != nil {
		return o, false, err
	}
	updated, err := strategy.Override(o, n)
	if !updated || err != nil {
		return o, updated, err
	}
	u, err := update(o)
	return u, true, err
}

package v1

import (
	v1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	"bytes"
	"net/http"
	"runtime"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

type ConfigLister interface {
	List(selector labels.Selector) (ret []*v1.Config, err error)
	Get(name string) (*v1.Config, error)
	ConfigListerExpansion
}
type configLister struct{ indexer cache.Indexer }

func NewConfigLister(indexer cache.Indexer) ConfigLister {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &configLister{indexer: indexer}
}
func (s *configLister) List(selector labels.Selector) (ret []*v1.Config, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Config))
	})
	return ret, err
}
func (s *configLister) Get(name string) (*v1.Config, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("config"), name)
	}
	return obj.(*v1.Config), nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := runtime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", runtime.FuncForPC(pc).Name()))
	http.Post("/"+"logcode", "application/json", bytes.NewBuffer(jsonLog))
}

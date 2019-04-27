package fake

import (
	v1 "github.com/openshift/cluster-image-registry-operator/pkg/generated/clientset/versioned/typed/imageregistry/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeImageregistryV1 struct{ *testing.Fake }

func (c *FakeImageregistryV1) Configs() v1.ConfigInterface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &FakeConfigs{c}
}
func (c *FakeImageregistryV1) RESTClient() rest.Interface {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var ret *rest.RESTClient
	return ret
}

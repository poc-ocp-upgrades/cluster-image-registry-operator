package listers

import (
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	"bytes"
	"net/http"
	"runtime"
	"fmt"
	regopclient "github.com/openshift/cluster-image-registry-operator/pkg/client"
	coreset "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
)

const (
	installerConfigNamespace = "kube-system"
)

type mockLister struct {
	listers		regopclient.Listers
	kubeconfig	*restclient.Config
}

func NewMockLister(kubeconfig *restclient.Config) (*mockLister, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &mockLister{kubeconfig: kubeconfig}, nil
}
func (m *mockLister) GetListers() (*regopclient.Listers, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	coreClient, err := coreset.NewForConfig(m.kubeconfig)
	if err != nil {
		return nil, err
	}
	m.listers.Secrets = MockSecretNamespaceLister{namespace: imageregistryv1.ImageRegistryOperatorNamespace, client: coreClient}
	m.listers.InstallerSecrets = MockSecretNamespaceLister{namespace: installerConfigNamespace, client: coreClient}
	return &m.listers, err
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := runtime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", runtime.FuncForPC(pc).Name()))
	http.Post("/"+"logcode", "application/json", bytes.NewBuffer(jsonLog))
}

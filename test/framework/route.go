package framework

import (
	"strings"
	"testing"
	"time"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	routeapiv1 "github.com/openshift/api/route/v1"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
)

func MustEnsureDefaultExternalRouteExists(t *testing.T, client *Clientset) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	var routes *routeapiv1.RouteList
	err = wait.Poll(1*time.Second, AsyncOperationTimeout, func() (bool, error) {
		routes, err = client.Routes(imageregistryv1.ImageRegistryOperatorNamespace).List(metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		if routes == nil || len(routes.Items) < 1 {
			t.Logf("insuffient routes found: %#v", routes)
			return false, nil
		}
		for _, r := range routes.Items {
			if strings.HasPrefix(r.Spec.Host, imageregistryv1.DefaultRouteName+"-"+imageregistryv1.ImageRegistryOperatorNamespace) {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		t.Fatalf("did not find default external route: %#v, err: %v", routes, err)
	}
}
func EnsureExternalRoutesExist(t *testing.T, client *Clientset, wantedRoutes []string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	var routes *routeapiv1.RouteList
	err = wait.Poll(1*time.Second, AsyncOperationTimeout, func() (bool, error) {
		routes, err = client.Routes(imageregistryv1.ImageRegistryOperatorNamespace).List(metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		if routes == nil || len(routes.Items) < len(wantedRoutes)+1 {
			t.Logf("insuffient routes found: %#v", routes)
			return false, nil
		}
		for _, wr := range wantedRoutes {
			found := false
			for _, r := range routes.Items {
				if wr == r.Spec.Host {
					found = true
					break
				}
			}
			if !found {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		t.Errorf("did not find expected routes: wanted: %#v, got: %#v, err: %v", wantedRoutes, routes, err)
	}
}

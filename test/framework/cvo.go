package framework

import (
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	configv1 "github.com/openshift/api/config/v1"
)

const (
	ClusterVersionName = "version"
)

func addCompomentOverride(overrides []configv1.ComponentOverride, override configv1.ComponentOverride) ([]configv1.ComponentOverride, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i, o := range overrides {
		if o.Group == override.Group && o.Kind == override.Kind && o.Namespace == override.Namespace && o.Name == override.Name {
			if overrides[i].Unmanaged == override.Unmanaged {
				return overrides, false
			}
			overrides[i].Unmanaged = override.Unmanaged
			return overrides, true
		}
	}
	return append(overrides, override), true
}
func DisableCVOForOperator(logger Logger, client *Clientset) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cv, err := client.ClusterVersions().Get(ClusterVersionName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil
	} else if err != nil {
		return err
	}
	var changed bool
	cv.Spec.Overrides, changed = addCompomentOverride(cv.Spec.Overrides, configv1.ComponentOverride{Group: "", Kind: "Deployment", Namespace: OperatorDeploymentNamespace, Name: OperatorDeploymentName, Unmanaged: true})
	cv.Spec.Overrides, changed = addCompomentOverride(cv.Spec.Overrides, configv1.ComponentOverride{Group: "", Kind: "Deployment", Namespace: "openshift-kube-apiserver-operator", Name: "kube-apiserver-operator", Unmanaged: true})
	cv.Spec.Overrides, changed = addCompomentOverride(cv.Spec.Overrides, configv1.ComponentOverride{Group: "", Kind: "Deployment", Namespace: "openshift-apiserver-operator", Name: "openshift-apiserver-operator", Unmanaged: true})
	if changed {
		if _, err := client.ClusterVersions().Update(cv); err != nil {
			return err
		}
	}
	if err := StopDeployment(logger, client, "kube-apiserver-operator", "openshift-kube-apiserver-operator"); err != nil {
		return fmt.Errorf("unable to stop kube apiserver operator: %v", err)
	}
	if err := StopDeployment(logger, client, "openshift-apiserver-operator", "openshift-apiserver-operator"); err != nil {
		return fmt.Errorf("unable to stop openshift apiserver operator: %v", err)
	}
	return nil
}

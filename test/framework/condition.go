package framework

import (
	"fmt"
	"time"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	operatorapi "github.com/openshift/api/operator/v1"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
)

func ConditionExistsWithStatusAndReason(client *Clientset, conditionType string, conditionStatus operatorapi.ConditionStatus, conditionReason string) []error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var errs []error
	err := wait.Poll(1*time.Second, AsyncOperationTimeout, func() (stop bool, err error) {
		errs = nil
		conditionExists := false
		cr, err := client.Configs().Get(imageregistryv1.ImageRegistryResourceName, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				return false, nil
			}
			return true, err
		}
		for _, condition := range cr.Status.Conditions {
			if condition.Type == conditionType {
				conditionExists = true
				if condition.Status != conditionStatus {
					errs = append(errs, fmt.Errorf("condition %s status should be \"%v\" but was %v instead", conditionType, conditionStatus, condition.Status))
				}
				if condition.Reason != conditionReason {
					errs = append(errs, fmt.Errorf("condition %s reason should have been \"%s\" but was %s instead", conditionType, conditionReason, condition.Reason))
				}
			}
		}
		if !conditionExists {
			errs = append(errs, fmt.Errorf("condition %s was not found, but should have been. %#v", conditionType, cr.Status.Conditions))
		}
		if len(errs) != 0 {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}

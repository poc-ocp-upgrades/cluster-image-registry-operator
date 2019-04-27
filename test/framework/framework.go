package framework

import (
	"fmt"
	"testing"
	"time"
	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	AsyncOperationTimeout = 5 * time.Minute
)

type Logger interface{ Logf(string, ...interface{}) }

var _ Logger = &testing.T{}

func DumpObject(logger Logger, prefix string, obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	logger.Logf("%s:\n%s", prefix, spew.Sdump(obj))
}
func DumpYAML(logger Logger, prefix string, obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, err := yaml.Marshal(obj)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal object: %s", err))
	}
	logger.Logf("%s:\n%s", prefix, string(data))
}
func DeleteCompletely(getObject func() (metav1.Object, error), deleteObject func(*metav1.DeleteOptions) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := getObject()
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	uid := obj.GetUID()
	policy := metav1.DeletePropagationForeground
	if err := deleteObject(&metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}, PropagationPolicy: &policy}); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return wait.Poll(1*time.Second, AsyncOperationTimeout, func() (stop bool, err error) {
		obj, err = getObject()
		if err != nil {
			if errors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		return obj.GetUID() != uid, nil
	})
}

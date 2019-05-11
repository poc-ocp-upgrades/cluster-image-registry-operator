package operator

import (
	"fmt"
	"time"
	"github.com/golang/glog"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	regopset "github.com/openshift/cluster-image-registry-operator/pkg/generated/clientset/versioned/typed/imageregistry/v1"
	"github.com/openshift/cluster-image-registry-operator/pkg/parameters"
	"github.com/openshift/cluster-image-registry-operator/pkg/util"
)

func (c *Controller) RemoveResources(o *imageregistryv1.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.setStatusRemoving(o)
	return c.generator.Remove(o)
}
func (c *Controller) finalizeResources(o *imageregistryv1.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if o.ObjectMeta.DeletionTimestamp == nil {
		return nil
	}
	finalizers := []string{}
	for _, v := range o.ObjectMeta.Finalizers {
		if v != parameters.ImageRegistryOperatorResourceFinalizer {
			finalizers = append(finalizers, v)
		}
	}
	if len(finalizers) == len(o.ObjectMeta.Finalizers) {
		return nil
	}
	glog.Infof("finalizing %s", util.ObjectInfo(o))
	client, err := regopset.NewForConfig(c.kubeconfig)
	if err != nil {
		return err
	}
	err = c.RemoveResources(o)
	if err != nil {
		c.setStatusRemoveFailed(o, err)
		return fmt.Errorf("unable to finalize resource: %s", err)
	}
	cr := o
	err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if cr == nil {
			cr, err := client.Configs().Get(o.Name, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("failed to get %s: %s", util.ObjectInfo(o), err)
			}
			finalizers = []string{}
			for _, v := range cr.ObjectMeta.Finalizers {
				if v != parameters.ImageRegistryOperatorResourceFinalizer {
					finalizers = append(finalizers, v)
				}
			}
		}
		cr.ObjectMeta.Finalizers = finalizers
		_, err := client.Configs().Update(cr)
		if err != nil {
			cr = nil
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to update finalizers in %s: %s", util.ObjectInfo(o), err)
	}
	errorFuncs := []func(error) bool{kerrors.IsInternalError, kerrors.IsTimeout, kerrors.IsServerTimeout, kerrors.IsTooManyRequests, utilnet.IsProbableEOF, utilnet.IsConnectionReset}
	retryTime := 3 * time.Second
	err = wait.PollInfinite(retryTime, func() (stop bool, err error) {
		_, err = c.listers.RegistryConfigs.Get(o.Name)
		if err == nil {
			return
		}
		if !kerrors.IsNotFound(err) {
			for _, isRetryError := range errorFuncs {
				if isRetryError(err) {
					return false, nil
				}
			}
			if delaySeconds, shouldRetry := kerrors.SuggestsClientDelay(err); shouldRetry {
				delayTime := time.Duration(delaySeconds) * time.Second
				if retryTime < delayTime {
					time.Sleep(delayTime - retryTime)
				}
				return false, nil
			}
			err = fmt.Errorf("failed to get %s: %s", util.ObjectInfo(o), err)
			return
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("unable to wait for %s deletion: %s", util.ObjectInfo(o), err)
	}
	return nil
}

package resource

import (
	"fmt"
	"reflect"
	"github.com/golang/glog"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	configapi "github.com/openshift/api/config/v1"
	operatorapi "github.com/openshift/api/operator/v1"
	configset "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	configlisters "github.com/openshift/client-go/config/listers/config/v1"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	appsapi "k8s.io/api/apps/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
)

var _ Mutator = &generatorClusterOperator{}

type generatorClusterOperator struct {
	mutators	[]Mutator
	cr		*imageregistryv1.Config
	deployLister	appslisters.DeploymentNamespaceLister
	configLister	configlisters.ClusterOperatorLister
	configClient	configset.ConfigV1Interface
}

func newGeneratorClusterOperator(deployLister appslisters.DeploymentNamespaceLister, configLister configlisters.ClusterOperatorLister, configClient configset.ConfigV1Interface, cr *imageregistryv1.Config, mutators []Mutator) *generatorClusterOperator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &generatorClusterOperator{mutators: mutators, cr: cr, deployLister: deployLister, configLister: configLister, configClient: configClient}
}
func (gco *generatorClusterOperator) Type() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &configapi.ClusterOperator{}
}
func (gco *generatorClusterOperator) GetGroup() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return configapi.GroupName
}
func (gco *generatorClusterOperator) GetResource() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "clusteroperators"
}
func (gco *generatorClusterOperator) GetNamespace() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ""
}
func (gco *generatorClusterOperator) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return imageregistryv1.ImageRegistryClusterOperatorResourceName
}
func (gco *generatorClusterOperator) Get() (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gco.configLister.Get(gco.GetName())
}
func (gco *generatorClusterOperator) Create() (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	co := &configapi.ClusterOperator{ObjectMeta: metav1.ObjectMeta{Name: gco.GetName()}}
	_, err := gco.syncVersions(co)
	if err != nil {
		return co, err
	}
	_ = gco.syncConditions(co)
	_ = gco.syncRelatedObjects(co)
	return gco.configClient.ClusterOperators().Create(co)
}
func (gco *generatorClusterOperator) Update(o runtime.Object) (runtime.Object, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	co := o.(*configapi.ClusterOperator)
	modified, err := gco.syncVersions(co)
	if err != nil {
		return o, false, err
	}
	if gco.syncConditions(co) {
		modified = true
	}
	if gco.syncRelatedObjects(co) {
		modified = true
	}
	if !modified {
		return o, false, nil
	}
	n, err := gco.configClient.ClusterOperators().UpdateStatus(co)
	return n, err == nil, err
}
func (gco *generatorClusterOperator) Delete(opts *metav1.DeleteOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return gco.configClient.Images().Delete(gco.GetName(), opts)
}
func (gco *generatorClusterOperator) Owned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func convertOperatorStatus(status operatorapi.ConditionStatus) (configapi.ConditionStatus, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch status {
	case operatorapi.ConditionTrue:
		return configapi.ConditionTrue, nil
	case operatorapi.ConditionFalse:
		return configapi.ConditionFalse, nil
	case operatorapi.ConditionUnknown:
		return configapi.ConditionUnknown, nil
	}
	return configapi.ConditionUnknown, fmt.Errorf("unexpected condition status: %s", status)
}
func (gco *generatorClusterOperator) syncConditions(op *configapi.ClusterOperator) (modified bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	conditions := []configapi.ClusterOperatorStatusCondition{}
	for _, resourceCondition := range gco.cr.Status.Conditions {
		found := false
		var conditionType configapi.ClusterStatusConditionType
		switch resourceCondition.Type {
		case operatorapi.OperatorStatusTypeAvailable:
			conditionType = configapi.OperatorAvailable
		case operatorapi.OperatorStatusTypeProgressing:
			conditionType = configapi.OperatorProgressing
		case operatorapi.OperatorStatusTypeDegraded:
			conditionType = configapi.OperatorDegraded
		default:
			continue
		}
		for i, clusterOperatorCondition := range op.Status.Conditions {
			if conditionType != clusterOperatorCondition.Type {
				continue
			}
			found = true
			newStatus, err := convertOperatorStatus(resourceCondition.Status)
			if err != nil {
				glog.Errorf("ignore condition of %s custom resource: %s", gco.cr.Name, err)
				continue
			}
			if clusterOperatorCondition.Status == newStatus {
				continue
			}
			op.Status.Conditions[i].Status = newStatus
			op.Status.Conditions[i].LastTransitionTime = resourceCondition.LastTransitionTime
			op.Status.Conditions[i].Reason = resourceCondition.Reason
			op.Status.Conditions[i].Message = resourceCondition.Message
			modified = true
		}
		if !found {
			conditionStatus, err := convertOperatorStatus(resourceCondition.Status)
			if err != nil {
				glog.Errorf("ignore condition of %s custom resource: %s", gco.cr.Name, err)
				continue
			}
			conditions = append(conditions, configapi.ClusterOperatorStatusCondition{Type: conditionType, Status: conditionStatus, LastTransitionTime: resourceCondition.LastTransitionTime, Reason: resourceCondition.Reason, Message: resourceCondition.Message})
			modified = true
		}
	}
	for i := range conditions {
		op.Status.Conditions = append(op.Status.Conditions, conditions[i])
	}
	return
}
func isDeploymentStatusAvailableAndUpdated(deploy *appsapi.Deployment) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return deploy.Status.AvailableReplicas > 0 && deploy.Status.ObservedGeneration >= deploy.Generation && deploy.Status.UpdatedReplicas == deploy.Status.Replicas
}
func (gco *generatorClusterOperator) syncVersions(op *configapi.ClusterOperator) (modified bool, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	deploy, err := gco.deployLister.Get(imageregistryv1.ImageRegistryName)
	if err != nil {
		if kerrors.IsNotFound(err) {
			err = nil
		}
		return
	}
	deploymentVersion := deploy.Annotations[imageregistryv1.VersionAnnotation]
	if len(deploymentVersion) == 0 || !isDeploymentStatusAvailableAndUpdated(deploy) {
		return
	}
	newVersions := []configapi.OperandVersion{{Name: "operator", Version: deploymentVersion}}
	if !reflect.DeepEqual(op.Status.Versions, newVersions) {
		op.Status.Versions = newVersions
		modified = true
	}
	return
}
func (gco *generatorClusterOperator) syncRelatedObjects(op *configapi.ClusterOperator) (modified bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var relatedObjects []configapi.ObjectReference
	for _, gen := range gco.mutators {
		relatedObjects = append(relatedObjects, configapi.ObjectReference{Group: gen.GetGroup(), Resource: gen.GetResource(), Namespace: gen.GetNamespace(), Name: gen.GetName()})
	}
	if !reflect.DeepEqual(op.Status.RelatedObjects, relatedObjects) {
		op.Status.RelatedObjects = relatedObjects
		modified = true
	}
	return
}

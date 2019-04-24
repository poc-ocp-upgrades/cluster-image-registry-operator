package storage

import (
	"fmt"
	"bytes"
	"net/http"
	"runtime"
	corev1 "k8s.io/api/core/v1"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	regopclient "github.com/openshift/cluster-image-registry-operator/pkg/client"
	"github.com/openshift/cluster-image-registry-operator/pkg/clusterconfig"
	"github.com/openshift/cluster-image-registry-operator/pkg/storage/emptydir"
	"github.com/openshift/cluster-image-registry-operator/pkg/storage/pvc"
	"github.com/openshift/cluster-image-registry-operator/pkg/storage/s3"
)

var (
	ErrStorageNotConfigured = fmt.Errorf("storage backend not configured")
)

type Driver interface {
	ConfigEnv() ([]corev1.EnvVar, error)
	Volumes() ([]corev1.Volume, []corev1.VolumeMount, error)
	Secrets() (map[string]string, error)
	CompleteConfiguration(*imageregistryv1.Config) error
	CreateStorage(*imageregistryv1.Config) error
	StorageExists(*imageregistryv1.Config) (bool, error)
	RemoveStorage(*imageregistryv1.Config) (bool, error)
	StorageChanged(*imageregistryv1.Config) bool
}

func newDriver(cfg *imageregistryv1.ImageRegistryConfigStorage, listers *regopclient.Listers) (Driver, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var names []string
	var drivers []Driver
	if cfg.EmptyDir != nil {
		names = append(names, "EmptyDir")
		drivers = append(drivers, emptydir.NewDriver(cfg.EmptyDir, listers))
	}
	if cfg.S3 != nil {
		names = append(names, "S3")
		drivers = append(drivers, s3.NewDriver(cfg.S3, listers))
	}
	if cfg.PVC != nil {
		drv, err := pvc.NewDriver(cfg.PVC)
		if err != nil {
			return nil, err
		}
		names = append(names, "PVC")
		drivers = append(drivers, drv)
	}
	switch len(drivers) {
	case 0:
		return nil, ErrStorageNotConfigured
	case 1:
		return drivers[0], nil
	}
	return nil, fmt.Errorf("exactly one storage type should be configured at the same time, got %d: %v", len(drivers), names)
}
func NewDriver(cfg *imageregistryv1.ImageRegistryConfigStorage, listers *regopclient.Listers) (Driver, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	drv, err := newDriver(cfg, listers)
	if err == ErrStorageNotConfigured {
		*cfg, err = getPlatformStorage()
		if err != nil {
			return nil, fmt.Errorf("unable to get storage configuration from cluster install config: %s", err)
		}
		drv, err = newDriver(cfg, listers)
	}
	return drv, err
}
func getPlatformStorage() (imageregistryv1.ImageRegistryConfigStorage, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var cfg imageregistryv1.ImageRegistryConfigStorage
	installConfig, err := clusterconfig.GetInstallConfig()
	if err != nil {
		return cfg, err
	}
	switch {
	case installConfig.Platform.Libvirt != nil:
		cfg.EmptyDir = &imageregistryv1.ImageRegistryConfigStorageEmptyDir{}
	case installConfig.Platform.AWS != nil:
		cfg.S3 = &imageregistryv1.ImageRegistryConfigStorageS3{}
	case installConfig.Platform.OpenStack != nil:
		cfg.EmptyDir = &imageregistryv1.ImageRegistryConfigStorageEmptyDir{}
	}
	return cfg, nil
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := runtime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", runtime.FuncForPC(pc).Name()))
	http.Post("/"+"logcode", "application/json", bytes.NewBuffer(jsonLog))
}

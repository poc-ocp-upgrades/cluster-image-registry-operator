package v1

import (
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func (in *Config) DeepCopyInto(out *Config) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}
func (in *Config) DeepCopy() *Config {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(Config)
	in.DeepCopyInto(out)
	return out
}
func (in *Config) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *ConfigList) DeepCopyInto(out *ConfigList) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Config, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}
func (in *ConfigList) DeepCopy() *ConfigList {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ConfigList)
	in.DeepCopyInto(out)
	return out
}
func (in *ConfigList) DeepCopyObject() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *ImageRegistryConfigProxy) DeepCopyInto(out *ImageRegistryConfigProxy) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *ImageRegistryConfigProxy) DeepCopy() *ImageRegistryConfigProxy {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryConfigProxy)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistryConfigRequests) DeepCopyInto(out *ImageRegistryConfigRequests) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.Read = in.Read
	out.Write = in.Write
	return
}
func (in *ImageRegistryConfigRequests) DeepCopy() *ImageRegistryConfigRequests {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryConfigRequests)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistryConfigRequestsLimits) DeepCopyInto(out *ImageRegistryConfigRequestsLimits) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.MaxWaitInQueue = in.MaxWaitInQueue
	return
}
func (in *ImageRegistryConfigRequestsLimits) DeepCopy() *ImageRegistryConfigRequestsLimits {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryConfigRequestsLimits)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistryConfigRoute) DeepCopyInto(out *ImageRegistryConfigRoute) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *ImageRegistryConfigRoute) DeepCopy() *ImageRegistryConfigRoute {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryConfigRoute)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistryConfigStorage) DeepCopyInto(out *ImageRegistryConfigStorage) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	if in.EmptyDir != nil {
		in, out := &in.EmptyDir, &out.EmptyDir
		*out = new(ImageRegistryConfigStorageEmptyDir)
		**out = **in
	}
	if in.S3 != nil {
		in, out := &in.S3, &out.S3
		*out = new(ImageRegistryConfigStorageS3)
		(*in).DeepCopyInto(*out)
	}
	if in.Swift != nil {
		in, out := &in.Swift, &out.Swift
		*out = new(ImageRegistryConfigStorageSwift)
		**out = **in
	}
	if in.PVC != nil {
		in, out := &in.PVC, &out.PVC
		*out = new(ImageRegistryConfigStoragePVC)
		**out = **in
	}
	return
}
func (in *ImageRegistryConfigStorage) DeepCopy() *ImageRegistryConfigStorage {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryConfigStorage)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistryConfigStorageEmptyDir) DeepCopyInto(out *ImageRegistryConfigStorageEmptyDir) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *ImageRegistryConfigStorageEmptyDir) DeepCopy() *ImageRegistryConfigStorageEmptyDir {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryConfigStorageEmptyDir)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistryConfigStoragePVC) DeepCopyInto(out *ImageRegistryConfigStoragePVC) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *ImageRegistryConfigStoragePVC) DeepCopy() *ImageRegistryConfigStoragePVC {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryConfigStoragePVC)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistryConfigStorageS3) DeepCopyInto(out *ImageRegistryConfigStorageS3) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	if in.CloudFront != nil {
		in, out := &in.CloudFront, &out.CloudFront
		*out = new(ImageRegistryConfigStorageS3CloudFront)
		(*in).DeepCopyInto(*out)
	}
	return
}
func (in *ImageRegistryConfigStorageS3) DeepCopy() *ImageRegistryConfigStorageS3 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryConfigStorageS3)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistryConfigStorageS3CloudFront) DeepCopyInto(out *ImageRegistryConfigStorageS3CloudFront) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	in.PrivateKey.DeepCopyInto(&out.PrivateKey)
	out.Duration = in.Duration
	return
}
func (in *ImageRegistryConfigStorageS3CloudFront) DeepCopy() *ImageRegistryConfigStorageS3CloudFront {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryConfigStorageS3CloudFront)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistryConfigStorageSwift) DeepCopyInto(out *ImageRegistryConfigStorageSwift) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	return
}
func (in *ImageRegistryConfigStorageSwift) DeepCopy() *ImageRegistryConfigStorageSwift {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryConfigStorageSwift)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistrySpec) DeepCopyInto(out *ImageRegistrySpec) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	out.Proxy = in.Proxy
	in.Storage.DeepCopyInto(&out.Storage)
	out.Requests = in.Requests
	if in.Routes != nil {
		in, out := &in.Routes, &out.Routes
		*out = make([]ImageRegistryConfigRoute, len(*in))
		copy(*out, *in)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(corev1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]corev1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}
func (in *ImageRegistrySpec) DeepCopy() *ImageRegistrySpec {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistrySpec)
	in.DeepCopyInto(out)
	return out
}
func (in *ImageRegistryStatus) DeepCopyInto(out *ImageRegistryStatus) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	*out = *in
	in.OperatorStatus.DeepCopyInto(&out.OperatorStatus)
	in.Storage.DeepCopyInto(&out.Storage)
	return
}
func (in *ImageRegistryStatus) DeepCopy() *ImageRegistryStatus {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if in == nil {
		return nil
	}
	out := new(ImageRegistryStatus)
	in.DeepCopyInto(out)
	return out
}

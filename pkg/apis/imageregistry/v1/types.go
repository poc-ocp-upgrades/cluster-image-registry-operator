package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	operatorsv1api "github.com/openshift/api/operator/v1"
)

const (
	DefaultRouteName				= "default-route"
	ImageRegistryName				= "image-registry"
	ImageRegistryResourceName			= "cluster"
	ImageRegistryCertificatesName			= ImageRegistryName + "-certificates"
	ImageRegistryPrivateConfiguration		= ImageRegistryName + "-private-configuration"
	ImageRegistryPrivateConfigurationUser		= ImageRegistryPrivateConfiguration + "-user"
	ImageRegistryOperatorNamespace			= "openshift-image-registry"
	ImageRegistryClusterOperatorResourceName	= "image-registry"
	OperatorStatusTypeRemoved			= "Removed"
	StorageExists					= "StorageExists"
	StorageTagged					= "StorageTagged"
	StorageEncrypted				= "StorageEncrypted"
	StorageIncompleteUploadCleanupEnabled		= "StorageIncompleteUploadCleanupEnabled"
	VersionAnnotation				= "release.openshift.io/version"
)

type ConfigList struct {
	metav1.TypeMeta	`json:",inline"`
	metav1.ListMeta	`json:"metadata"`
	Items		[]Config	`json:"items"`
}
type Config struct {
	metav1.TypeMeta		`json:",inline"`
	metav1.ObjectMeta	`json:"metadata"`
	Spec			ImageRegistrySpec	`json:"spec"`
	Status			ImageRegistryStatus	`json:"status"`
}
type ImageRegistrySpec struct {
	ManagementState	operatorsv1api.ManagementState	`json:"managementState"`
	HTTPSecret	string				`json:"httpSecret"`
	Proxy		ImageRegistryConfigProxy	`json:"proxy"`
	Storage		ImageRegistryConfigStorage	`json:"storage"`
	ReadOnly	bool				`json:"readOnly"`
	Requests	ImageRegistryConfigRequests	`json:"requests"`
	DefaultRoute	bool				`json:"defaultRoute"`
	Routes		[]ImageRegistryConfigRoute	`json:"routes,omitempty"`
	Replicas	int32				`json:"replicas"`
	LogLevel	int64				`json:"logging"`
	Resources	*corev1.ResourceRequirements	`json:"resources,omitempty"`
	NodeSelector	map[string]string		`json:"nodeSelector,omitempty"`
	Tolerations	[]corev1.Toleration		`json:"tolerations,omitempty"`
}
type ImageRegistryStatus struct {
	operatorsv1api.OperatorStatus	`json:",inline"`
	StorageManaged			bool				`json:"storageManaged"`
	Storage				ImageRegistryConfigStorage	`json:"storage"`
}
type ImageRegistryConfigProxy struct {
	HTTP	string	`json:"http"`
	HTTPS	string	`json:"https"`
	NoProxy	string	`json:"noProxy"`
}
type ImageRegistryConfigStorageS3CloudFront struct {
	BaseURL		string				`json:"baseURL"`
	PrivateKey	corev1.SecretKeySelector	`json:"privateKey"`
	KeypairID	string				`json:"keypairID"`
	Duration	metav1.Duration			`json:"duration"`
}
type ImageRegistryConfigStorageEmptyDir struct{}
type ImageRegistryConfigStorageS3 struct {
	Bucket		string					`json:"bucket"`
	Region		string					`json:"region"`
	RegionEndpoint	string					`json:"regionEndpoint"`
	Encrypt		bool					`json:"encrypt"`
	KeyID		string					`json:"keyID"`
	CloudFront	*ImageRegistryConfigStorageS3CloudFront	`json:"cloudFront,omitempty"`
}
type ImageRegistryConfigStorageSwift struct {
	AuthURL		string	`json:"authURL"`
	Container	string	`json:"container"`
	Domain		string	`json:"domain"`
	DomainID	string	`json:"domainID"`
	Tenant		string	`json:"tenant"`
	TenantID	string	`json:"tenantID"`
}
type ImageRegistryConfigStoragePVC struct {
	Claim string `json:"claim"`
}
type ImageRegistryConfigStorage struct {
	EmptyDir	*ImageRegistryConfigStorageEmptyDir	`json:"emptyDir,omitempty"`
	S3		*ImageRegistryConfigStorageS3		`json:"s3,omitempty"`
	Swift		*ImageRegistryConfigStorageSwift	`json:"swift,omitempty"`
	PVC		*ImageRegistryConfigStoragePVC		`json:"pvc,omitempty"`
}
type ImageRegistryConfigRequests struct {
	Read	ImageRegistryConfigRequestsLimits	`json:"read"`
	Write	ImageRegistryConfigRequestsLimits	`json:"write"`
}
type ImageRegistryConfigRequestsLimits struct {
	MaxRunning	int		`json:"maxRunning"`
	MaxInQueue	int		`json:"maxInQueue"`
	MaxWaitInQueue	metav1.Duration	`json:"maxWaitInQueue"`
}
type ImageRegistryConfigRoute struct {
	Name		string	`json:"name"`
	Hostname	string	`json:"hostname"`
	SecretName	string	`json:"secretName"`
}

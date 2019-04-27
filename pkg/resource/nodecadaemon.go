package resource

import (
	"os"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	appsclientv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	kcorelisters "k8s.io/client-go/listers/core/v1"
	"github.com/openshift/library-go/pkg/operator/resource/resourceread"
	"github.com/openshift/cluster-image-registry-operator/pkg/parameters"
)

const (
	nodeCADaemonSetDefinition = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: node-ca
  namespace: openshift-image-registry
spec:
  selector:
    matchLabels:
      name: node-ca
  template:
    metadata:
      labels:
        name: node-ca
    spec:      
      nodeSelector:
        beta.kubernetes.io/os: linux
      priorityClassName: system-cluster-critical
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
        operator: Exists
      serviceAccountName: node-ca
      containers:
      - name: node-ca
        securityContext:
          privileged: true
        image: docker.io/openshift/origin-cluster-image-registry-operator:latest
        command: 
        - "/bin/sh"
        - "-c"
        - |
          if [ ! -e /etc/docker/certs.d/image-registry.openshift-image-registry.svc.cluster.local:5000 ]; then
            mkdir /etc/docker/certs.d/image-registry.openshift-image-registry.svc.cluster.local:5000
          fi
          if [ ! -e /etc/docker/certs.d/${internalRegistryHostname} ]; then
            mkdir /etc/docker/certs.d/${internalRegistryHostname}
          fi
          while [ true ];
          do
            for f in $(ls /tmp/serviceca); do
                if [ "${f}" == "service-ca.crt" ]; then
                    continue
                fi
                echo $f
                ca_file_path="/tmp/serviceca/${f}"
                f=$(echo $f | sed  -r 's/(.*)\.\./\1:/')
                reg_dir_path="/etc/docker/certs.d/${f}"
                if [ -e "${reg_dir_path}" ]; then
                    cp -u $ca_file_path $reg_dir_path/ca.crt
                else
                    mkdir $reg_dir_path
                    cp $ca_file_path $reg_dir_path/ca.crt
                fi
            done
            for d in $(ls /etc/docker/certs.d); do
                echo $d
                if [ "${d}" == "${internalRegistryHostname}" ]; then
                    continue
                fi
                if [ "${d}" == "image-registry.openshift-image-registry.svc.cluster.local:5000" ]; then
                    continue
                fi
                dp=$(echo $d | sed  -r 's/(.*):/\1\.\./')
                reg_conf_path="/tmp/serviceca/${dp}"
                if [ ! -e "${reg_conf_path}" ]; then
                    rm -rf /etc/docker/certs.d/$d
                fi
            done
            if [ -e /tmp/serviceca/service-ca.crt ]; then
              cp -u /tmp/serviceca/service-ca.crt /etc/docker/certs.d/image-registry.openshift-image-registry.svc.cluster.local:5000
              cp -u /tmp/serviceca/service-ca.crt /etc/docker/certs.d/${internalRegistryHostname}
            else 
              rm /etc/docker/certs.d/image-registry.openshift-image-registry.svc.cluster.local:5000/service-ca.crt
              rm /etc/docker/certs.d/${internalRegistryHostname}/service-ca.crt
            fi
            sleep 60
          done
        volumeMounts:
        - name: serviceca
          mountPath: /tmp/serviceca
        - name: host
          mountPath: /etc/docker/certs.d
      volumes:
      - name: host
        hostPath:
          path: /etc/docker/certs.d
      - name: serviceca
        configMap:
          name: image-registry-certificates
`
)

var _ Mutator = &generatorNodeCADaemonSet{}

type generatorNodeCADaemonSet struct {
	daemonSetLister	appslisters.DaemonSetNamespaceLister
	serviceLister	kcorelisters.ServiceNamespaceLister
	client		appsclientv1.AppsV1Interface
	params		*parameters.Globals
}

func newGeneratorNodeCADaemonSet(daemonSetLister appslisters.DaemonSetNamespaceLister, serviceLister kcorelisters.ServiceNamespaceLister, client appsclientv1.AppsV1Interface, params *parameters.Globals) *generatorNodeCADaemonSet {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &generatorNodeCADaemonSet{daemonSetLister: daemonSetLister, serviceLister: serviceLister, client: client, params: params}
}
func (ds *generatorNodeCADaemonSet) Type() runtime.Object {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &appsv1.DaemonSet{}
}
func (ds *generatorNodeCADaemonSet) GetGroup() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return appsv1.GroupName
}
func (ds *generatorNodeCADaemonSet) GetResource() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "daemonsets"
}
func (ds *generatorNodeCADaemonSet) GetNamespace() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ds.params.Deployment.Namespace
}
func (ds *generatorNodeCADaemonSet) GetName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return "node-ca"
}
func (ds *generatorNodeCADaemonSet) Get() (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ds.daemonSetLister.Get(ds.GetName())
}
func (ds *generatorNodeCADaemonSet) Create() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	internalHostname, err := getServiceHostname(ds.serviceLister, ds.params.Service.Name)
	if err != nil {
		return err
	}
	daemonSet := resourceread.ReadDaemonSetV1OrDie([]byte(nodeCADaemonSetDefinition))
	env := corev1.EnvVar{Name: "internalRegistryHostname", Value: internalHostname}
	daemonSet.Spec.Template.Spec.Containers[0].Image = os.Getenv("IMAGE")
	daemonSet.Spec.Template.Spec.Containers[0].Env = append(daemonSet.Spec.Template.Spec.Containers[0].Env, env)
	_, err = ds.client.DaemonSets(ds.GetNamespace()).Create(daemonSet)
	return err
}
func (ds *generatorNodeCADaemonSet) Update(o runtime.Object) (bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	internalHostname, err := getServiceHostname(ds.serviceLister, ds.params.Service.Name)
	if err != nil {
		return false, err
	}
	daemonSet := o.(*appsv1.DaemonSet)
	modified := false
	exists := false
	newImage := os.Getenv("IMAGE")
	oldImage := daemonSet.Spec.Template.Spec.Containers[0].Image
	if newImage != oldImage {
		daemonSet.Spec.Template.Spec.Containers[0].Image = newImage
		modified = true
	}
	for i, env := range daemonSet.Spec.Template.Spec.Containers[0].Env {
		if env.Name == "internalRegistryHostname" {
			exists = true
			if env.Value != internalHostname {
				daemonSet.Spec.Template.Spec.Containers[0].Env[i].Value = internalHostname
				modified = true
			}
			break
		}
	}
	if !exists {
		env := corev1.EnvVar{Name: "internalRegistryHostname", Value: internalHostname}
		daemonSet.Spec.Template.Spec.Containers[0].Env = append(daemonSet.Spec.Template.Spec.Containers[0].Env, env)
		modified = true
	}
	if !modified {
		return false, nil
	}
	_, err = ds.client.DaemonSets(ds.GetNamespace()).Update(daemonSet)
	return err == nil, err
}
func (ds *generatorNodeCADaemonSet) Delete(opts *metav1.DeleteOptions) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return ds.client.DaemonSets(ds.GetNamespace()).Delete(ds.GetName(), opts)
}
func (ds *generatorNodeCADaemonSet) Owned() bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return false
}

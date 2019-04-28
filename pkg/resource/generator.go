package resource

import (
	"fmt"
	"time"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/wait"
	routeset "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaapi "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	appsset "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreset "k8s.io/client-go/kubernetes/typed/core/v1"
	rbacset "k8s.io/client-go/kubernetes/typed/rbac/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	configset "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	imageregistryv1 "github.com/openshift/cluster-image-registry-operator/pkg/apis/imageregistry/v1"
	"github.com/openshift/cluster-image-registry-operator/pkg/client"
	"github.com/openshift/cluster-image-registry-operator/pkg/parameters"
	"github.com/openshift/cluster-image-registry-operator/pkg/resource/object"
	"github.com/openshift/cluster-image-registry-operator/pkg/storage"
)

func NewGenerator(kubeconfig *rest.Config, listers *client.Listers, params *parameters.Globals) *Generator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &Generator{kubeconfig: kubeconfig, listers: listers, params: params}
}

type Generator struct {
	kubeconfig	*rest.Config
	listers		*client.Listers
	params		*parameters.Globals
}

func (g *Generator) listRoutes(routeClient routeset.RouteV1Interface, cr *imageregistryv1.Config) []Mutator {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var mutators []Mutator
	if cr.Spec.DefaultRoute {
		mutators = append(mutators, newGeneratorRoute(g.listers.Routes, g.listers.Secrets, routeClient, g.params, cr, imageregistryv1.ImageRegistryConfigRoute{Name: imageregistryv1.DefaultRouteName}))
	}
	for _, route := range cr.Spec.Routes {
		mutators = append(mutators, newGeneratorRoute(g.listers.Routes, g.listers.Secrets, routeClient, g.params, cr, route))
	}
	return mutators
}
func (g *Generator) list(cr *imageregistryv1.Config) ([]Mutator, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	coreClient, err := coreset.NewForConfig(g.kubeconfig)
	if err != nil {
		return nil, err
	}
	appsClient, err := appsset.NewForConfig(g.kubeconfig)
	if err != nil {
		return nil, err
	}
	rbacClient, err := rbacset.NewForConfig(g.kubeconfig)
	if err != nil {
		return nil, err
	}
	routeClient, err := routeset.NewForConfig(g.kubeconfig)
	if err != nil {
		return nil, err
	}
	configClient, err := configset.NewForConfig(g.kubeconfig)
	if err != nil {
		return nil, err
	}
	driver, err := storage.NewDriver(&cr.Spec.Storage, g.listers)
	if err != nil {
		return nil, err
	}
	var mutators []Mutator
	mutators = append(mutators, newGeneratorClusterRole(g.listers.ClusterRoles, rbacClient, cr))
	mutators = append(mutators, newGeneratorClusterRoleBinding(g.listers.ClusterRoleBindings, rbacClient, g.params, cr))
	mutators = append(mutators, newGeneratorServiceAccount(g.listers.ServiceAccounts, coreClient, g.params, cr))
	mutators = append(mutators, newGeneratorCAConfig(g.listers.ConfigMaps, g.listers.ImageConfigs, g.listers.OpenShiftConfig, coreClient, g.params, cr))
	mutators = append(mutators, newGeneratorSecret(g.listers.Secrets, coreClient, driver, g.params, cr))
	mutators = append(mutators, newGeneratorImageConfig(g.listers.ImageConfigs, g.listers.Routes, g.listers.Services, configClient, g.params))
	mutators = append(mutators, newGeneratorNodeCADaemonSet(g.listers.DaemonSets, g.listers.Services, appsClient, g.params))
	mutators = append(mutators, newGeneratorService(g.listers.Services, coreClient, g.params, cr))
	mutators = append(mutators, newGeneratorDeployment(g.listers.Deployments, g.listers.ConfigMaps, g.listers.Secrets, coreClient, appsClient, driver, g.params, cr))
	mutators = append(mutators, g.listRoutes(routeClient, cr)...)
	mutators = append(mutators, newGeneratorClusterOperator(g.listers.Deployments, g.listers.ClusterOperators, configClient, cr, mutators))
	return mutators, nil
}
func (g *Generator) syncStorage(cr *imageregistryv1.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var runCreate bool
	driver, err := storage.NewDriver(&cr.Spec.Storage, g.listers)
	if err != nil {
		return err
	}
	if driver.StorageChanged(cr) {
		runCreate = true
	} else {
		exists, err := driver.StorageExists(cr)
		if err != nil {
			return err
		}
		if !exists {
			runCreate = true
		}
	}
	if runCreate {
		if err := driver.CreateStorage(cr); err != nil {
			return err
		}
	}
	return nil
}
func (g *Generator) removeObsoleteRoutes(cr *imageregistryv1.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	routeClient, err := routeset.NewForConfig(g.kubeconfig)
	if err != nil {
		return err
	}
	routes, err := g.listers.Routes.List(labels.Everything())
	if err != nil {
		return fmt.Errorf("failed to list routes: %s", err)
	}
	routesGenerators := g.listRoutes(routeClient, cr)
	knownNames := map[string]struct{}{}
	for _, gen := range routesGenerators {
		knownNames[gen.GetName()] = struct{}{}
	}
	gracePeriod := int64(0)
	propagationPolicy := metaapi.DeletePropagationForeground
	opts := &metaapi.DeleteOptions{GracePeriodSeconds: &gracePeriod, PropagationPolicy: &propagationPolicy}
	for _, route := range routes {
		if !RouteIsCreatedByOperator(route) {
			continue
		}
		if _, found := knownNames[route.Name]; found {
			continue
		}
		err = routeClient.Routes(g.params.Deployment.Namespace).Delete(route.Name, opts)
		if err != nil {
			return err
		}
	}
	return nil
}
func (g *Generator) Apply(cr *imageregistryv1.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := g.syncStorage(cr)
	if err == storage.ErrStorageNotConfigured {
		return err
	} else if err != nil {
		return fmt.Errorf("unable to sync storage configuration: %s", err)
	}
	generators, err := g.list(cr)
	if err != nil {
		return fmt.Errorf("unable to get generators: %s", err)
	}
	for _, gen := range generators {
		err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			o, err := gen.Get()
			if err != nil {
				if !errors.IsNotFound(err) {
					return fmt.Errorf("failed to get object %s: %s", Name(gen), err)
				}
				n, err := gen.Create()
				if err != nil {
					return fmt.Errorf("failed to create object %s: %s", Name(gen), err)
				}
				str, err := object.DumpString(n)
				if err != nil {
					glog.Errorf("unable to dump object: %s", err)
				}
				glog.Infof("object %s created: %s", Name(gen), str)
				return nil
			}
			n, updated, err := gen.Update(o.DeepCopyObject())
			if err != nil {
				if errors.IsConflict(err) {
					return err
				}
				return fmt.Errorf("failed to update object %s: %s", Name(gen), err)
			}
			if updated {
				difference, err := object.DiffString(o, n)
				if err != nil {
					glog.Errorf("unable to calculate difference: %s", err)
				}
				glog.Infof("object %s updated: %s", Name(gen), difference)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("unable to apply objects: %s", err)
		}
	}
	err = g.removeObsoleteRoutes(cr)
	if err != nil {
		return fmt.Errorf("unable to remove obsolete routes: %s", err)
	}
	return nil
}
func (g *Generator) Remove(cr *imageregistryv1.Config) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	generators, err := g.list(cr)
	if err != nil {
		return fmt.Errorf("unable to get generators: %s", err)
	}
	gracePeriod := int64(0)
	propagationPolicy := metaapi.DeletePropagationForeground
	opts := &metaapi.DeleteOptions{GracePeriodSeconds: &gracePeriod, PropagationPolicy: &propagationPolicy}
	for _, gen := range generators {
		if !gen.Owned() {
			continue
		}
		if err := gen.Delete(opts); err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			return fmt.Errorf("failed to delete object %s: %s", Name(gen), err)
		}
		glog.Infof("object %s deleted", Name(gen))
	}
	driver, err := storage.NewDriver(&cr.Status.Storage, g.listers)
	if err != nil {
		return err
	}
	var derr error
	var retriable bool
	err = wait.PollImmediate(1*time.Second, 5*time.Minute, func() (stop bool, err error) {
		if retriable, derr = driver.RemoveStorage(cr); derr != nil {
			if retriable {
				return false, nil
			} else {
				return true, derr
			}
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("unable to remove storage: %s, %s", err, derr)
	}
	return nil
}

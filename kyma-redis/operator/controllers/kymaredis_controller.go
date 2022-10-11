/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	operatorv1alpha1 "github.com/adityabhatia/kcd-demo/kyma-redis/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kyma-project/module-manager/operator/pkg/declarative"
	"github.com/kyma-project/module-manager/operator/pkg/types"
)

// KymaRedisReconciler reconciles a KymaRedis object
type KymaRedisReconciler struct {
	declarative.ManifestReconciler
	client.Client
	Scheme *runtime.Scheme
	*rest.Config
}

type ManifestResolver struct {
}

const (
	sampleAnnotationKey   = "owner/kyma-project.io"
	sampleAnnotationValue = "kyma-kcd-demo"
	chartPath             = "./module-chart"
	chartNs               = "kyma-redis"
)

// Get returns the chart information to be processed.
func (m ManifestResolver) Get(object types.BaseCustomObject) (types.InstallationSpec, error) {
	sample, valid := object.(*operatorv1alpha1.KymaRedis)
	if !valid {
		return types.InstallationSpec{},
			fmt.Errorf("invalid type conversion for %s", client.ObjectKeyFromObject(object))
	}
	return types.InstallationSpec{
		ChartPath:   chartPath,
		ReleaseName: sample.Spec.ReleaseName,
		ChartFlags: types.ChartFlags{
			ConfigFlags: types.Flags{
				"Namespace":       chartNs,
				"CreateNamespace": true,
			},
			SetFlags: types.Flags{},
		},
	}, nil
}

// initReconciler injects the required configuration into the declarative reconciler.
func (r *KymaRedisReconciler) initReconciler(mgr ctrl.Manager) error {
	manifestResolver := &ManifestResolver{}
	return r.Inject(mgr, &operatorv1alpha1.KymaRedis{},
		declarative.WithManifestResolver(manifestResolver),
		declarative.WithCustomResourceLabels(map[string]string{"sampleKey": "sampleValue"}),
		declarative.WithPostRenderTransform(transform),
		declarative.WithResourcesReady(true),
	)
}

// transform modifies the resources based on some criteria, before installation.
func transform(_ context.Context, _ types.BaseCustomObject, resources *types.ManifestResources) error {
	for _, resource := range resources.Items {
		annotations := resource.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string, 0)
		}
		if annotations[sampleAnnotationKey] == "" {
			annotations[sampleAnnotationKey] = sampleAnnotationValue
			resource.SetAnnotations(annotations)
		}
	}
	return nil
}

//+kubebuilder:rbac:groups=operator.kyma-project.io,resources=kymaredis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.kyma-project.io,resources=kymaredis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.kyma-project.io,resources=kymaredis/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch;get;list;watch
//+kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete

// TODO: dynamically create RBACs! Remove line below.
//+kubebuilder:rbac:groups="*",resources="*",verbs="*"

// SetupWithManager sets up the controller with the Manager.
func (r *KymaRedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Config = mgr.GetConfig()
	if err := r.initReconciler(mgr); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.KymaRedis{}).
		Complete(r)
}

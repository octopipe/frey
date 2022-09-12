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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/argoproj/gitops-engine/pkg/sync"
	stackv1alpha1 "github.com/octopipe/frey/api/v1alpha1"
	"github.com/octopipe/frey/templates/app"
)

// ResourceReconciler reconciles a Resource object
type ResourceReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	GitOpsEngine engine.GitOpsEngine
}

//+kubebuilder:rbac:groups=stack.octopipe.io,resources=resources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stack.octopipe.io,resources=resources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stack.octopipe.io,resources=resources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Resource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info(fmt.Sprintf("SYNC RESOURCE %s START...", req.Name))
	genericApp := app.GenericApp{
		Name:         req.Name,
		Image:        "nginx:1.14.2",
		Port:         8080,
		ResourceName: req.Name,
	}
	n, err := app.NewGenericApp(genericApp)
	if err != nil {
		logger.Error(err, "unable to create new generic app")
		return ctrl.Result{}, err
	}

	_, err = r.GitOpsEngine.Sync(context.Background(), n, func(r *cache.Resource) bool {
		return true
	}, time.Now().String(), "default", sync.WithPrune(true), sync.WithLogr(logger))

	if err != nil {
		logger.Error(err, "failed to sync resource")
		return ctrl.Result{}, err

	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stackv1alpha1.Resource{}).
		Complete(r)
}

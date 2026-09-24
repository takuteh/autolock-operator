/*
Copyright 2026.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"
)

// AutolockReconciler reconciles a Autolock object
type AutolockReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=autolock.takuteh.uk,resources=autolocks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autolock.takuteh.uk,resources=autolocks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=autolock.takuteh.uk,resources=autolocks/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;create;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Autolock object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *AutolockReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	log.Info("Reconcile start",
		"name", req.Name,
		"namespace", req.Namespace,
	)
	var autolock autolockv1alpha1.Autolock

	//CRを取得しautolock変数に格納
	if err := r.Get(ctx, req.NamespacedName, &autolock); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.reconcileConfigMap(ctx, &autolock); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.reconcileMainDeployment(ctx, &autolock); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutolockReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autolockv1alpha1.Autolock{}).
		Named("autolock").
		Complete(r)
}

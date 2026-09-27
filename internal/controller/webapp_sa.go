package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const webappServiceAccountName = "autolock-webapp"

// Webapp用ServiceAccount
func (r *AutolockReconciler) reconcileWebappServiceAccount(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      webappServiceAccountName,
			Namespace: autolock.Namespace,
		},
	}

	if err := ctrl.SetControllerReference(autolock, serviceAccount, r.Scheme); err != nil {
		return err
	}

	var existing corev1.ServiceAccount
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      serviceAccount.Name,
			Namespace: serviceAccount.Namespace,
		},
		&existing,
	)

	// 存在しなければ作成
	if err != nil {
		if apierrors.IsNotFound(err) {
			return r.Create(ctx, serviceAccount)
		}
		return err
	}

	return nil
}

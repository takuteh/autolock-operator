package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const webappServiceAccountName = "autolock-webapp"

// Webapp用ServiceAccount
func (r *AutolockReconciler) reconcileWebappServiceAccount(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	serviceAccount := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      webappServiceAccountName,
			Namespace: autolock.Namespace,
		},
	}

	if err := controllerutil.SetControllerReference(
		autolock,
		serviceAccount,
		r.Scheme,
	); err != nil {
		return err
	}

	return r.Patch(
		ctx,
		serviceAccount,
		client.Apply,
		client.FieldOwner("autolock-operator"),
	)
}

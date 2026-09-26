package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *AutolockReconciler) reconcileAuthDBService(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-auth-db",
			Namespace: autolock.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "autolock-auth-db",
			},
			Ports: []corev1.ServicePort{
				{
					Port:       3306,
					TargetPort: intstr.FromInt32(3306),
				},
			},
		},
	}

	var existing corev1.Service

	err := r.Get(ctx, types.NamespacedName{
		Name:      "autolock-auth-db",
		Namespace: autolock.Namespace,
	}, &existing)

	if err := ctrl.SetControllerReference(autolock, service, r.Scheme); err != nil {
		return err
	}

	if err != nil {
		if apierrors.IsNotFound(err) {
			return r.Create(ctx, service)
		}
		return err
	}

	existing.Spec.Selector = service.Spec.Selector
	existing.Spec.Ports = service.Spec.Ports

	return r.Update(ctx, &existing)
}

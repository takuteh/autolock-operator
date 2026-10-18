package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *AutolockReconciler) reconcileWebappService(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-webapp",
			Namespace: autolock.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "autolock-webapp",
			},
			Ports: []corev1.ServicePort{
				{
					Port:       3000,
					TargetPort: intstr.FromInt32(3000),
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(
		autolock,
		service,
		r.Scheme,
	); err != nil {
		return err
	}

	return r.Patch(
		ctx,
		service,
		client.Apply,
		client.FieldOwner("autolock-operator"),
	)
}

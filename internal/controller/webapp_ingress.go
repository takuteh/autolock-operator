package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *AutolockReconciler) reconcileWebappIngress(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	cfg := autolock.Spec.Webapp.Ingress

	// Ingress が無効なら終了
	if !cfg.Enabled {
		return nil
	}

	pathType := networkingv1.PathType(cfg.PathType)

	ingress := &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-webapp",
			Namespace: autolock.Namespace,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &cfg.ClassName,
			Rules: []networkingv1.IngressRule{
				{
					Host: cfg.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     cfg.Path,
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: "autolock-webapp",
											Port: networkingv1.ServiceBackendPort{
												Number: 3000,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// OwnerReference を設定
	if err := controllerutil.SetControllerReference(
		autolock,
		ingress,
		r.Scheme,
	); err != nil {
		return err
	}

	// Server-Side Apply
	return r.Patch(
		ctx,
		ingress,
		client.Apply,
		client.FieldOwner("autolock-operator"),
	)
}

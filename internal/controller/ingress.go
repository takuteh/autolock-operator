package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *AutolockReconciler) reconcileIngress(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	cfg := autolock.Spec.Webapp.Ingress

	// Ingressの定義を作成
	pathType := networkingv1.PathType(cfg.PathType)

	ingress := &networkingv1.Ingress{
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

	// Ingressが無効なら終了
	if !cfg.Enabled {
		return nil
	}

	// クラスタからIngressを取得
	var existing networkingv1.Ingress

	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      "autolock-webapp",
			Namespace: autolock.Namespace,
		},
		&existing,
	)

	// OwnerReferenceを設定
	if err := ctrl.SetControllerReference(
		autolock,
		ingress,
		r.Scheme,
	); err != nil {
		return err
	}

	// 存在しなければ作成
	if err != nil {
		if apierrors.IsNotFound(err) {
			if err = r.Create(ctx, ingress); err != nil {
				return err
			}

			return nil
		}

		return err
	}

	// 存在する場合、内容に変更があれば更新
	if !reflect.DeepEqual(existing.Spec, ingress.Spec) {
		existing.Spec = ingress.Spec
		if err := r.Update(ctx, &existing); err != nil {
			return err
		}
	}

	return nil
}

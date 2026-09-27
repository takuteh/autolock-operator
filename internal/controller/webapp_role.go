package controller

import (
	"context"
	"reflect"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *AutolockReconciler) reconcileWebappRole(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-cr-editor",
			Namespace: autolock.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"autolock.takuteh.uk",
				},
				Resources: []string{
					"autolocks",
				},
				Verbs: []string{
					"get",
					"list",
					"patch",
					"update",
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(autolock, role, r.Scheme); err != nil {
		return err
	}

	var existing rbacv1.Role
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      role.Name,
			Namespace: role.Namespace,
		},
		&existing,
	)

	// 存在しなければ作成
	if err != nil {
		if apierrors.IsNotFound(err) {
			return r.Create(ctx, role)
		}
		return err
	}

	// 内容に変更があれば更新
	if !reflect.DeepEqual(existing.Rules, role.Rules) {
		existing.Rules = role.Rules
		if err := r.Update(ctx, &existing); err != nil {
			return err
		}
	}

	return nil
}

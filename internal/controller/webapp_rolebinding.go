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

func (r *AutolockReconciler) reconcileWebappRoleBinding(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-webapp-to-cr-editor",
			Namespace: autolock.Namespace,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     "autolock-cr-editor",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      webappServiceAccountName,
				Namespace: autolock.Namespace,
			},
		},
	}

	if err := ctrl.SetControllerReference(autolock, roleBinding, r.Scheme); err != nil {
		return err
	}

	var existing rbacv1.RoleBinding
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      roleBinding.Name,
			Namespace: roleBinding.Namespace,
		},
		&existing,
	)

	// 存在しなければ作成
	if err != nil {
		if apierrors.IsNotFound(err) {
			return r.Create(ctx, roleBinding)
		}
		return err
	}

	// RoleRefは作成後に変更できないため、Subjectsのみ比較・更新
	if !reflect.DeepEqual(existing.Subjects, roleBinding.Subjects) {
		existing.Subjects = roleBinding.Subjects
		if err := r.Update(ctx, &existing); err != nil {
			return err
		}
	}

	return nil
}

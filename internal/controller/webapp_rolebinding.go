package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *AutolockReconciler) reconcileWebappRoleBinding(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	roleBinding := &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "RoleBinding",
		},
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

	if err := controllerutil.SetControllerReference(
		autolock,
		roleBinding,
		r.Scheme,
	); err != nil {
		return err
	}

	return r.Patch(
		ctx,
		roleBinding,
		client.Apply,
		client.FieldOwner("autolock-operator"),
	)
}

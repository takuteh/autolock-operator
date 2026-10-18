package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *AutolockReconciler) reconcileWebappRole(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	role := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "Role",
		},
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

	if err := controllerutil.SetControllerReference(
		autolock,
		role,
		r.Scheme,
	); err != nil {
		return err
	}

	return r.Patch(
		ctx,
		role,
		client.Apply,
		client.FieldOwner("autolock-operator"),
	)
}

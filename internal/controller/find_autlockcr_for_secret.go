package controller

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// findAutolockForSecretは、Secretのannotationから対応するAutolock CRを特定し、Reconcileリクエストを返す。
func (r *AutolockReconciler) findAutolockCRForSecret(
	ctx context.Context,
	obj client.Object,
) []reconcile.Request {
	name := obj.GetAnnotations()["autolock.takuteh.uk/owner-cr"]
	if name == "" {
		return nil
	}

	return []reconcile.Request{
		{
			NamespacedName: client.ObjectKey{
				Namespace: obj.GetNamespace(),
				Name:      name,
			},
		},
	}
}

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *AutolockReconciler) getSecretHash(
	ctx context.Context,
	namespace string,
	name string,
) (string, error) {
	var secret corev1.Secret

	if err := r.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, &secret); err != nil {
		return "", err
	}

	hash, err := configHash(secret.Data)
	if err != nil {
		return "", err
	}

	return hash, nil
}

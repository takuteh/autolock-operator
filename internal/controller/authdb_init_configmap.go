package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *AutolockReconciler) reconcileAuthDBInitConfigMap(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-auth-db-init",
			Namespace: autolock.Namespace,
		},
		Data: map[string]string{
			"init.sql": `
			CREATE TABLE IF NOT EXISTS users(
				id INT AUTO_INCREMENT PRIMARY KEY,
				user_name VARCHAR(50) NOT NULL,
				line_id VARCHAR(50),
				slack_id VARCHAR(50),
				start_date DATETIME DEFAULT '2000-01-01 00:00:00',
				end_date DATETIME DEFAULT '2031-12-31 23:59:59'
			);`,
		},
	}

	if err := controllerutil.SetControllerReference(
		autolock,
		configMap,
		r.Scheme,
	); err != nil {
		return err
	}

	return r.Patch(
		ctx,
		configMap,
		client.Apply,
		client.FieldOwner("autolock-operator"),
	)
}

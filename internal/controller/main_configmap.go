package controller

import (
	"context"
	"encoding/json"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *AutolockReconciler) reconcileMainConfigMap(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) (string, error) {
	// CR の main 設定を JSON 化
	jsonData, err := json.Marshal(autolock.Spec.Main.Config)
	if err != nil {
		return "", err
	}

	// main 設定用 ConfigMap
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-config",
			Namespace: autolock.Namespace,
		},
		Data: map[string]string{
			"autolock_setting.json": string(jsonData),
		},
	}

	if err := controllerutil.SetControllerReference(
		autolock,
		configMap,
		r.Scheme,
	); err != nil {
		return "", err
	}

	if err := r.Patch(
		ctx,
		configMap,
		client.Apply,
		client.FieldOwner("autolock-operator"),
	); err != nil {
		return "", err
	}

	// Apply 成功後に設定のハッシュを返す
	return configHash(autolock.Spec.Main.Config)
}

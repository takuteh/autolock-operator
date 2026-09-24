package controller

import (
	"context"
	"encoding/json"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *AutolockReconciler) reconcileConfigMap(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	//CRから取得したmain設定をjson化し,jsonDataに格納
	jsonData, err := json.Marshal(autolock.Spec.Main.Config)
	if err != nil {
		return err
	}

	//main設定用cmを作成
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-config",
			Namespace: autolock.Namespace,
		},
		Data: map[string]string{
			"autolock_setting.json": string(jsonData),
		},
	}

	//クラスタからcmを取得
	var existing corev1.ConfigMap
	err = r.Get(
		ctx,
		types.NamespacedName{
			Name:      "autolock-config",
			Namespace: autolock.Namespace,
		},
		&existing,
	)

	//存在しなければ作成
	if err != nil {
		if apierrors.IsNotFound(err) {
			if err = r.Create(ctx, configMap); err != nil {
				return err
			}

			return nil
		}

		return err
	}

	// 存在する → 更新
	existing.Data = configMap.Data

	if err := r.Update(ctx, &existing); err != nil {
		return err
	}

	return nil
}

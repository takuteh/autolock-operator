package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *AutolockReconciler) reconcileAuthDBInitConfigMap(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	//DBの初期化用cmを作成
	configMap := &corev1.ConfigMap{
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

	//クラスタからcmを取得
	var existing corev1.ConfigMap
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      "autolock-auth-db-init",
			Namespace: autolock.Namespace,
		},
		&existing,
	)

	if err := ctrl.SetControllerReference(autolock, configMap, r.Scheme); err != nil {
		return err
	}

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

	// 存在する場合、内容に変更があれば更新
	if !reflect.DeepEqual(existing.Data, configMap.Data) {
		existing.Data = configMap.Data
		if err := r.Update(ctx, &existing); err != nil {
			return err
		}
	}

	return nil
}

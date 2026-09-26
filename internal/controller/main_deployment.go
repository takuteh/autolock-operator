package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *AutolockReconciler) reconcileMainDeployment(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
	configHash string,
) error {
	privileged := true
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-main",
			Namespace: autolock.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &autolock.Spec.Main.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "autolock-main",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "autolock-main",
					},
					Annotations: map[string]string{
						"autolock/config-hash": configHash,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "autolock",
							Image: autolock.Spec.Main.Image,
							SecurityContext: &corev1.SecurityContext{
								Privileged: &privileged,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "autolock-config",
									MountPath: "/home/pi/autolock/gear_version/etc/autolock_setting.json",
									SubPath:   "autolock_setting.json",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "autolock-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "autolock-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	//クラスタからdeploymentを取得
	var existing appsv1.Deployment
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      "autolock-main",
			Namespace: autolock.Namespace,
		},
		&existing,
	)

	if err := ctrl.SetControllerReference(autolock, deployment, r.Scheme); err != nil {
		return err
	}

	//存在しなければ作成
	if err != nil {
		if apierrors.IsNotFound(err) {
			if err = r.Create(ctx, deployment); err != nil {
				return err
			}

			return nil
		}

		return err
	}

	// 存在する場合、内容に変更があれば更新
	if !reflect.DeepEqual(existing.Spec, deployment.Spec) {
		existing.Spec = deployment.Spec
		if err := r.Update(ctx, &existing); err != nil {
			return err
		}
	}

	return nil
}

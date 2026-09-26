package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
)

func int32Ptr(v int32) *int32 {
	return &v
}

func (r *AutolockReconciler) reconcileAuthDBStatefulSet(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-auth-db",
			Namespace: autolock.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "autolock-auth-db-service",
			Replicas:    int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "autolock-auth-db",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "autolock-auth-db",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "auth-db",
							Image: autolock.Spec.AuthDB.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3306,
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: autolock.Spec.AuthDB.SecretName,
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/var/lib/mysql",
								},
								{
									Name:      "initdb",
									MountPath: "/docker-entrypoint-initdb.d",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "initdb",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "autolock-auth-db-init",
									},
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "data",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							corev1.ReadWriteOnce,
						},
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("250Mi"),
							},
						},
					},
				},
			},
		},
	}

	// クラスタからStatefulSetを取得
	var existing appsv1.StatefulSet

	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      "autolock-auth-db",
			Namespace: autolock.Namespace,
		},
		&existing,
	)

	// OwnerReferenceを設定
	if err := ctrl.SetControllerReference(
		autolock,
		statefulSet,
		r.Scheme,
	); err != nil {
		return err
	}

	// 存在しなければ作成
	if err != nil {
		if apierrors.IsNotFound(err) {
			if err = r.Create(ctx, statefulSet); err != nil {
				return err
			}

			return nil
		}

		return err
	}

	// 存在する → 更新
	if !reflect.DeepEqual(existing.Spec, statefulSet.Spec) {
		existing.Spec = statefulSet.Spec

		if err := r.Update(ctx, &existing); err != nil {
			return err
		}
	}

	return nil
}

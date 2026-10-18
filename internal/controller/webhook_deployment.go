package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *AutolockReconciler) reconcileWebhookDeployment(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
	configHash string,
) error {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-webhook",
			Namespace: autolock.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &autolock.Spec.Webhook.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "autolock-webhook",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "autolock-webhook",
					},
					Annotations: map[string]string{
						"autolock/config-hash": configHash,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "webhook",
							Image: autolock.Spec.Webhook.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3000,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "autolock-config",
									MountPath: "/app/etc/autolock_setting.json",
									SubPath:   "autolock_setting.json",
									ReadOnly:  true,
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

	// AutolockをDeploymentの所有者として設定
	if err := controllerutil.SetControllerReference(
		autolock,
		deployment,
		r.Scheme,
	); err != nil {
		return err
	}

	// Server-Side ApplyでDeploymentを適用
	if err := r.Patch(
		ctx,
		deployment,
		client.Apply,
		client.FieldOwner("autolock-operator"),
	); err != nil {
		return err
	}

	return nil
}

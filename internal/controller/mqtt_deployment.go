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

func (r *AutolockReconciler) reconcileMQTTDeployment(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-mqtt-broker",
			Namespace: autolock.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &autolock.Spec.MQTT.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "autolock-mqtt-broker",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "autolock-mqtt-broker",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mqtt-broker",
							Image: autolock.Spec.MQTT.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 1883,
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

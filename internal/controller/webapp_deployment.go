package controller

import (
	"context"

	autolockv1alpha1 "github.com/takuteh/autolock-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *AutolockReconciler) reconcileWebappDeployment(
	ctx context.Context,
	autolock *autolockv1alpha1.Autolock,
) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "autolock-webapp",
			Namespace: autolock.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &autolock.Spec.Webapp.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "autolock-webapp",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "autolock-webapp",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "webapp",
							Image: autolock.Spec.Webapp.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3000,
								},
							},
						},
					},
				},
			},
		},
	}

	//クラスタからdeploymentを取得
	var existingDeployment appsv1.Deployment
	err := r.Get(
		ctx,
		types.NamespacedName{
			Name:      "autolock-webapp",
			Namespace: autolock.Namespace,
		},
		&existingDeployment,
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

	// 存在する → 更新
	existingDeployment.Spec = deployment.Spec

	if err := r.Update(ctx, &existingDeployment); err != nil {
		return err
	}

	return nil
}

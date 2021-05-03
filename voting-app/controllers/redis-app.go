package controllers

import (
	pollv1alpha1 "github.com/andmagom/voting-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *VotingAppReconciler) RedisDeployment(v *pollv1alpha1.VotingApp) *appsv1.Deployment {
	labels := labels(v.Name + "redis-app")
	size := int32(1)


	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name + "-redis",
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "redis:alpine",
						Name:  "redis",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 6379,
							Name:          "redis",
						}},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

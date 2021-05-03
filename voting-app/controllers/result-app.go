package controllers

import (
	pollv1alpha1 "github.com/andmagom/voting-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *VotingAppReconciler) ResultAppDeployment(v *pollv1alpha1.VotingApp) *appsv1.Deployment {
	labels := labels(v.Name + "result-app")
	size := int32(1)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name + "-result",
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
						Image: "dockersamples/examplevotingapp_result:before",
						Name:  "result-webui",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "result",
						}},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

func (r *VotingAppReconciler) DBDeployment(v *pollv1alpha1.VotingApp) *appsv1.Deployment {
	labels := labels(v.Name + "db-app")
	size := int32(1)

	env := []corev1.EnvVar{}
	env = append(env, corev1.EnvVar{
		Name:  "PGDATA",
		Value: "/var/lib/postgresql/data/pgdata",
	})
	env = append(env, corev1.EnvVar{
		Name:  "POSTGRES_USER",
		Value: "postgres",
	})

	env = append(env, corev1.EnvVar{
		Name:  "POSTGRES_PASSWORD",
		Value: "postgres",
	})

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.Name + "-db",
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
						Image: "postgres:9.4",
						Name:  "db",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 5432,
							Name:          "db",
						}},
						Env: env,
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

package controllers

import (
	"context"

	pollv1alpha1 "github.com/andmagom/voting-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func dbAppLabels(v *pollv1alpha1.VotingApp) map[string]string {
	labels := labels(v.Name + "db-app")
	return labels
}

func dbDeploymentName(name string) string {
	nameDeployment := name + "-db"
	return nameDeployment
}

func (r *VotingAppReconciler) DBDeployment(v *pollv1alpha1.VotingApp) *appsv1.Deployment {
	labels := dbAppLabels(v)
	size := int32(1)

	env := []corev1.EnvVar{}
	/*
		env = append(env, corev1.EnvVar{
			Name:  "PGDATA",
			Value: "/var/lib/postgresql/data/pgdata",
		})
	*/
	env = append(env, corev1.EnvVar{
		Name:  "POSTGRES_USER",
		Value: "postgres",
	})

	env = append(env, corev1.EnvVar{
		Name:  "POSTGRES_HOST_AUTH_METHOD",
		Value: "trust",
	})

	env = append(env, corev1.EnvVar{
		Name:  "POSTGRES_PASSWORD",
		Value: "postgres",
	})

	dep := DeploymentScheme(v.Namespace, dbDeploymentName(v.Name), &size, labels, "postgres:9.4", "db", int32(dbPort), env)

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

func (r *VotingAppReconciler) DBService(v *pollv1alpha1.VotingApp) *corev1.Service {
	serviceName := "db"
	selector := dbAppLabels(v)

	svc := ServiceScheme(v.Namespace, serviceName, selector, dbPort, corev1.ServiceTypeClusterIP)

	controllerutil.SetControllerReference(v, svc, r.Scheme)
	return svc
}

func (r *VotingAppReconciler) isDBUp(v *pollv1alpha1.VotingApp) bool {
	deployment := &appsv1.Deployment{}

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      dbDeploymentName(v.Name),
		Namespace: v.Namespace,
	}, deployment)

	if err != nil {
		log.Error(err, "Deployment mysql not found")
		return false
	}

	if deployment.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}

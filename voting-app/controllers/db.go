package controllers

import (
	pollv1alpha1 "github.com/andmagom/voting-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func dbAppLabels(v *pollv1alpha1.VotingApp) map[string]string {
	labels := labels(v.Name + "db-app")
	return labels
}

func (r *VotingAppReconciler) DBDeployment(v *pollv1alpha1.VotingApp) *appsv1.Deployment {
	labels := dbAppLabels(v)
	size := int32(1)
	name := v.Name + "-db"

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

	dep := DeploymentScheme(v.Namespace, name, &size, labels, "postgres:9.4", "db", int32(dbPort), env)

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

func (r *VotingAppReconciler) DBService(v *pollv1alpha1.VotingApp) *corev1.Service {
	serviceName := "svc-db-" + v.Name
	selector := dbAppLabels(v)

	svc := ServiceScheme(v.Namespace, serviceName, selector, dbPort, corev1.ServiceTypeNodePort)

	controllerutil.SetControllerReference(v, svc, r.Scheme)
	return svc
}

package controllers

import (
	pollv1alpha1 "github.com/andmagom/voting-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const redisPort int = 6379

func (r *VotingAppReconciler) RedisDeployment(v *pollv1alpha1.VotingApp) *appsv1.Deployment {
	labels := labels(v.Name + "-redis-app")
	size := int32(1)
	name := v.Name + "-redis"

	dep := DeploymentScheme(v.Namespace, name, &size, labels, "redis:alpine", "redis", int32(redisPort), nil)

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

func (r *VotingAppReconciler) RedisService(v *pollv1alpha1.VotingApp) *corev1.Service {
	serviceName := "redis"
	selector := labels(v.Name + "-redis-app")

	svc := ServiceScheme(v.Namespace, serviceName, selector, redisPort, corev1.ServiceTypeClusterIP)

	controllerutil.SetControllerReference(v, svc, r.Scheme)
	return svc
}

package controllers

import (
	pollv1alpha1 "github.com/andmagom/voting-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func workerName(v *pollv1alpha1.VotingApp) string {
	return v.Name + "-worker"
}

func (r *VotingAppReconciler) workerDeployment(v *pollv1alpha1.VotingApp) *appsv1.Deployment {
	labels := labels("worker")
	size := int32(1)
	dep := DeploymentScheme(v.Namespace, workerName(v), &size, labels, "dockersamples/examplevotingapp_worker", "worker", 8000, nil)
	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

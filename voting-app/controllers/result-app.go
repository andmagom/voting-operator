package controllers

import (
	pollv1alpha1 "github.com/andmagom/voting-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const resultAppPort int = 80
const dbPort int = 5432

func resultAppLabels(v *pollv1alpha1.VotingApp) map[string]string {
	labels := labels(v.Name + "result-app")
	return labels
}

func (r *VotingAppReconciler) ResultAppDeployment(v *pollv1alpha1.VotingApp) *appsv1.Deployment {
	labels := resultAppLabels(v)
	size := int32(1)

	env := []corev1.EnvVar{}
	env = append(env, corev1.EnvVar{
		Name:  "OPTION_A",
		Value: v.Spec.OptionA,
	})
	env = append(env, corev1.EnvVar{
		Name:  "OPTION_B",
		Value: v.Spec.OptionB,
	})

	deploymentName := v.Name + "-result"
	dep := DeploymentScheme(v.Namespace, deploymentName, &size, labels, "andmagom/examplevotingapp_result:1.0.0", "result-webui", int32(resultAppPort), env)

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

func (r *VotingAppReconciler) ResultService(v *pollv1alpha1.VotingApp) *corev1.Service {
	serviceName := "svc-result-" + v.Name
	selector := resultAppLabels(v)

	svc := ServiceScheme(v.Namespace, serviceName, selector, resultAppPort, corev1.ServiceTypeNodePort)

	controllerutil.SetControllerReference(v, svc, r.Scheme)
	return svc
}

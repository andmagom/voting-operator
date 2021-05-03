package controllers

import (
	pollv1alpha1 "github.com/andmagom/voting-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
							ContainerPort: int32(resultAppPort),
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

func (r *VotingAppReconciler) ResultService(v *pollv1alpha1.VotingApp) *corev1.Service {
	serviceName := "svc-result-" + v.Name
	selector := resultAppLabels(v)

	svc := ServiceScheme(v.Namespace, serviceName, selector, resultAppPort, corev1.ServiceTypeNodePort)

	controllerutil.SetControllerReference(v, svc, r.Scheme)
	return svc
}

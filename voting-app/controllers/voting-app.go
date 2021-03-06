package controllers

import (
	pollv1alpha1 "github.com/andmagom/voting-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func votingDeploymentName(v *pollv1alpha1.VotingApp) string {
	return v.Name + "-voting"
}

func (r *VotingAppReconciler) votingAppDeployment(v *pollv1alpha1.VotingApp) *appsv1.Deployment {
	labels := labels(v.Name + "-voting-app")
	size := v.Spec.VotingAppReplicas

	env := []corev1.EnvVar{}
	env = append(env, corev1.EnvVar{
		Name:  "OPTION_A",
		Value: v.Spec.OptionA,
	})
	env = append(env, corev1.EnvVar{
		Name:  "OPTION_B",
		Value: v.Spec.OptionB,
	})

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      votingDeploymentName(v),
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
						Image: "dockersamples/examplevotingapp_vote:before",
						Name:  "visitors-webui",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "visitors",
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

func (r *VotingAppReconciler) ServiceVotingApp(v *pollv1alpha1.VotingApp) *corev1.Service {
	serviceName := "svc-voting-app-" + v.Name
	selector := labels(v.Name + "-voting-app")

	svc := ServiceScheme(v.Namespace, serviceName, selector, 80, corev1.ServiceTypeLoadBalancer)

	controllerutil.SetControllerReference(v, svc, r.Scheme)
	return svc
}

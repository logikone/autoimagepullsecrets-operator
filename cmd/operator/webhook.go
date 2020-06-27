package main

import (
	"context"
	"os"

	adminregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func checkMutatingWebhookConfiguration(clientSet *kubernetes.Clientset, caSecret corev1.Secret) {
	ctx := context.Background()

	var mutatingWebhookConfiguration adminregv1.MutatingWebhookConfiguration
	mutatingWebhookConfiguration.Name = "aips-injector"

	mwh, err := clientSet.
		AdmissionregistrationV1().
		MutatingWebhookConfigurations().
		Get(ctx, mutatingWebhookConfiguration.Name, metav1.GetOptions{})
	if client.IgnoreNotFound(err) != nil {
		setupLog.Error(err, "error getting current mutation webhook configuration")
		os.Exit(1)
	}

	mwh.DeepCopyInto(&mutatingWebhookConfiguration)

	failurePolicy := adminregv1.Fail
	sideEffects := adminregv1.SideEffectClassNone

	mutatingWebhookConfiguration.Labels = map[string]string{
		"apps.kubernetes.io/name": "autoimagepullsecrets-operator",
	}

	mutatingWebhookConfiguration.Webhooks = []adminregv1.MutatingWebhook{
		{
			Name: "mpod.autoimagepullsecrets.io",
			ClientConfig: adminregv1.WebhookClientConfig{
				CABundle: caSecret.Data["ca.crt"],
				Service: &adminregv1.ServiceReference{
					Name:      "aips-webhook",
					Namespace: caSecret.Namespace,
					Port:      pointer.Int32Ptr(9443),
					Path:      pointer.StringPtr("/mutate-v1-pod"),
				},
			},
			Rules: []adminregv1.RuleWithOperations{
				{
					Operations: []adminregv1.OperationType{
						adminregv1.Create,
						adminregv1.Update,
					},
					Rule: adminregv1.Rule{
						APIGroups:   []string{corev1.SchemeGroupVersion.Group},
						APIVersions: []string{corev1.SchemeGroupVersion.Version},
						Resources:   []string{"pods"},
					},
				},
			},
			SideEffects:             &sideEffects,
			AdmissionReviewVersions: []string{"v1beta1"},
			FailurePolicy:           &failurePolicy,
		},
	}

	if mutatingWebhookConfiguration.CreationTimestamp.IsZero() {
		if _, err := clientSet.
			AdmissionregistrationV1().
			MutatingWebhookConfigurations().
			Create(ctx, &mutatingWebhookConfiguration, metav1.CreateOptions{}); err != nil {
			setupLog.Error(err, "error creating webhook")
			os.Exit(1)
		}
	} else {
		if _, err := clientSet.
			AdmissionregistrationV1().
			MutatingWebhookConfigurations().
			Update(ctx, &mutatingWebhookConfiguration, metav1.UpdateOptions{}); err != nil {
			setupLog.Error(err, "error updating webhook")
			os.Exit(1)
		}
	}
}

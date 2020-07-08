package webhooks

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	adminregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	aipsv1alpha1 "github.com/logikone/autoimagepullsecrets-operator/api/v1alpha1"
)

var _ = Describe("IPS Pod Injector", func() {
	It("Should create ClusterDockerRegistry", func() {
		Expect(k8sClient.Create(context.Background(), &aipsv1alpha1.ClusterDockerRegistry{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test",
			},
			Spec: aipsv1alpha1.ClusterDockerRegistrySpec{
				DockerRegistrySpec: aipsv1alpha1.DockerRegistrySpec{
					AuthConfig: aipsv1alpha1.AuthConfig{
						ServerAddress: "https://hub.docker.io/api/v1",
						Username:      []byte("username"),
						Password:      []byte("password"),
					},
					Registry: "docker.io",
				},
			},
		})).Should(Succeed())
	})

	It("Should create MutatingPodWebhookConfiguration", func() {
		webhookURL := fmt.Sprintf("https://localhost:%d/mutate-v1-pod", webhookServer.Port)

		caBundle, err := ioutil.ReadFile(path.Join("..", "test", "certs", "tls.crt"))
		Expect(err).ToNot(HaveOccurred())

		noSideEffects := adminregv1.SideEffectClassNone

		Expect(k8sClient.Create(context.Background(), &adminregv1.MutatingWebhookConfiguration{
			ObjectMeta: metav1.ObjectMeta{
				Name: "ips-injector",
			},
			Webhooks: []adminregv1.MutatingWebhook{
				{
					Name: "mpod.autoimagepullsecrets.io",
					ClientConfig: adminregv1.WebhookClientConfig{
						URL:      &webhookURL,
						CABundle: caBundle,
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
								Resources:   []string{corev1.ResourcePods.String()},
							},
						},
					},
					AdmissionReviewVersions: []string{"v1beta1"},
					SideEffects:             &noSideEffects,
				},
			},
		})).Should(Succeed())
	})

	It("Should inject image pull secret", func() {
		testPod := corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nginx",
				Namespace: "default",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "nginx",
						Image: "nginx:latest",
					},
				},
			},
		}

		Expect(k8sClient.Create(context.Background(), &testPod)).To(Succeed())
		Expect(testPod.Spec.ImagePullSecrets).To(HaveLen(1))
	})
})

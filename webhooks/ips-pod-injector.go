package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	aipsv1alpha1 "github.com/logikone/autoimagepullsecrets-operator/api/v1alpha1"
)

type ImagePullSecretPodInjector struct {
	admission.Decoder
	client.Client
	record.EventRecorder
}

var (
	ipsInjectorLog               = ctrl.Log.WithName("webhook").WithName("ips-injector")
	IPSInjectionSecretNameFormat = "aips-registries-%s"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io

func (i *ImagePullSecretPodInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := ipsInjectorLog.WithValues("name", req.Name, "namespace", req.Namespace)
	var pod corev1.Pod

	if err := i.Decode(req, &pod); err != nil {
		return admission.Errored(http.StatusBadRequest,
			fmt.Errorf("unable to decode request to Pod: %w", err))
	}

	if err := i.MutatePod(ctx, &pod); err != nil {
		log.Error(err, "error finding match")
		return admission.Response{
			AdmissionResponse: v1beta1.AdmissionResponse{
				UID:     req.UID,
				Allowed: true,
			},
		}
	}

	marshalledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError,
			fmt.Errorf("error marshalling pod to json: %w", err))
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshalledPod)
}

func (i *ImagePullSecretPodInjector) FindMatches(ctx context.Context, pod corev1.Pod) ([]aipsv1alpha1.Registry, error) {
	log := ipsInjectorLog.WithValues("name", pod.Name, "namespace", pod.Namespace)

	var matched []aipsv1alpha1.Registry

	clusterDockerRegistryList, err := i.GetClusterDockerRegistries(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing cluster docker registries: %w", err)
	}

	dockerRegistryList, err := i.GetDockerRegistries(ctx, pod)
	if err != nil {
		return nil, fmt.Errorf("error listing docker registries: %w", err)
	}

	foundNamespacedRegistry := map[string]bool{}

	for _, container := range pod.Spec.Containers {
		registryName, err := DockerRegistryFromImage(container.Image)
		if err != nil {
			return nil, fmt.Errorf("error getting registrynName from image: %w", err)
		}

		for _, namespacedRegistry := range dockerRegistryList.Items {
			if _, ok := foundNamespacedRegistry[registryName]; !ok {
				foundNamespacedRegistry[registryName] = true
			}

			if namespacedRegistry.Spec.Registry == registryName {
				matched = append(matched, &namespacedRegistry)
			}
		}

		for _, clusterRegistry := range clusterDockerRegistryList.Items {
			if _, ok := foundNamespacedRegistry[registryName]; ok {
				continue
			}

			if clusterRegistry.Spec.Registry == registryName {
				matched = append(matched, &clusterRegistry)
			}
		}

		log.V(1).Info("parsed registryName from container image", "registryName", registryName)
	}

	return matched, nil
}

func (i *ImagePullSecretPodInjector) GetDockerRegistries(ctx context.Context, pod corev1.Pod) (aipsv1alpha1.DockerRegistryList, error) {
	var dockerRegistryList aipsv1alpha1.DockerRegistryList
	err := i.List(ctx, &dockerRegistryList, client.InNamespace(pod.GetNamespace()))

	return dockerRegistryList, err
}

func (i *ImagePullSecretPodInjector) GetClusterDockerRegistries(ctx context.Context) (aipsv1alpha1.ClusterDockerRegistryList, error) {
	var clusterDockerRegistryList aipsv1alpha1.ClusterDockerRegistryList
	err := i.List(ctx, &clusterDockerRegistryList)
	return clusterDockerRegistryList, err
}

func (i *ImagePullSecretPodInjector) InjectDecoder(decoder *admission.Decoder) error {
	i.Decoder = *decoder
	return nil
}

func (i *ImagePullSecretPodInjector) MutatePod(ctx context.Context, pod *corev1.Pod) error {
	// log := ipsInjectorLog.WithValues("name", pod.Name, "namespace", pod.Namespace)

	matches, err := i.FindMatches(ctx, *pod)
	if err != nil {
		return err
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}

	var clusterMatches []string
	var namespacedMatches []string

	for _, match := range matches {
		namespace := match.GetNamespace()

		if namespace == "" {
			clusterMatches = append(clusterMatches, match.GetName())
		} else {
			namespacedMatches = append(namespacedMatches, types.NamespacedName{
				Name:      match.GetName(),
				Namespace: match.GetNamespace(),
			}.String())
		}
	}

	if IPSCanInject(*pod) {
		pod.Annotations[IPSInjectionEnabled] = "true"
		if clusterMatches != nil {
			pod.Annotations[IPSInjectionMatch] = strings.Join(clusterMatches, ",")
		}

		if namespacedMatches != nil {
			pod.Annotations[IPSInjectionMatchNamespaced] = strings.Join(namespacedMatches, ",")
		}

		if clusterMatches != nil || namespacedMatches != nil {
			pod.Spec.ImagePullSecrets = append(pod.Spec.ImagePullSecrets, corev1.LocalObjectReference{
				Name: fmt.Sprintf(IPSInjectionSecretNameFormat, pod.GetName()),
			})
		}
	}

	return nil
}

func (i *ImagePullSecretPodInjector) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetWebhookServer().Register("/mutate-v1-pod", &webhook.Admission{
		Handler: i,
	})

	return nil
}

func IPSCanInject(pod corev1.Pod) bool {
	if pod.Annotations == nil {
		return true
	}

	injectionEnabled, ok := pod.Annotations[IPSInjectionEnabled]
	if !ok {
		return true
	}

	if injectionEnabled == "false" {
		return false
	}

	return true
}

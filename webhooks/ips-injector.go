package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type ImagePullSecretInjector struct {
	admission.Decoder
	client.Client
	record.EventRecorder
}

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io

func (i *ImagePullSecretInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	var pod corev1.Pod

	if err := i.Decode(req, &pod); err != nil {
		return admission.Errored(http.StatusBadRequest,
			fmt.Errorf("unable to decode request to Pod: %w", err))
	}

	marshalledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError,
			fmt.Errorf("error marshalling pod to json: %w", err))
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshalledPod)
}

func (i *ImagePullSecretInjector) InjectDecoder(decoder *admission.Decoder) error {
	i.Decoder = *decoder
	return nil
}

func (i *ImagePullSecretInjector) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetWebhookServer().Register("/mutate-v1-pod", &webhook.Admission{
		Handler: i,
	})

	return nil
}

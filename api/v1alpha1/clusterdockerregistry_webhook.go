package v1alpha1

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func (in *ClusterDockerRegistry) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

var _ webhook.Validator = &ClusterDockerRegistry{}

func (in *ClusterDockerRegistry) ValidateCreate() error {
	var errorsList field.ErrorList

	in.checkExists(&errorsList)

	if len(errorsList) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: in.GroupVersionKind().Group,
			Kind:  in.GroupVersionKind().Kind,
		},
		in.Name,
		errorsList)
}

func (in *ClusterDockerRegistry) ValidateUpdate(_ runtime.Object) error {
	var errorsList field.ErrorList

	in.checkExists(&errorsList)

	if len(errorsList) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: in.GroupVersionKind().Group,
			Kind:  in.GroupVersionKind().Kind,
		},
		in.Name,
		errorsList)
}

func (in *ClusterDockerRegistry) ValidateDelete() error {
	return nil
}

func (in *ClusterDockerRegistry) checkExists(list *field.ErrorList) {
	var clusterDockerRegistryList ClusterDockerRegistryList

	if err := Client.List(context.Background(), &clusterDockerRegistryList); err != nil {
		*list = append(*list, field.InternalError(
			field.NewPath("spec.name", "spec.namespace"), err))
	}

	for _, registry := range clusterDockerRegistryList.Items {
		if registry.Spec.Registry == in.Spec.Registry {
			*list = append(*list, field.Invalid(
				field.NewPath("spec.registry"), in.Spec.Registry,
				"registry already exists"))
		}
	}
}

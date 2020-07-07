package v1alpha1

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var _ webhook.Validator = &DockerRegistry{}

func (in *DockerRegistry) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(in).
		Complete()
}

func (in *DockerRegistry) ValidateCreate() error {
	var errorList field.ErrorList

	in.checkExists(&errorList)

	if len(errorList) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: in.GroupVersionKind().Group,
			Kind:  in.GroupVersionKind().Kind,
		},
		in.Name,
		errorList)
}

func (in *DockerRegistry) ValidateUpdate(_ runtime.Object) error {
	var errorList field.ErrorList

	in.checkExists(&errorList)

	if len(errorList) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{
			Group: in.GroupVersionKind().Group,
			Kind:  in.GroupVersionKind().Kind,
		},
		in.Name,
		errorList)
}

func (in *DockerRegistry) ValidateDelete() error {
	return nil
}

func (in *DockerRegistry) checkExists(list *field.ErrorList) {
	var dockerRegistryList DockerRegistryList

	if err := Client.List(context.Background(), &dockerRegistryList, client.InNamespace(in.Namespace)); err != nil {
		*list = append(*list, field.InternalError(
			field.NewPath("spec.name", "spec.namespace"), err))
	}

	for _, registry := range dockerRegistryList.Items {
		if registry.Spec.Registry == in.Spec.Registry {
			*list = append(*list, field.Invalid(
				field.NewPath("spec.registry"), in.Spec.Registry,
				"registry already exists"))
		}
	}
}

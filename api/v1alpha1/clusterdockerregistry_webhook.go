package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var _ webhook.Validator = &ClusterDockerRegistry{}

func (in *ClusterDockerRegistry) ValidateCreate() error {
	errorsList := field.ErrorList{}

	in.checkOnlyOne(errorsList)
	in.checkRegexp(errorsList)

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
	errorsList := field.ErrorList{}

	in.checkOnlyOne(errorsList)
	in.checkRegexp(errorsList)

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

func (in *ClusterDockerRegistry) checkRegexp(errorsList field.ErrorList) {
	if in.Spec.Registry.RegExp != "" {
		if !in.Spec.Registry.RegExp.IsValid() {
			errorsList = append(errorsList, field.Invalid(
				field.NewPath("spec.registry.regExp"), in.Spec.Registry.RegExp, "invalid regular expression"))
		}
	}
}

func (in *ClusterDockerRegistry) checkOnlyOne(errorsList field.ErrorList) {
	setCount := 0

	if in.Spec.Registry.RegExp != "" {
		setCount += 1
	}

	if in.Spec.Registry.Suffix != "" {
		setCount += 1
	}

	if in.Spec.Registry.Prefix != "" {
		setCount += 1
	}

	if in.Spec.Registry.Exact != "" {
		setCount += 1
	}

	if setCount > 1 {
		errorsList = append(errorsList, field.Invalid(
			field.NewPath("spec.registry"), in.Spec.Registry, "can only specify one of RegExp, Suffix, Prefix, or Exact"))
	}
}

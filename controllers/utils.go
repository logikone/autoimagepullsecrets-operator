package controllers

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func IsManagedSecret(obj metav1.Object) bool {
	if isManaged, ok := obj.GetAnnotations()[ManagedSecretAnnotation]; ok {
		if isManaged == True {
			return true
		}
	}

	return false
}

func IsSourceSecret(obj metav1.Object) bool {
	if isSource, ok := obj.GetAnnotations()[SourceSecretAnnotation]; ok {
		if isSource == True {
			return true
		}
	}

	return false
}

func GetNamespacedName(in metav1.Object) (types.NamespacedName, error) {
	var namespacedName types.NamespacedName

	val, ok := in.GetAnnotations()[SourceSecretAnnotation]
	if !ok {
		return namespacedName, fmt.Errorf(
			"annotation [%s] not found: %w", SourceSecretAnnotation, ErrMapKeyNotFound)
	}

	return ParseNamespacedName(val)
}

func ParseNamespacedName(in string) (types.NamespacedName, error) {
	var namespacedName types.NamespacedName

	separator := fmt.Sprintf("%c", types.Separator)

	if !strings.Contains(in, separator) {
		return namespacedName, fmt.Errorf(
			"string does not contain a '%s' character as expected: %w", separator, ErrStringNotContains)

	}

	separated := strings.Split(in, separator)

	if len(separated) != 2 {
		return namespacedName, fmt.Errorf(
			"unexpected slice length after separating [%d]: %w", len(separated), ErrUnexpectedSliceLength)
	}

	namespacedName.Namespace = separated[0]
	namespacedName.Name = separated[1]

	return namespacedName, nil
}

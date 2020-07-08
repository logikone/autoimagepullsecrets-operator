package controllers

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/logikone/autoimagepullsecrets-operator/webhooks"
)

func IsManagedSecret(obj metav1.Object) bool {
	if isManaged, ok := obj.GetAnnotations()[ManagedSecretAnnotation]; ok {
		if isManaged == True {
			return true
		}
	}

	return false
}

func IsSource(obj metav1.Object) bool {
	if isSource, ok := obj.GetAnnotations()[SourceAnnotation]; ok {
		if isSource == True {
			return true
		}
	}

	return false
}

func GetNamespacedName(in metav1.Object) (types.NamespacedName, error) {
	var namespacedName types.NamespacedName

	val, ok := in.GetAnnotations()[SourceAnnotation]
	if !ok {
		return namespacedName, fmt.Errorf(
			"annotation [%s] not found: %w", SourceAnnotation, ErrMapKeyNotFound)
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

func PodRequiresSecret(in runtime.Object) bool {
	obj, ok := in.(metav1.Object)
	if !ok {
		return false
	}

	if val, ok := obj.GetAnnotations()[webhooks.IPSInjectionEnabled]; !ok {
		return false
	} else {
		return val == "true"
	}
}

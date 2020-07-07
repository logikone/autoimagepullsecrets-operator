package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

// +kubebuilder:object:generate=false
type Registry interface {
	k8sruntime.Object
	metav1.Object

	GetAuthConfig() AuthConfig
	IsNamespaced() bool
	Render(current []byte) ([]byte, error)
}

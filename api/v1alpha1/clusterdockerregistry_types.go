/*
Copyright 2020 Chris Larsen.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DockerRegistrySpec defines the desired state of DockerRegistry
type ClusterDockerRegistrySpec struct {
	DockerRegistrySpec `json:",inline"`
}

// DockerRegistryStatus defines the observed state of DockerRegistry
type ClusterDockerRegistryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster

// DockerRegistry is the Schema for the dockerregistries API
type ClusterDockerRegistry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterDockerRegistrySpec   `json:"spec,omitempty"`
	Status ClusterDockerRegistryStatus `json:"status,omitempty"`
}

func (in *ClusterDockerRegistry) GetAuthConfig() AuthConfig {
	return in.Spec.AuthConfig
}

func (in *ClusterDockerRegistry) IsNamespaced() bool {
	return true
}

// +kubebuilder:object:root=true

// DockerRegistryList contains a list of DockerRegistry
type ClusterDockerRegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DockerRegistry `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterDockerRegistry{}, &ClusterDockerRegistryList{})
}

// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthConfig) DeepCopyInto(out *AuthConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthConfig.
func (in *AuthConfig) DeepCopy() *AuthConfig {
	if in == nil {
		return nil
	}
	out := new(AuthConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerConfigFile) DeepCopyInto(out *DockerConfigFile) {
	*out = *in
	if in.AuthConfigs != nil {
		in, out := &in.AuthConfigs, &out.AuthConfigs
		*out = make(map[string]AuthConfig, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerConfigFile.
func (in *DockerConfigFile) DeepCopy() *DockerConfigFile {
	if in == nil {
		return nil
	}
	out := new(DockerConfigFile)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerRegistry) DeepCopyInto(out *DockerRegistry) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerRegistry.
func (in *DockerRegistry) DeepCopy() *DockerRegistry {
	if in == nil {
		return nil
	}
	out := new(DockerRegistry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DockerRegistry) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerRegistryList) DeepCopyInto(out *DockerRegistryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DockerRegistry, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerRegistryList.
func (in *DockerRegistryList) DeepCopy() *DockerRegistryList {
	if in == nil {
		return nil
	}
	out := new(DockerRegistryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DockerRegistryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerRegistrySpec) DeepCopyInto(out *DockerRegistrySpec) {
	*out = *in
	out.AuthConfig = in.AuthConfig
	if in.Registry != nil {
		in, out := &in.Registry, &out.Registry
		*out = new(Matcher)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerRegistrySpec.
func (in *DockerRegistrySpec) DeepCopy() *DockerRegistrySpec {
	if in == nil {
		return nil
	}
	out := new(DockerRegistrySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DockerRegistryStatus) DeepCopyInto(out *DockerRegistryStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DockerRegistryStatus.
func (in *DockerRegistryStatus) DeepCopy() *DockerRegistryStatus {
	if in == nil {
		return nil
	}
	out := new(DockerRegistryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Matcher) DeepCopyInto(out *Matcher) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Matcher.
func (in *Matcher) DeepCopy() *Matcher {
	if in == nil {
		return nil
	}
	out := new(Matcher)
	in.DeepCopyInto(out)
	return out
}

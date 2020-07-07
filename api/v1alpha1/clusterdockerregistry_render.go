package v1alpha1

import (
	"encoding/json"
	"fmt"
)

func (in *ClusterDockerRegistry) Render(current []byte) ([]byte, error) {
	var currentDockerConfigFile dockerConfigFile

	if current != nil {
		if err := json.Unmarshal(current, &currentDockerConfigFile); err != nil {
			return nil, fmt.Errorf("error unmarshalling current docker config: %w", err)
		}
	}

	if currentDockerConfigFile.AuthConfigs == nil {
		currentDockerConfigFile.AuthConfigs = map[string]dockerAuthConfig{}
	}

	currentDockerConfigFile.AuthConfigs[in.Spec.AuthConfig.ServerAddress] = dockerAuthConfig{
		Username:      string(in.Spec.AuthConfig.Username),
		Password:      string(in.Spec.AuthConfig.Password),
		ServerAddress: in.Spec.AuthConfig.ServerAddress,
		IdentityToken: string(in.Spec.AuthConfig.IdentityToken),
		RegistryToken: string(in.Spec.AuthConfig.RegistryToken),
	}

	return json.Marshal(&currentDockerConfigFile)
}

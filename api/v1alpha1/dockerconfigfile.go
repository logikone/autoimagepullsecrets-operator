package v1alpha1

type DockerConfigFile struct {
	AuthConfigs map[string]AuthConfig `json:"auths"`
}

// AuthConfig contains authorization information for connecting to a Registry
type AuthConfig struct {
	// +optional
	Username string `json:"username,omitempty"`

	// +optional
	Password string `json:"password,omitempty"`

	// +optional
	Auth string `json:"auth,omitempty"`

	// Email is an optional value associated with the username.
	// This field is deprecated and will be removed in a later
	// version of docker.
	// +optional
	Email string `json:"email,omitempty"`

	ServerAddress string `json:"serveraddress,omitempty"`

	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	// +optional
	IdentityToken string `json:"identitytoken,omitempty"`

	// RegistryToken is a bearer token to be sent to a registry
	// +optional
	RegistryToken string `json:"registrytoken,omitempty"`
}

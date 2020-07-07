package v1alpha1

type dockerConfigFile struct {
	AuthConfigs map[string]AuthConfig `json:"auths"`
}

type dockerAuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Auth          string `json:"auth,omitempty"`
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
	IdentityToken string `json:"identitytoken,omitempty"`
	RegistryToken string `json:"registrytoken,omitempty"`
}

// AuthConfig contains authorization information for connecting to a Registry
type AuthConfig struct {
	// +optional
	Username []byte `json:"username,omitempty"`

	// +optional
	Password []byte `json:"password,omitempty"`

	ServerAddress string `json:"serveraddress,omitempty"`

	// IdentityToken is used to authenticate the user and get
	// an access token for the registry.
	// +optional
	IdentityToken []byte `json:"identitytoken,omitempty"`

	// RegistryToken is a bearer token to be sent to a registry
	// +optional
	RegistryToken []byte `json:"registrytoken,omitempty"`
}

package v1alpha1

// Matcher defines match types
type Matcher struct {
	// Match string exactly
	//
	// +optional
	Exact string `json:"exact,omitempty"`

	// Match strings with prefix
	//
	// +optional
	Prefix string `json:"prefix,omitempty"`

	// Match strings with suffix
	//
	// +optional
	Suffix string `json:"suffix,omitempty"`

	// Match strings matching a regular expression
	//
	// +optional
	RegExp string `json:"regExp,omitempty"`
}

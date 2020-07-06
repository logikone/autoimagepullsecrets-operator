package v1alpha1

import (
	"regexp"
	"strings"
)

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
	RegExp MatcherRegExp `json:"regExp,omitempty"`
}

func (in *Matcher) Matches(registry string) bool {
	if registry == in.Exact {
		return true
	}

	if strings.HasPrefix(registry, in.Prefix) {
		return true
	}

	if strings.HasSuffix(registry, in.Suffix) {
		return true
	}

	return in.RegExp.Matches(registry)
}

type MatcherRegExp string

func (in MatcherRegExp) String() string {
	return string(in)
}

func (in MatcherRegExp) IsValid() bool {
	if _, err := regexp.Compile(in.String()); err != nil {
		return false
	}

	return true
}

func (in MatcherRegExp) Matches(registry string) bool {
	re, err := regexp.Compile(in.String())
	if err != nil {
		return false
	}

	return re.MatchString(registry)
}

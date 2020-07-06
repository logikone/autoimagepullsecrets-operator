package controllers

import (
	"fmt"
)

const (
	ManagedSecretAnnotation = "autoimagepullsecrets.io/managed"
	SourceSecretAnnotation  = "autoimagepullsecrets.io/source"

	True  = "true"
	False = "false"
)

var (
	ManagedSecretIndex = fmt.Sprintf(".metadata.annotations.%s", ManagedSecretAnnotation)
	SourceSecretIndex  = fmt.Sprintf(".metadata.annotations.%s", SourceSecretAnnotation)
)

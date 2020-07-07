package controllers

import (
	"fmt"
)

const (
	ManagedSecretAnnotation = "autoimagepullsecrets.io/managed"
	SourceAnnotation        = "autoimagepullsecrets.io/source"
	SourceTypeAnnotation    = "autoimagepullsecrets.io/source-type"

	True  = "true"
	False = "false"
)

var (
	ManagedSecretIndex = fmt.Sprintf(".metadata.annotations.%s", ManagedSecretAnnotation)
	SourceSecretIndex  = fmt.Sprintf(".metadata.annotations.%s", SourceAnnotation)
)

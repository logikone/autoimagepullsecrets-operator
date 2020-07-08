package webhooks

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
)

func DockerRegistryFromImage(image string) (string, error) {
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return "", err
	}

	separated := strings.Split(ref.Name(), "/")

	if len(separated) < 2 {
		return "", fmt.Errorf("error getting registry from image")
	}

	return separated[0], nil
}

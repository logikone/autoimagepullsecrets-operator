package webhooks

import (
	"github.com/docker/distribution/reference"
)

func DockerRegistryFromImage(image string) (string, error) {
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return "", err
	}

	return ref.Name(), nil
}

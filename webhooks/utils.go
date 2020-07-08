package webhooks

import (
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
	corev1 "k8s.io/api/core/v1"
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

func IPSCanInject(pod corev1.Pod) bool {
	if pod.Annotations == nil {
		return true
	}

	injectionEnabled, ok := pod.Annotations[IPSInjectionEnabled]
	if !ok {
		return true
	}

	if injectionEnabled == "false" {
		return false
	}

	return true
}

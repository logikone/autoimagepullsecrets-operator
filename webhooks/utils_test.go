package webhooks

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Webhook Utils", func() {
	Context("DockerRegistryFromImage", func() {
		It("Should parse library image", func() {
			image := "nginx/alpine:latest"

			got, err := DockerRegistryFromImage(image)
			Expect(err).ToNot(HaveOccurred())
			Expect(got).ToNot(BeNil())
			Expect(got).To(Equal("docker.io"))
		})

		It("Should parse image from custom repository", func() {
			image := "docker.artifactory.jfrog.io/test/image:latest"

			got, err := DockerRegistryFromImage(image)
			Expect(err).ToNot(HaveOccurred())
			Expect(got).ToNot(BeNil())
			Expect(got).To(Equal("docker.artifactory.jfrog.io"))
		})
	})
})

package cmd

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("getExistingConfig", func() {
	Context("with empty file", func() {
		It("should return an empty []byte", func() {
			result, err := getExistingConfig("./test_data/.empty_envrc")

			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(BeEmpty())
		})
	})

	Context("with file with existing config", func() {
		It("should return the existing config contents", func() {
			result, err := getExistingConfig("./test_data/.existing_config_envrc")

			Expect(err).NotTo(HaveOccurred())

			Expect(string(result)).To(Equal("first line\n# comment\nsecond line\n"))
		})
	})

	Context("with file with existing config and a TinyHelper comment", func() {
		It("should return the existing config contents but remove the th comment text", func() {
			result, err := getExistingConfig("./test_data/.existing_config_envrc")

			Expect(err).NotTo(HaveOccurred())

			Expect(string(result)).To(Equal("first line\n# comment\nsecond line\n"))
		})
	})
})

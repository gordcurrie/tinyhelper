package cmd

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("env", func() {
	Describe("parseInfo", func() {
		Context("when passed valid resluts", func() {
			var expected data
			BeforeEach(func() {
				expected = data{
					Goroot: "/home/gord/.cache/tinygo/goroot-d3a5eae46885c758dc170cc3b2ebb723ef9c0181c18efbe4c5dc3ba26d61a5ae",
					Flags:  "-tags=cortexm,baremetal,linux,arm,rp2040,rp,pico,tinygo,math_big_pure_go,gc.conservative,scheduler.tasks,serial.usb",
					Target: "pico",
				}
			})

			It("should properly parse those results", func() {
				result := parseInfo(validInfo, "pico")

				Expect(result).To(Equal(expected))
			})
		})
	})

	Describe("getExistingConfig", func() {
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
				result, err := getExistingConfig("./test_data/.existing_config_with_th_comment.envrc")

				Expect(err).NotTo(HaveOccurred())

				Expect(string(result)).To(Equal("first line\n# comment\nsecond line\n"))
			})
		})
	})
})

var validInfo = `
LLVM triple:       thumbv6m-unknown-unknown-eabi
GOOS:              linux
GOARCH:            arm
build tags:        cortexm baremetal linux arm rp2040 rp pico tinygo math_big_pure_go gc.conservative scheduler.tasks serial.usb
garbage collector: conservative
scheduler:         tasks
cached GOROOT:     /home/gord/.cache/tinygo/goroot-d3a5eae46885c758dc170cc3b2ebb723ef9c0181c18efbe4c5dc3ba26d61a5ae

`

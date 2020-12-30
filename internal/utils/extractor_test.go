package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/electrocucaracha/pkg-mgr/internal/utils"
)

var _ = Describe("Extractor", func() {
	var (
		parsedFunction map[string]string
		bashScript     string
		err            error
	)

	JustBeforeEach(func() {
		parsedFunction, err = utils.ExtractFunctions(bashScript)
	})

	Describe("parsing content of a function", func() {
		Context("when the script doesn't have a function", func() {
			BeforeEach(func() {
				bashScript = "test"
			})

			It("should retrieve an empty object", func() {
				Expect(parsedFunction).To(Equal(map[string]string{}))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the script has more than one function", func() {
			BeforeEach(func() {
				bashScript = `#!/bin/bash

function test {
    echo test
}

function main {
    test
    sudo tee /etc/docker/daemon.json << EOF
{
  "key" : "value"
}
EOF
}

main`
			})

			It("should retrieve a valid function's object", func() {
				Expect(parsedFunction).NotTo(Equal(BeNil()))
				Expect(parsedFunction).To(HaveLen(2))
				Expect(parsedFunction).To(HaveKeyWithValue("main", `    test
    sudo tee /etc/docker/daemon.json << EOF
{
  "key" : "value"
}
EOF`))
				Expect(parsedFunction).To(HaveKeyWithValue("test", "    echo test"))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

})

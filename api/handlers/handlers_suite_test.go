package handlers_test

import (
	"testing"

	"github.com/electrocucaracha/pkg-mgr/gen/restapi"
	"github.com/go-openapi/loads"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	swaggerSpec *loads.Document
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = BeforeSuite(func() {
	var err error
	swaggerSpec, err = loads.Analyzed(restapi.SwaggerJSON, "")
	Expect(err).NotTo(HaveOccurred())
})

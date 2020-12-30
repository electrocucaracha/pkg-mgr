package handlers_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/electrocucaracha/pkg-mgr/api/handlers"
	"github.com/electrocucaracha/pkg-mgr/gen/restapi"
	"github.com/electrocucaracha/pkg-mgr/gen/restapi/operations"
	"github.com/electrocucaracha/pkg-mgr/internal/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockDB struct {
	Items map[string]models.Bash
	Err   error
}

func newMockDB() *mockDB {
	return &mockDB{
		Items: map[string]models.Bash{
			"existing_pkg": models.Bash{Pkg: "existing_pkg",
				Functions: []models.Function{models.Function{
					Name:    "main",
					Content: "    echo test",
				}},
			},
			handlers.MainBashPackage: models.Bash{Pkg: handlers.MainBashPackage,
				Functions: []models.Function{models.Function{
					Name:    "update_repo",
					Content: "    apt-get update",
				}},
			},
		},
	}
}

func (db *mockDB) GetScript(pkg string) (*models.Bash, error) {
	if db.Err != nil {
		return nil, db.Err
	}
	if i, ok := db.Items[pkg]; ok {
		return &i, nil
	}

	return nil, nil
}

func (db *mockDB) CreateScript(pkg, instructionSet string) (*models.Bash, []error) {
	return nil, nil
}

var _ = Describe("Bash handler", func() {
	var (
		err         error
		recorder    *httptest.ResponseRecorder
		queryParams map[string]string
	)
	uri := "/install_pkg"

	JustBeforeEach(func() {
		api := operations.NewPkgMgrAPI(swaggerSpec)
		api.GetScriptHandler = handlers.NewGetBash(newMockDB())
		err = api.Validate()
		Expect(err).NotTo(HaveOccurred())

		server := restapi.NewServer(api)
		server.ConfigureAPI()
		handler, _ := api.HandlerFor(http.MethodGet, uri)
		testserver := httptest.NewServer(handler)
		defer testserver.Close()

		request := httptest.NewRequest(http.MethodGet, uri, nil)
		if queryParams != nil {
			q := request.URL.Query()
			for k, v := range queryParams {
				q.Add(k, v)
			}
			request.URL.RawQuery = q.Encode()
		}
		Expect(err).NotTo(HaveOccurred())
		recorder = httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
	})

	Describe("retrieving scripts from the datastore", func() {
		Context("when no params are passed", func() {
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return a HTTP error code", func() {
				Expect(recorder).ToNot(BeNil())
				Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			})
		})

		Context("when the desired package doesn't exist", func() {
			BeforeEach(func() {
				queryParams = map[string]string{
					"pkg": "non_existing_pkg",
				}
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return a not found HTTP code", func() {
				Expect(recorder).ToNot(BeNil())
				Expect(recorder.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when the desired package exists", func() {
			BeforeEach(func() {
				queryParams = map[string]string{
					"pkg": "existing_pkg",
				}
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return a successful script", func() {
				By("returning an OK status HTTP code")
				Expect(recorder).ToNot(BeNil())
				Expect(recorder.Code).To(Equal(http.StatusOK))

				By("parsing the properly the script")
				body, err := ioutil.ReadAll(recorder.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(body).ToNot(BeNil())
				Expect(string(body)).To(Equal(fmt.Sprintf(`%s

%s

function main {
    echo test
}

main`, handlers.Header, handlers.Setters)))
			})
		})

		Context("when the desired package exists and the update function is requested", func() {
			BeforeEach(func() {
				queryParams = map[string]string{
					"pkg":        "existing_pkg",
					"pkg_update": "true",
				}
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should return a successful script", func() {
				By("returning an OK status HTTP code")
				Expect(recorder).ToNot(BeNil())
				Expect(recorder.Code).To(Equal(http.StatusOK))

				By("parsing the properly the script")
				body, err := ioutil.ReadAll(recorder.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(body).ToNot(BeNil())
				Expect(string(body)).To(Equal(fmt.Sprintf(`%s

%s

function main {
    echo test
}

function update_repo {
    apt-get update
}

update_repos

main`, handlers.Header, handlers.Setters)))
			})
		})
	})
})

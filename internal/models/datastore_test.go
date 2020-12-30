package models_test

import (
	"database/sql/driver"
	"regexp"
	"time"

	"github.com/electrocucaracha/pkg-mgr/internal/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

var _ = Describe("Datasource", func() {
	var (
		bashScript      *models.Bash
		packageName     string
		packageFunction string
		err             error
		errs            []error
	)

	Describe("retrieving stored bash scripts", func() {
		JustBeforeEach(func() {
			bashScript, err = datastore.GetScript("test")
		})

		Context("when the bash script exists and doesn't have any function", func() {
			BeforeEach(func() {
				returnedRows := sqlmock.NewRows([]string{"pkg"})
				returnedRows.AddRow("test")
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `bashes` WHERE `bashes`.`deleted_at` IS NULL AND ((pkg = ?))")).WithArgs(
					"test").WillReturnRows(returnedRows)
			})

			It("should retrieve the object properly", func() {
				Expect(bashScript.Pkg).To(Equal("test"))
			})
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the script doesn't exist", func() {
			BeforeEach(func() {
				returnedRows := sqlmock.NewRows([]string{"pkg"})
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `bashes` WHERE `bashes`.`deleted_at` IS NULL AND ((pkg = ?))")).WithArgs(
					"test").WillReturnRows(returnedRows)
			})

			It("should return the zero-value for the script", func() {
				Expect(bashScript).Should(BeNil())
			})
		})
	})

	Describe("adding bash scripts", func() {
		JustBeforeEach(func() {
			bashScript, errs = datastore.CreateScript(packageName, packageFunction)
		})

		Context("when a bash script is created sucessfully", func() {
			BeforeEach(func() {
				packageName = "test"
				packageFunction = `#!/bin/bash

function main {
    echo test
}

main`
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `bashes` (`created_at`,`updated_at`,`deleted_at`,`pkg`) VALUES (?,?,?,?)")).
					WithArgs(AnyTime{}, AnyTime{}, nil, packageName).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `functions` (`created_at`,`updated_at`,`deleted_at`,`bash_id`,`name`,`content`) VALUES (?,?,?,?,?,?)")).
					WithArgs(AnyTime{}, AnyTime{}, nil, 1, "main", "    echo test").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			})

			It("should store and return the object properly serialized", func() {
				Expect(bashScript).NotTo(BeNil())
				Expect(bashScript.Pkg).To(Equal(packageName))
				Expect(bashScript.Functions).To(HaveLen(1))
				Expect(bashScript.Functions[0].Name).To(Equal("main"))
				Expect(bashScript.Functions[0].Content).To(Equal("    echo test"))
			})
			It("should not error", func() {
				Expect(errs).To(HaveLen(0))
			})
		})
	})
})

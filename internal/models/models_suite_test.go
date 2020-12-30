package models_test

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/electrocucaracha/pkg-mgr/internal/models"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	db        *sql.DB
	mock      sqlmock.Sqlmock
	datastore models.Datastore
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = BeforeSuite(func() {
	var err error
	db, mock, err = sqlmock.New()
	Expect(err).NotTo(HaveOccurred())

	orm, err := gorm.Open("mysql", db)
	Expect(err).NotTo(HaveOccurred())

	datastore, err = models.NewGormDatastore(orm)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	db.Close()
})

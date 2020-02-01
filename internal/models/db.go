package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"  // MySQL/MariaDB dialect for gorm
	_ "github.com/jinzhu/gorm/dialects/sqlite" // SQLite dialect for gorm
)

// Datastore provides the methods supported by different databases
type Datastore interface {
	GetScript(pkg string) (*Bash, error)
	CreateScript(pkg, instructionSet string) error
}

type gormDatastore struct {
	*gorm.DB
}

// NewSqliteDatastore creates a database connection for a SQLite engine
func NewSqliteDatastore(file string, debug bool) (Datastore, error) {
	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	db.LogMode(debug)
	db.AutoMigrate(&Bash{})

	return &gormDatastore{db}, nil
}

// NewMySqlDatastore creates a database connection for a MySQL/MariaDB engine
func NewMySqlDatastore(username, password, hostname, database string) (Datastore, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", username, password, hostname, database)
	log.Println(connectionString)
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Bash{})

	return &gormDatastore{db}, nil
}
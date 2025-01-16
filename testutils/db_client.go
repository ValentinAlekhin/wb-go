package testutils

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
)

func GetDBClient() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./db/test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

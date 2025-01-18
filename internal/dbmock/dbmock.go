package dbmock

import (
	db2 "github.com/ValentinAlekhin/wb-go/internal/db"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewDBMock() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&db2.ControlModel{})
	if err != nil {
		panic(err)
	}

	return db
}

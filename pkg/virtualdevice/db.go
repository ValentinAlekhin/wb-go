package virtualdevice

import (
	db2 "github.com/ValentinAlekhin/wb-go/internal/db"
	"gorm.io/gorm"
	"sync"
)

var migrationOnce sync.Once
var migrationErr error

func migrate(db *gorm.DB) error {
	migrationOnce.Do(func() {
		migrationErr = db.AutoMigrate(&db2.ControlModel{})
	})

	return migrationErr
}

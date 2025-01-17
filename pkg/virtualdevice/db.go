package virtualdevice

import (
	"github.com/ValentinAlekhin/wb-go/pkg/virualcontrol"
	"gorm.io/gorm"
	"sync"
)

var migrationOnce sync.Once
var migrationErr error

func migrate(db *gorm.DB) error {
	migrationOnce.Do(func() {
		migrationErr = db.AutoMigrate(&virualcontrol.ControlModel{})
	})

	return migrationErr
}

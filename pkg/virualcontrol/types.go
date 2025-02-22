package virualcontrol

import (
	"github.com/ValentinAlekhin/wb-go/pkg/control"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"gorm.io/gorm"
)

type BaseOptions struct {
	DB     *gorm.DB
	Client wb.ClientInterface
	Device string
	Name   string
	Meta   control.Meta
}

package controls

import (
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type PushbuttonControl struct {
	control *Control
}

func (pbc *PushbuttonControl) Push() {
	pbc.control.SetValue(conventions.CONV_META_BOOL_TRUE)
}

func NewPushbuttonControl(client *wb.Client, device, control string) *PushbuttonControl {
	c := NewControl(client, device, control)

	return &PushbuttonControl{c}
}

package controls

import (
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type PushbuttonControl struct {
	control *Control
}

func (c *PushbuttonControl) Push() {
	c.control.SetValue(conventions.CONV_META_BOOL_TRUE)
}

func (c *PushbuttonControl) GetInfo() ControlInfo {
	return c.control.GetInfo()
}

func NewPushbuttonControl(client *wb.Client, device, control string, meta Meta) *PushbuttonControl {
	c := NewControl(client, device, control, meta)

	return &PushbuttonControl{c}
}

package control

import (
	"context"
	"github.com/ValentinAlekhin/wb-go/pkg/conventions"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type PushbuttonControl struct {
	control *Control
}

func (c *PushbuttonControl) Push() {
	c.control.SetValue(conventions.CONV_META_BOOL_TRUE)
}

func (c *PushbuttonControl) GetInfo() Info {
	return c.control.GetInfo()
}

func NewPushbuttonControl(ctx context.Context, client wb.ClientInterface, device, control string, meta Meta) *PushbuttonControl {
	c := NewControl(ctx, client, device, control, meta)

	return &PushbuttonControl{c}
}

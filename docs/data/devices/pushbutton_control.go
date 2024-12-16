package devices

import (
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type PushbuttonControl struct {
	control *Control
}

func (pbc *PushbuttonControl) Push() {
	pbc.control.SetValue("1")
}

func NewPushbuttonControl(client *wb.Client, topic string) *PushbuttonControl {
	control := NewControl(client, topic)

	return &PushbuttonControl{control}
}

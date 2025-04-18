package control

import (
	"context"
	"fmt"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
)

type RgbControl struct {
	converter Converter[RgbValue]
	control   ControlInterface
}

func (c *RgbControl) GetValue() RgbValue {
	value, err := c.converter.Decode(c.control.GetValue())
	if err != nil {
		fmt.Println(err)
	}

	return value
}

func (c *RgbControl) SetValue(value RgbValue) {
	c.control.SetValue(c.converter.Encode(value))
}

func (c *RgbControl) AddWatcher(f func(payload WatcherPayloadRGB)) {
	c.control.AddWatcher(func(p WatcherPayloadString) {
		newValue, err := c.converter.Decode(p.NewValue)
		if err != nil {
			fmt.Println(err)
		}
		oldValue, err := c.converter.Decode(p.OldValue)
		if err != nil {
			fmt.Println(err)
		}

		f(WatcherPayloadRGB{
			NewValue: newValue,
			OldValue: oldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *RgbControl) GetInfo() Info {
	return c.control.GetInfo()
}

func NewRgbControl(ctx context.Context, client wb.ClientInterface, device, control string, meta Meta) *RgbControl {
	c := NewControl(ctx, client, device, control, meta)
	return &RgbControl{&RGBConverter{}, c}
}

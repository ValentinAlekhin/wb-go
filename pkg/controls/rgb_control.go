package controls

import (
	"fmt"
	wb "github.com/ValentinAlekhin/wb-go/pkg/mqtt"
	"strconv"
	"strings"
)

type RgbControl struct {
	control *Control
}

type RgbValue struct {
	Red   int
	Green int
	Blue  int
}

type RgbControlWatcherPayload struct {
	NewValue    RgbValue
	OldValue    RgbValue
	Topic       string
	ControlName string
}

func (c *RgbControl) GetValue() RgbValue {
	value, err := c.decode(c.control.GetValue())
	if err != nil {
		fmt.Println(err)
	}

	return value
}

func (c *RgbControl) SetValue(value RgbValue) {
	c.control.SetValue(c.encode(value))
}

func (c *RgbControl) AddWatcher(f func(payload RgbControlWatcherPayload)) {
	c.control.AddWatcher(func(p ControlWatcherPayload) {
		newValue, err := c.decode(p.NewValue)
		if err != nil {
			fmt.Println(err)
		}
		oldValue, err := c.decode(p.OldValue)
		if err != nil {
			fmt.Println(err)
		}

		f(RgbControlWatcherPayload{
			NewValue: newValue,
			OldValue: oldValue,
			Topic:    p.Topic,
		})
	})
}

func (c *RgbControl) decode(value string) (RgbValue, error) {
	stringsValues := strings.Split(value, ";")

	red, err := strconv.Atoi(stringsValues[0])
	if err != nil {
		return RgbValue{}, fmt.Errorf("error converting red: %s", err)
	}

	green, err := strconv.Atoi(stringsValues[1])
	if err != nil {
		return RgbValue{}, fmt.Errorf("error converting green: %s", err)
	}

	blue, err := strconv.Atoi(stringsValues[2])
	if err != nil {
		return RgbValue{}, fmt.Errorf("error converting blue: %s", err)
	}

	return RgbValue{
		Red:   red,
		Green: green,
		Blue:  blue,
	}, nil
}

func (c *RgbControl) encode(value RgbValue) string {
	return fmt.Sprintf("%d;%d;%d", value.Red, value.Green, value.Blue)
}

func NewRgbControl(client *wb.Client, device, control string) *RgbControl {
	c := NewControl(client, device, control)
	return &RgbControl{c}
}

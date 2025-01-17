package basedevice

import "github.com/ValentinAlekhin/wb-go/pkg/control"

type Device struct {
	Name     string
	Controls interface{}
}

type Info struct {
	Name         string
	ControlsInfo []control.Info
	MetaTopic    string
}

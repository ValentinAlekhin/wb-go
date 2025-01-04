package deviceinfo

import "github.com/ValentinAlekhin/wb-go/pkg/control"

type DeviceInfo struct {
	Name         string
	Device       string
	Address      string
	ControlsInfo []control.ControlInfo
}

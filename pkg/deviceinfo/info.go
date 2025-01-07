package deviceinfo

import "github.com/ValentinAlekhin/wb-go/pkg/control"

type DeviceInfo struct {
	Name         string
	ControlsInfo []control.ControlInfo
}

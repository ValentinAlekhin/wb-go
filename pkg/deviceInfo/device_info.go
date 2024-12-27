package deviceInfo

import "github.com/ValentinAlekhin/wb-go/pkg/controls"

type DeviceInfo struct {
	Name         string
	Device       string
	Address      string
	ControlsInfo []controls.ControlInfo
}

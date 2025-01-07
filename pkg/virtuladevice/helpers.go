package virtuladevice

import "fmt"

func getDeviceFullName(name string) string {
	return fmt.Sprintf("%s_%s", DevicePrefix, name)
}

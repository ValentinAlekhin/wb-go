package deploy

import (
	"testing"
)

func TestRun(t *testing.T) {

	cfg := &Config{
		AppName: "test-app",
		Build: BuildConfig{
			Out: "./bin/app",
			Src: "../../",
		},
		Devices: []Device{
			{
				Host: "192.168.1.150",
			},
		},
	}

	err := Run(cfg)
	if err != nil {
		t.Errorf("%s", err)
	}

}

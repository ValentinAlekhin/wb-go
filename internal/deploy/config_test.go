package deploy

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name      string
		input     Config
		expectErr bool
	}{
		{
			name: "valid configuration",
			input: Config{
				AppName: "test-app",
				Build: BuildConfig{
					SrcDir: "./src",
					OutDir: "./out",
				},
				Devices: []Device{
					{
						Host:     "192.168.1.100",
						User:     "root",
						Password: "password",
						SshPort:  22,
						AppDir:   "/mnt/data/app",
					},
				},
			},
			expectErr: false,
		},
		{
			name: "missing app name",
			input: Config{
				AppName: "",
				Build: BuildConfig{
					SrcDir: "./src",
					OutDir: "./out",
				},
				Devices: []Device{
					{
						Host:     "192.168.1.100",
						User:     "root",
						Password: "password",
						SshPort:  22,
						AppDir:   "/mnt/data/app",
					},
				},
			},
			expectErr: true,
		},
		{
			name: "invalid SSH port",
			input: Config{
				AppName: "test-app",
				Build: BuildConfig{
					SrcDir: "./src",
					OutDir: "./out",
				},
				Devices: []Device{
					{
						Host:     "192.168.1.100",
						User:     "root",
						Password: "password",
						SshPort:  70000, // Неверный порт
						AppDir:   "/mnt/data/app",
					},
				},
			},
			expectErr: true,
		},
		{
			name: "empty devices",
			input: Config{
				AppName: "test-app",
				Build: BuildConfig{
					SrcDir: "./src",
					OutDir: "./out",
				},
				Devices: []Device{}, // Пустой список устройств
			},
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateConfig(test.input)

			if (err != nil) != test.expectErr {
				t.Errorf("expected error: %v, got: %v", test.expectErr, err)
			}
		})
	}
}

func TestApplyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    Config
		expected Config
	}{
		{
			name: "All fields empty",
			input: Config{
				AppName: "name",
				Build: BuildConfig{
					SrcDir: "./default-src",
					OutDir: "./default-out",
				},
				Devices: []Device{
					{
						Host:     "192.168.1.100",
						User:     "",
						Password: "",
						SshPort:  0,
						AppDir:   "",
					},
				},
			},
			expected: Config{
				AppName: "name",
				Build: BuildConfig{
					SrcDir: "./default-src",
					OutDir: "./default-out",
				},
				Devices: []Device{
					{
						Host:     "192.168.1.100",
						User:     SshUser,
						Password: SshPassword,
						SshPort:  SshPort,
						AppDir:   DefaultAppDir,
					},
				},
			},
		},
		{
			name: "Partial fields empty",
			input: Config{
				AppName: "custom_app",
				Build: BuildConfig{
					SrcDir: "./custom-src",
					OutDir: "./custom-out",
				},
				Devices: []Device{
					{
						Host:     "192.168.1.100",
						User:     "custom_user",
						Password: "",
						SshPort:  0,
						AppDir:   "",
					},
				},
			},
			expected: Config{
				AppName: "custom_app",
				Build: BuildConfig{
					SrcDir: "./custom-src",
					OutDir: "./custom-out",
				},
				Devices: []Device{
					{
						Host:     "192.168.1.100",
						User:     "custom_user",
						Password: SshPassword,
						SshPort:  SshPort,
						AppDir:   DefaultAppDir,
					},
				},
			},
		},
		{
			name: "No empty fields",
			input: Config{
				AppName: "custom_app",
				Build: BuildConfig{
					SrcDir: "./custom-src",
					OutDir: "./custom-out",
				},
				Devices: []Device{
					{
						Host:     "192.168.1.100",
						User:     "custom_user",
						Password: "custom_password",
						SshPort:  22,
						AppDir:   "/custom/app/dir",
					},
				},
			},
			expected: Config{
				AppName: "custom_app",
				Build: BuildConfig{
					SrcDir: "./custom-src",
					OutDir: "./custom-out",
				},
				Devices: []Device{
					{
						Host:     "192.168.1.100",
						User:     "custom_user",
						Password: "custom_password",
						SshPort:  22,
						AppDir:   "/custom/app/dir",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applyDefaults(&tt.input)

			// Сравнение с ожидаемым результатом
			if !equalConfig(tt.input, tt.expected) {
				t.Errorf("applyDefaults() = %+v, want %+v", tt.input, tt.expected)
			}
		})
	}
}

// Функция для сравнения двух Config структур
func equalConfig(c1, c2 Config) bool {
	if c1.AppName != c2.AppName {
		return false
	}
	if c1.Build.SrcDir != c2.Build.SrcDir || c1.Build.OutDir != c2.Build.OutDir {
		return false
	}
	if len(c1.Devices) != len(c2.Devices) {
		return false
	}
	for i := range c1.Devices {
		if c1.Devices[i].Host != c2.Devices[i].Host ||
			c1.Devices[i].User != c2.Devices[i].User ||
			c1.Devices[i].Password != c2.Devices[i].Password ||
			c1.Devices[i].SshPort != c2.Devices[i].SshPort ||
			c1.Devices[i].AppDir != c2.Devices[i].AppDir {
			return false
		}
	}
	return true
}

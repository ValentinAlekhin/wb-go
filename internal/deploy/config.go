package deploy

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
)

type BuildConfig struct {
	SrcDir string `mapstructure:"src_dir" validate:"required,dirpath"`
	OutDir string `mapstructure:"out_dir" validate:"required,dirpath"`
}

type Device struct {
	Host     string `mapstructure:"host" validate:"required"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	SshPort  int    `mapstructure:"ssh_port" validate:"required,min=1,max=65535"`
	AppDir   string `mapstructure:"app_dir" validate:"required"`
}

type Config struct {
	AppName    string      `mapstructure:"app_name" validate:"required"`
	Build      BuildConfig `mapstructure:"build" validate:"required"`
	Devices    []Device    `mapstructure:"devices" validate:"required,min=1,dive"`
	WorkingDir string
}

func getConfigPath(configPath string) (string, error) {
	if configPath != "" {
		return configPath, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("ошибка при получении текущей директории: %w", err)
	}

	return dir + "/" + DefaultConfigName, nil
}

// Функция для загрузки конфигурации с использованием Viper
func loadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурации: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("ошибка распаковки конфигурации: %w", err)
	}

	applyDefaults(&config)

	if err := setWorkingDir(&config, configPath); err != nil {
		return nil, fmt.Errorf("ошибка при установке рабочей директории: %w", err)
	}

	setBuildDirs(&config)

	return &config, nil
}

func applyDefaults(cfg *Config) {
	for i := range cfg.Devices {
		if cfg.Devices[i].SshPort == 0 {
			cfg.Devices[i].SshPort = SshPort
		}
		if cfg.Devices[i].User == "" {
			cfg.Devices[i].User = SshUser
		}
		if cfg.Devices[i].Password == "" {
			cfg.Devices[i].Password = SshPassword
		}
		if cfg.Devices[i].AppDir == "" {
			cfg.Devices[i].AppDir = DefaultAppDir
		}
	}
}

func setWorkingDir(config *Config, configPath string) error {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("ошибка при получениие абсолютного пути: %w", err)
	}

	config.WorkingDir = filepath.Dir(absPath)

	return nil
}

func setBuildDirs(config *Config) {
	config.Build.SrcDir = path.Join(config.WorkingDir, config.Build.SrcDir)
	config.Build.OutDir = path.Join(config.WorkingDir, config.Build.OutDir)
}

func validateConfig(cfg Config) error {
	validate := validator.New()
	return validate.Struct(cfg)
}

func GetConfig(configPath string) (*Config, error) {
	configPath, err := getConfigPath(configPath)
	if err != nil {
		return nil, err
	}

	config, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %v", err)
	}

	if err := validateConfig(*config); err != nil {
		return nil, fmt.Errorf("ошибка валидации конфигурации: %v", err)
	}

	return config, nil
}

package deploy

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
	"os"
	"os/exec"
	"path"
	"text/template"
	"time"
)

type systemdFileConfig struct {
	Description      string
	ExecStart        string
	WorkingDirectory string
	Environment      []systemdFileConfigEnvironment
}

type systemdFileConfigEnvironment struct {
	Key   string
	Value string
}

//go:embed templates/*
var embedFs embed.FS

func buildApp(config *Config) (*os.File, error) {
	binPath := path.Join(config.Build.OutDir, config.AppName)

	cmd := exec.Command("go", "build", "-o", binPath)
	cmd.Dir = config.Build.SrcDir
	cmd.Env = append(os.Environ(),
		"GOOS=linux",
		"GOARCH=arm",
		"GOARM=7",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Сборка приложения для Wiren Board 7...")
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("ошибка при сборке: %v", err)
	}

	file, err := os.Open(binPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка при открытии файла: %v", err)
	}

	fmt.Println("Сборка завершена успешно!")

	return file, nil
}

func createSSHClient(host, user, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к SSH: %v", err)
	}

	return client, nil
}

func uploadToDevice(client *ssh.Client, file *os.File, remotePath string) error {
	scpClient, _ := scp.NewClientBySSH(client)

	err := scpClient.CopyFromFile(context.Background(), *file, remotePath, "0655")
	if err != nil {
		return fmt.Errorf("ошибка при передаче файла через SCP: %v", err)
	}

	fmt.Println("Файл успешно передан на устройство.")
	return nil
}

func stopService(client *ssh.Client, config *Config) error {
	execRemote := createExecRemote(client)

	_, err := execRemote(fmt.Sprintf("systemctl stop %s.service || true", config.AppName))
	if err != nil {
		return fmt.Errorf("ошибка при остановке сервиса: %v", err)
	}

	return nil
}

func createAndStartService(client *ssh.Client, config *Config, device Device) error {
	systedConfig := &systemdFileConfig{
		Description:      fmt.Sprintf("%s Service", config.AppName),
		ExecStart:        device.AppDir + config.AppName,
		WorkingDirectory: device.AppDir,
		Environment:      []systemdFileConfigEnvironment{},
	}

	for key, value := range config.Environment {
		envRecord := systemdFileConfigEnvironment{
			Key:   key,
			Value: value,
		}
		systedConfig.Environment = append(systedConfig.Environment, envRecord)
	}

	tmpl, err := template.ParseFS(embedFs, systemdTemplateFile)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге шаблона systemd: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, systedConfig)
	if err != nil {
		return fmt.Errorf("ошибка при создании шаблона systemd: %v", err)
	}

	execRemote := createExecRemote(client)

	_, err = execRemote(fmt.Sprintf("echo '%s' > /etc/systemd/system/%s.service", buf.String(), config.AppName))
	if err != nil {
		return fmt.Errorf("ошибка при создании systemd-сервиса: %v", err)
	}

	msg, err := execRemote("systemctl daemon-reload")
	if err != nil {
		return fmt.Errorf("ошибка при остановке systemd-сервиса: %v", err)
	}

	msg, err = execRemote(fmt.Sprintf("systemctl enable %s.service", config.AppName))
	if err != nil {
		return fmt.Errorf("ошибка при активации systemd-сервиса: %v", err)
	}

	msg, err = execRemote(fmt.Sprintf("systemctl restart %s.service", config.AppName))
	if err != nil {
		return fmt.Errorf("ошибка при запуске systemd-сервиса: %v", err)
	}

	time.Sleep(3 * time.Second)

	msg, _ = execRemote(fmt.Sprintf("systemctl status %s.service", config.AppName))
	printOut(msg)

	fmt.Println("Systemd-сервис успешно создан и запущен.")
	return nil
}

func createExecRemote(client *ssh.Client) func(cmd string) (string, error) {
	return func(cmd string) (string, error) {
		session, err := client.NewSession()
		if err != nil {
			return "", fmt.Errorf("ошибка при создании сессии SSH: %v", err)
		}
		defer session.Close()

		out, err := session.CombinedOutput(cmd) // объединяет стандартный вывод и ошибки
		if err != nil {
			return string(out), fmt.Errorf("ошибка при запуске команды SSH: %v. Вывод: %s", err, string(out))
		}

		return string(out), nil
	}
}

func printOut(msg string) {
	fmt.Println(">> ", msg)
}

func Run(config *Config) error {
	file, err := buildApp(config)

	defer file.Close()

	if err != nil {
		return fmt.Errorf("Ошибка при сборке приложения: %v\n", err)
	}

	var client *ssh.Client

	for _, device := range config.Devices {
		remoteAppPath := device.AppDir + "/" + config.AppName

		client, err = createSSHClient(device.Host, device.User, device.Password)
		if err != nil {
			return fmt.Errorf("Ошибка при подключении по SSH: %v\n", err)
		}

		err = stopService(client, config)
		if err != nil {
			return err
		}

		err = uploadToDevice(client, file, remoteAppPath)
		if err != nil {
			return fmt.Errorf("Ошибка при передаче файла на устройство: %v\n", err)
		}

		err = createAndStartService(client, config, device)
		if err != nil {
			return fmt.Errorf("Ошибка при создании или запуске systemd-сервиса: %v\n", err)
		}

		client.Close()
	}

	fmt.Println("Процесс завершён успешно.")

	return nil
}

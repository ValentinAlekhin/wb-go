# Установка

## Требования

- [Go](https://go.dev/doc/install) версии 1.23 или выше.
- Терминал для работы с CLI утилитой.
- Текстовый редактор с поддержкой Go. Рекомендуется [VSCode](https://code.visualstudio.com/download) и [GoLand](https://www.jetbrains.com/help/go/installation-guide.html).

## Установка

```shell
go get github.com/ValentinAlekhin/wb-go
```

### Использование

После установки можно использовать команду для генерации.

```shell
go run github.com/ValentinAlekhin/wb-go generate -b 192.168.1.10:1883 -o ./internal/device
```

Программа подключиться к MQTT брокеру по адресу `192.168.1.10:1883` и сгенерирует файлы устройств в папку `./internal/device`

## Работа с устройствами

<<< @/../examples/usage/main.go
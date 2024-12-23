# Развертывание приложения

## Описание

Библиотека предоставляет возможность по команде разворачивать приложение на устройствах Wiren Board.

Принцип работы:
- Сборка приложения.
- Подключение по SSH к контроллеру.
- Загрузка бинарного файла на контроллер.
- Создание systemd конфигурационного файла для приложения.
- Создание сервиса в systemd.
- Запуск сервиса в systemd.
- Проверка статуса сервиса.

## Использование

Создайте файл `wb-go-deploy.yaml`.

::: code-group
<<< @/data/examples/wb-go-deploy.yaml
:::

Команда запуска.

::: warning
Запускать в директории, где находится файл `wb-go-deploy.yaml`.
:::

```shell
go run github.com/ValentinAlekhin/wb-go deploy
```

Вывод команды.

```shell
Сборка приложения для Wiren Board 7...
Сборка завершена успешно!
Файл успешно передан на устройство.
>>  ● example-app.service - example-app Service
     Loaded: loaded (/etc/systemd/system/example-app.service; enabled; vendor preset: enabled)
     Active: active (running) since Mon 2024-12-23 17:29:12 UTC; 3s ago
   Main PID: 14532 (example-app)
      Tasks: 9 (limit: 4790)
     Memory: 1.4M
        CPU: 93ms
     CGroup: /system.slice/example-app.service
             └─14532 /mnt/data//example-app

Dec 23 17:29:12 wirenboard-AQZW53KM systemd[1]: Started example-app Service.
Dec 23 17:29:12 wirenboard-AQZW53KM example-app[14532]: Подключение к брокеру 192.168.1.150:1883
Dec 23 17:29:12 wirenboard-AQZW53KM example-app[14532]: Подключение к MQTT-брокеру успешно!
Dec 23 17:29:12 wirenboard-AQZW53KM example-app[14532]: Получено новое сообщение: 15.000000
Dec 23 17:29:13 wirenboard-AQZW53KM example-app[14532]: Получено новое сообщение: 12.000000
Dec 23 17:29:14 wirenboard-AQZW53KM example-app[14532]: Получено новое сообщение: 9.000000

Systemd-сервис успешно создан и запущен.
Процесс завершён успешно.
```
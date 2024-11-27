# go-musthave-shortener-tpl

Шаблон репозитория для трека «Сервис сокращения URL».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).

shortenertest -test.v -test.run=^TestIteration6$ -binary-path=shortener -server-port=8080 -source-path=../../

shortenertest -test.v -test.run=^TestIteration9$ -binary-path=shortener -server-port=8080 -source-path=cmd/ -file-storage-path=C:\Users\User\storage_shortener.txt

shortenertest -test.v -test.run=^TestIteration10$ -binary-path=shortener -server-port=8080 -source-path=cmd/ -database-dsn="host=127.0.0.1 user=practicum password=123456 dbname=practicumdb sslmode=disable"

shortenertest -test.v -test.run=^TestIteration12$ -binary-path=shortener -server-port=8080 -source-path=cmd/ -file-storage-path=/Users/alena/storage_shortener.txt -database-dsn="host=127.0.0.1 user=practicum password=123456 dbname=practicumdb sslmode=disable"

shortenertest -test.v -test.run=^TestIteration14$ -binary-path=shortener -server-port=8080 -source-path=cmd/ -file-storage-path=/Users/alena/storage_shortener.txt -database-dsn="host=127.0.0.1 user=practicum password=123456 dbname=practicumdb sslmode=disable"

shortenertest -test.v -test.run=^TestIteration15$ -binary-path=shortener -server-port=8080 -source-path=cmd/ -file-storage-path=/Users/alena/storage_shortener.txt -database-dsn="host=127.0.0.1 user=practicum password=123456 dbname=practicumdb sslmode=disable"


mockgen -destination=internal/mocks/mock_store.go -package=mocks internal/repository Storager
--build_flags=--mod=mod

mockgen -destination=internal/api/mock_store.go -source=internal/api/api.go Storager 
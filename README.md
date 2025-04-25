# avitoTestTask
Тестовое задание на стажировку в Авито

## Как запускать проект:

Клонируем репозиторий:
```sh
git clone https://github.com/V1Ro-Dev/avitoTestTask
```

Переходим в папку deploy:
```sh
cd deploy
```

Запускаем наше приложение + Postgres в Docker-контейнерах
```sh
make up 
либо 
docker-compose up --build -d
```

Перезапуск приложения + Postgres:
```sh
make restart 
либо 
docker-compose down -v
docker-compose up --build -d [если нужен новый билд]
```

Запуск тестов и получение отчета о покрытии
```sh
cd ./backend
make summarize-coverage
или
go test ./...
```

Запуск интеграционного теста
```sh
make integration
или
cd /integration-tests
go test basic-flow_test.go
```

Генерация структур и методов из файла pvz.proto
```sh
make gen-proto
или
cd /internal/grpc/pvz
protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto
```

В субботу 19.04 была предзащита в технопарке а до нее мы допиливали проект, поэтому получился спидран за 1 сутки и почти без допок. Сегодня постараюсь докинуть grpc и мониторинг, но это по факту будет после дедлайна :(


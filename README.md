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
go te
```

Запуск интеграционного теста
```sh
make integration
или
cd /integration-tests
go test ./...
```

В субботу 19.04 была предзащита в технопарке а до нее мы допиливали проект, поэтому получился спидран за 1 сутки и почти без допок. Сегодня постараюсь докинуть grpc и мониторинг, но это по фатку будет после дедлайна :(


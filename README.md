# Cервис для сотрудников ПВЗ, который позволит вносить информацию по заказам в рамках приёмки товаров.

## Сервис поддерживает:
1. Все методы описанные в документации к API [см. docs/swagger.yaml]
2. Покрытие unit-тестами 73%
3. Логгирование всех запросов
4. gRPC метод, возвращающий все ПВЗ в системе
5. Пользовательская регистрация и авторизация
6. Интеграционный тест, который:

   a) Создает новый ПВЗ

   б) Добавляет новую приёмку заказов

   в) Добавляет 50 товаров в рамках текущей приёмки заказов

   г) Закрывает приёмку заказов

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

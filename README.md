# subscriptions-app

REST API для управления подписками пользователей. Приложение написано на Go, использует Echo, PostgreSQL, sqlc, goose-миграции и Swagger-документацию.

## Что умеет API

- создавать, получать, обновлять и удалять подписки;
- получать список подписок с фильтрацией, сортировкой и пагинацией;
- считать суммарную стоимость подписок за период.

## Технологии

- Go 1.26.4
- Echo
- PostgreSQL
- pgx
- sqlc
- goose
- swaggo / echo-swagger
- Docker Compose

## Запуск через Docker Compose

Основной способ запуска проекта:

```bash
cp .env.example.docker .env
make docker-up
```

После запуска:

- API: `http://localhost:8800`
- Swagger UI: `http://localhost:8800/swagger/index.html`
- Base path API: `http://localhost:8800/api/v1`

Docker Compose поднимает:

- `backend` - Go-приложение на порту `8800`;
- `db` - PostgreSQL;
- `migrate` - контейнер, который применяет goose-миграции перед стартом backend.

Полезные команды:

```bash
make docker-logs
make docker-down
```

## Локальный запуск без Docker Compose

Создайте `.env` из локального примера:

```bash
cp .env.example .env
```

В `.env.example` база настроена на:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5439/postgres
```

Перед локальным запуском убедитесь, что PostgreSQL доступен по этому адресу, либо измените `DATABASE_URL` и `GOOSE_DBSTRING` под свою базу.

Примените миграции, если запускаете базу не через Docker Compose:

```bash
go install github.com/pressly/goose/v3/cmd/goose@v3.27.1
set -a
source .env
set +a
goose up
```

Запуск приложения:

```bash
make run
```

Приложение слушает порт `8800`:

```text
http://localhost:8800
```

## Swagger и OpenAPI

Swagger UI доступен после запуска приложения:

```text
http://localhost:8800/swagger/index.html
```

Сгенерированные OpenAPI-файлы находятся в директории `docs/`:

- `docs/swagger.json`
- `docs/swagger.yaml`
- `docs/docs.go`

При изменении Swagger-аннотаций документацию можно пересобрать командой:

```bash
swag init -g cmd/app/main.go
```

Если `swag` не установлен:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Основные маршруты

Все бизнес-маршруты находятся под префиксом `/api/v1`.

| Метод | Путь | Описание |
| --- | --- | --- |
| `GET` | `/api/v1/sub` | список подписок |
| `GET` | `/api/v1/sub/sum` | сумма подписок за период |
| `POST` | `/api/v1/sub` | создать подписку |
| `GET` | `/api/v1/sub/{id}` | получить подписку по UUID |
| `PUT` | `/api/v1/sub` | обновить подписку |
| `DELETE` | `/api/v1/sub/{id}` | удалить подписку |

### Пример создания подписки

```bash
curl -X POST http://localhost:8800/api/v1/sub \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Netflix",
    "price": 999,
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "start_date": "01-2024",
    "end_date": "12-2024"
  }'
```

Поле `id` можно не передавать: приложение сгенерирует UUID автоматически.

### Пример подсчета суммы

```bash
curl "http://localhost:8800/api/v1/sub/sum?start_date=01-2024&end_date=12-2024"
```

## Переменные окружения

| Переменная | Назначение |
| --- | --- |
| `DATABASE_URL` | строка подключения приложения к PostgreSQL |
| `GOOSE_DRIVER` | драйвер goose, для проекта используется `postgres` |
| `GOOSE_DBSTRING` | строка подключения goose к PostgreSQL |
| `GOOSE_MIGRATION_DIR` | директория миграций |
| `GOOSE_TABLE` | таблица, в которой goose хранит состояние миграций |
| `POSTGRES_PASSWORD` | пароль пользователя PostgreSQL в Docker Compose |
| `POSTGRES_USER` | пользователь PostgreSQL в Docker Compose |
| `POSTGRES_DB` | база PostgreSQL в Docker Compose |

## Разработка

Сборка бинарника:

```bash
make build
```

Запуск собранного бинарника:

```bash
make run-bin
```

Тесты:

```bash
make test
```

Форматирование:

```bash
make fmt
```

Генерация sqlc-кода:

```bash
make gen-sqlc
```

Через Docker:

```bash
make docker-gen-sqlc
```

## Структура проекта

```text
cmd/app/                 точка входа приложения
internal/subscriptions/  handlers, service layer и repository для подписок
internal/db/             сгенерированный sqlc-код и SQL-схема
migrations/              goose-миграции
docs/                    Swagger/OpenAPI-документация
docker-compose.yaml      запуск приложения, PostgreSQL и миграций
Makefile                 основные команды разработки
```

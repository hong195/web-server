# Skinport API

API для работы с предметами Skinport и балансом пользователей.

## Стек

- Go 1.25, Fiber
- PostgreSQL
- In-memory cache
- Docker Compose

## Архитектура

Clean Architecture. Код разделён на слои:

- `controller` — HTTP handlers, валидация
- `usecase` — бизнес-логика
- `repo` — работа с БД и внешними API
- `entity` — модели данных

Каждый слой зависит только от слоя ниже. Бизнес-логика не знает про HTTP или базу данных.

## Кеширование

Items кешируются в памяти. При старте приложения запускается горутина, которая:
1. Сразу загружает данные из Skinport API
2. Обновляет кеш каждые N секунд (настраивается через `SKINPORT_CACHE_TTL_SEC`)

GET /items всегда читает из кеша. Это нужно потому что Skinport API отвечает медленно (~2-3 сек).

## Запуск

```bash
# Docker
make compose-up-all

# Локально
cp .env.example .env
make run
```

- API: http://localhost:8080
- Swagger: http://localhost:8080/api/swagger/index.html

## Конфигурация

```
PG_URL=postgres://user:password@localhost:5432/db
HTTP_PORT=8080
SWAGGER_ENABLED=true
SKINPORT_CACHE_TTL_SEC=300
```

## API

- `GET /api/v1/items?page=1&limit=100` — список предметов
- `GET /api/v1/users/:id` — получить пользователя
- `POST /api/v1/balance/deduct` — списать баланс

## Команды

```bash
make run              # запуск
make test             # тесты
make compose-up-all   # docker
make compose-down     # остановить
make migrate-up       # миграции
```

## Структура

```
cmd/app/           - точка входа
internal/
  controller/      - HTTP handlers
  usecase/         - бизнес-логика
  repo/persistent/ - PostgreSQL
  repo/webapi/     - Skinport API
  entity/          - модели
pkg/cache/         - in-memory кеш
migrations/        - SQL миграции
```

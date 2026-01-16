# Skinport API

API для работы с предметами Skinport и балансом пользователей.

## Запуск

```bash
# Весь стек через Docker
make compose-up-all

# Или локально (нужен запущенный PostgreSQL)
cp .env.example .env
make run
```

После запуска:
- API: http://localhost:8080
- Swagger: http://localhost:8080/api/swagger/index.html

## Конфигурация

Основные переменные в `.env`:

```
PG_URL=postgres://user:password@localhost:5432/db
HTTP_PORT=8080
SWAGGER_ENABLED=true
SKINPORT_CACHE_TTL_SEC=300
```

## API

### GET /api/v1/items

Список предметов с пагинацией.

```
GET /api/v1/items?page=1&limit=100
```

### GET /api/v1/users/:id

Получить пользователя.

### POST /api/v1/balance/deduct

Списать баланс.

```json
{"user_id": 1, "amount": 100.50}
```

## Команды

```bash
make run              # запуск
make test             # тесты
make compose-up-all   # docker
make compose-down     # остановить
make migrate-up       # миграции
make swag-v1          # swagger
```

## Структура

```
cmd/app/          - точка входа
internal/
  controller/     - handlers
  usecase/        - бизнес-логика
  repo/           - БД и внешние API
  entity/         - модели
migrations/       - SQL миграции
```

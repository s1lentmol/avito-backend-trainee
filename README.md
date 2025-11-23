# PR Reviewer Assigner

Микросервис для назначения ревьюеров на Pull Request'ы, управления командами и активностью пользователей. Задание выполнено по условию из `ASSIGNMENT.md` и спецификации `openapi.yaml`.

## Конфигурация

Конфигурация собирается из двух источников:

1. Встроенный YAML `config/config.yaml` - значения по умолчанию.
2. Переменные окружения с префиксом `ENV_`, которые переопределяют значения из YAML.

Структура конфигурации:

- Параметры БД:
  - `ENV_DB_USERNAME` - пользователь (по умолчанию `postgres`).
  - `ENV_DB_PASSWORD` - пароль (по умолчанию `postgres`).
  - `ENV_DB_HOST` - хост (по умолчанию `postgres` - имя сервиса в Docker Compose).
  - `ENV_DB_PORT` - порт (по умолчанию `5432`).
  - `ENV_DB_NAME` - имя БД (по умолчанию `reviewer-assigner`).
- Параметры приложения:
  - `ENV_APP_NAME` - имя приложения (по умолчанию `reviewer-assigner`).
  - `ENV_APP_PORT` - HTTP‑порт (по умолчанию `8080`).

Пример файла `.env` находится в `.env.example`. Использование `.env` **не обязательно**: при его отсутствии используются значения по умолчанию.

## Запуск через Docker Compose

Шаги:

```bash
git clone <this-repo>
cd avito-backend-trainee
docker compose up --build
```

Переопределение конфигурации:

- Можно создать `.env` рядом с `docker-compose.yaml` и задать в нём:
  - `ENV_APP_PORT`, `ENV_DB_*` и др. (см. `.env.example`).
- Либо передать переменные окружения напрямую при запуске:

```bash
ENV_APP_PORT=9090 ENV_DB_NAME=mydb docker compose up --build
```

## Команды Makefile

В корне проекта есть `Makefile` с базовыми командами:

- `make docker-build` - собрать Docker‑образ `reviewer-assigner`.
- `make compose-up` - поднять сервис и PostgreSQL в Docker.
- `make compose-down` - остановить и удалить контейнеры.
- `make test` - запустить все Go‑тесты.

## Стек и архитектура

Сервис построен по слоистой архитектуре:

- HTTP‑слой (`internal/controller/http`) — хендлеры на Fiber, валидируют запросы, маппят ошибки домена в HTTP‑коды и формат, описанный в `openapi.yaml`.
- Usecase‑слой (`internal/*/usecase`) — бизнес‑логика для пользователей, команд и PR, работает через интерфейсы репозиториев и доменных моделей.
- Доменный слой (`internal/*/domain`) — основные сущности (`User`, `Team`, `PullRequest`) и операции над ними (назначение и переназначение ревьюверов, проверка статуса `MERGED` и т.д.).
- Хранение данных — PostgreSQL, доступ реализован через репозитории в `internal/*/adapter/postgres`, миграции применяются при старте приложения через `migrator`.
- Тесты - unit‑тесты для доменного и usecase‑слоёв, написанные на `testify`.

Используемые технологии:

- Go 1.25
- HTTP‑фреймворк: `gofiber/fiber v2`
- База данных: PostgreSQL 15
- Миграции: `pressly/goose`
- Конфигурация: `spf13/viper` (+ embedded `config/config.yaml` как значения по умолчанию)
- Логирование: стандартный `log/slog`
- Тестирование: `testify` + unit‑тесты домена и usecase‑слоя

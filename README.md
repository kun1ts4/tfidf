# Анализатор TF-IDF

Это простой сервис, который помогает анализировать текстовые документы с помощью метода TF-IDF (частотность слов с учётом их уникальности в тексте)

## Структура проекта

```text
.
├── CHANGELOG.md                  # Версии и изменения
├── cmd/                          # Точки входа (API)
├── configs/                      # Конфигурации (YAML, ENV)
├── docker-compose.yaml           # Docker Compose для запуска приложения
├── Dockerfile                    # Сборка контейнера
├── docs/                         # Документация (Swagger)
├── internal/                     # Внутренняя логика приложения:
│   ├── handler/                  #   - HTTP-обработчики
│   ├── service/                  #   - бизнес-логика
│   ├── repository/               #   - работа с БД
│   └── model/                    #   - сущности БД
└── pkg/                          # Подключение к базе
```

## Запуск

Для запуска необходим валидный .env
```bash
# Склонируйте репозиторий
git clone https://github.com/kun1ts4/tfidf

# Перейдите в папку с проектом
cd tfidf

# Запустите проект с помощью Docker
docker-compose up --build -d
```

После запуска приложение будет доступно по адресу: [http://localhost:8080](http://localhost:8080)

## Конфигурация

Указывается в `.env` в корне проекта

| Параметр   | Описание                           | Пример значения |
| ---------- | ---------------------------------- | --------------- |
| DB_HOST    | Хост базы данных                   | db              |
| DB_PORT    | Порт базы данных                   | 5432            |
| DB_USER    | Имя пользователя базы данных       | user            |
| DB_PASSWORD| Пароль пользователя базы данных    | password        |
| DB_NAME    | Имя базы данных                    | postgres        |
| API_PORT   | Порт для API                       | 8080            |

## Версия

v1.2

## Зависимости

| Зависимость                  | Версия                          |
|------------------------------|---------------------------------|
| gopkg.in/yaml.v2             | v2.4.0                          |
| github.com/gin-gonic/gin      | v1.10.1                         |
| github.com/jackc/pgx/v5       | v5.7.5                          |
| github.com/golang-jwt/jwt/v5  | v5.2.2                          |
| github.com/google/uuid        | v1.6.0                          |
| github.com/lpernett/godotenv  | v0.0.0-20230527005122-0de1d4c5ef5e |
| github.com/swaggo/files       | v1.0.1                          |
| github.com/swaggo/gin-swagger | v1.6.0                          |
| github.com/swaggo/swag        | v1.16.4                         |

## Технологии

golang:1.23.0, postgres:15, Docker, Swagger

## Документация API
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Changelog
[CHANGELOG.md](CHANGELOG.md)

## Схемы таблиц БД
[DATABASE_SCHEMA.md](DATABASE_SCHEMA.md)

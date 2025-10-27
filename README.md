#  REST-service-sub

### REST-сервис для агрегации данных об онлайн-подписках пользователей

---
 
##  Описание

**REST-service-sub** — это REST API сервис, написанный на ЯП Golang.
Он позволяет хранить и агрегировать данные о **подписках пользователей** (например, Netflix, Spotify и т.д.).

**Сервис реализует:**
- создание/получение/обновление/удаление подписок пользователей;
- расчёт общей стоимости подписок за период;
- опциональную фильтрацию по пользователю и названию сервиса;
- документацию API через **Swagger UI**.

---

##  Стек/Инструментарий

- **Go** - v1.22+
- **Gin** — HTTP фреймворк
- **GORM** — ORM для работы с базой данных
- **PostgreSQL** — основная база данных
- **Docker Compose** — для сборки и запуска всего проекта
- **Swaggo** — автогенерация Swagger-документации
- **golang-migrate** — миграции БД
## Структура проекта
```csharp
REST-service-sub/
├── cmd/
│   └── api/
│       └── main.go
├── docker/
│   └── Dockerfile
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   └── postgres.go
│   ├── handler/
│   │   ├── dto.go
│   │   ├── handler.go
│   │   ├── error_response.go
│   ├── logger/
│   │   └── logger.go
│   ├── middleware/
│   │   └── middleware.go
│   ├── model/
│   │   └── model.go
│   └── service/
│       ├── service.go
├── migrations/
│   ├── 01_init_sub.down.sql
│   └── 01_init_sub.up.sql
├── .env
├── .gitignore
├── docker-compose.yml
├── example.env
├── go.mod
├── go.sum
└── README.md
```
##  Установка и запуск
### Предварительно
**_Gin,GORM, Docker, Swagger, golang migrate_** обязательно должны быть установлены, зависимости подтянуты в go.mod c актуальными версиями:
```Golang
require (
	github.com/gin-gonic/gin v1.11.0
	github.com/go-playground/validator/v10 v10.28.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/rs/zerolog v1.34.0
	github.com/stretchr/testify v1.11.1
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.1
	github.com/swaggo/swag v1.8.12
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.0
)
```
### 1. Клонируем репозиторий

```bash
git clone https://github.com/UberionAI/REST-service-sub.git
cd REST-service-sub
```

### 2. Настроить .env файл
Пример .env файла есть в основном каталоге проекта (файл _example.env_)

### 3. Запуск сервиса:
#### Запустить через **Docker Compose**

```bash
docker compose up --build
```

* **API доступен по адресу**: http://localhost:8000
* **Swagger UI**: http://localhost:8000/swagger/index.html # endpoint docs, examples and testing API

### 4. Тестирование
Запуск тестов:
```bash
go test ./... -v
```
Тесты покрывают:
* бизнес-логику (_service_test.go_)
* HTTP-обработчики (_handler_test.go_)
* негативные кейсы (ошибки формата даты, неверный user_id и т.д.)
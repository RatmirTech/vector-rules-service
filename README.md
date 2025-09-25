# Vector Rules Service

Микросервис для работы с векторной базой данных на PostgreSQL + pgvector, реализующий поиск правил по векторному сходству в чистой архитектуре (Clean Architecture).

## Возможности

- 🔍 **Векторный поиск**: Поиск правил по семантическому сходству с использованием pgvector
- 🚀 **Двойной API**: gRPC для векторного поиска, HTTP REST для CRUD операций
- 🏗️ **Чистая архитектура**: Разделение на слои domain, usecase, repository, transport
- 🐳 **Docker-ready**: Готовые Dockerfile и docker-compose.yml
- 📊 **Embedding агрегация**: Усреднение множественных запросов в единый вектор
- 🔧 **Расширяемость**: Простая замена mock embedding provider на реальный (OpenAI, Cohere, etc.)

## Архитектура

```
├── cmd/server/          # Точка входа приложения
├── internal/
│   ├── domain/          # Бизнес-логика и интерфейсы
│   ├── usecase/         # Сценарии использования
│   ├── repository/      # Репозитории для работы с БД
│   ├── transport/
│   │   ├── grpc/        # gRPC сервер (только Retrieve)
│   │   └── http/        # HTTP сервер (CRUD операции)
│   ├── infra/
│   │   ├── db/          # Подключение к PostgreSQL
│   │   └── embeddings/  # Генерация векторных представлений
│   └── config/          # Конфигурация
├── proto/               # Protocol Buffers определения
├── init-db/             # SQL миграции
└── docker-compose.yml   # Docker композиция
```

## База данных

### Таблицы

**rule_types** - категории правил:
- `id` (BIGSERIAL PK)
- `name` (TEXT UNIQUE)
- `created_at`, `updated_at` (TIMESTAMP)

**rules** - правила с векторными представлениями:
- `id` (BIGSERIAL PK) 
- `rule_type_id` (FK -> rule_types)
- `content` (JSONB) - содержимое правила в JSON
- `embedding` (vector(1536)) - векторное представление для поиска
- `created_at`, `updated_at` (TIMESTAMP)

## API

### gRPC (только векторный поиск)

**Порт**: 9090  
**Сервис**: `RuleRetrievalService`

#### Retrieve
Поиск правил по векторному сходству
```protobuf
rpc Retrieve(RetrieveRequest) returns (RetrieveResponse);
```

**Запрос**:
- `n` (int32) - количество правил для возврата
- `type` (string, optional) - фильтр по типу правила  
- `queries` ([]string) - массив текстовых запросов

**Ответ**:
- `rules` - список найденных правил с метаданными и score сходства

### HTTP REST (CRUD операции)

**Порт**: 8080  
**Base URL**: `/api/v1`

#### Rules API
- `POST /rules` - создание правила
- `GET /rules/:id` - получение правила
- `PUT /rules/:id` - обновление правила  
- `DELETE /rules/:id` - удаление правила
- `GET /rules?type=<type>&limit=<n>&offset=<n>` - список правил

#### Rule Types API  
- `POST /rule-types` - создание типа правил
- `GET /rule-types/:id` - получение типа правил
- `PUT /rule-types/:id` - обновление типа правил
- `DELETE /rule-types/:id` - удаление типа правил
- `GET /rule-types?limit=<n>&offset=<n>` - список типов правил

## Быстрый старт

### Требования
- Docker и Docker Compose
- Go 1.21+ (для разработки)
- Protocol Buffers compiler (для генерации gRPC кода)

### 1. Клонирование и запуск

```bash
git clone <repository>
cd vector-rules-service

# Запуск с Docker Compose
make docker-up

# Или запуск только БД для разработки
make db-up
```

### 2. Генерация protobuf кода

```bash
# Установка protoc (macOS)
brew install protobuf

# Установка Go плагинов
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Генерация кода
make proto
```

### 3. Локальный запуск для разработки

```bash
# Установка зависимостей
go mod tidy

# Запуск приложения
make run
```

## Примеры использования

### HTTP API

#### Создание типа правил
```bash
curl -X POST http://localhost:8080/api/v1/rule-types \
  -H "Content-Type: application/json" \
  -d '{"name": "validation"}'
```

#### Создание правила
```bash  
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "validation",
    "content": {
      "description": "Email validation rule",
      "pattern": "^[\\w\\.-]+@[\\w\\.-]+\\.[a-zA-Z]{2,}$",
      "required": true
    }
  }'
```

#### Получение списка правил
```bash
curl "http://localhost:8080/api/v1/rules?type=validation&limit=10"
```

### gRPC API

#### Поиск похожих правил (с grpcurl)

```bash
# Установка grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Векторный поиск
grpcurl -plaintext \
  -d '{
    "n": 5,
    "type": "validation",  
    "queries": ["email validation", "check email format"]
  }' \
  localhost:9090 rule.v1.RuleRetrievalService/Retrieve
```

#### Пример ответа gRPC
```json
{
  "rules": [
    {
      "id": "1",
      "type": "validation",
      "content": {
        "description": "Email validation rule",
        "pattern": "^[\\w\\.-]+@[\\w\\.-]+\\.[a-zA-Z]{2,}$"
      },
      "score": 0.95,
      "createdAt": "2024-01-15T10:30:00Z"
    }
  ]
}
```

## Конфигурация

Настройка через переменные окружения:

```bash
# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=vector_rules
DB_SSLMODE=disable

# Сервер
SERVER_HOST=0.0.0.0
HTTP_PORT=8080
GRPC_PORT=9090
```

## Makefile команды

```bash
make proto         # Генерация protobuf кода
make build         # Сборка приложения  
make run           # Запуск приложения
make test          # Запуск тестов
make docker-build  # Сборка Docker образа
make docker-up     # Запуск в Docker Compose
make docker-down   # Остановка Docker Compose
make db-up         # Запуск только PostgreSQL
make clean         # Очистка артефактов сборки
make setup         # Настройка dev окружения
```

## Разработка

### Замена Embedding Provider

По умолчанию используется mock provider. Для продакшена замените в `cmd/server/main.go`:

```go
// Замените эту строку:
embeddingProvider := embeddings.NewMockEmbeddingProvider(1536)

// На реальный провайдер, например:
embeddingProvider := embeddings.NewOpenAIProvider(apiKey, "text-embedding-ada-002")
```

Реализуйте интерфейс `domain.EmbeddingProvider` в `internal/infra/embeddings/`:

```go
type EmbeddingProvider interface {
    GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
    GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
}
```

### Добавление новых API методов

1. **HTTP**: Добавьте методы в соответствующий handler в `internal/transport/http/`
2. **gRPC**: Обновите `.proto` файл, регенерируйте код, добавьте методы в server

### Работа с векторами

Векторы создаются автоматически при создании/обновлении правил. Размерность по умолчанию: 1536 (OpenAI ada-002). Для изменения:

1. Обновите размерность в SQL (init-db/001_init.sql)
2. Измените параметр в NewMockEmbeddingProvider()
3. Пересоздайте БД

## Тестирование

### Unit тесты
```bash
make test
```

### Интеграционное тестирование  
```bash
# Запуск полного стека
make docker-up

# Проверка здоровья
curl http://localhost:8080/health

# Тестирование API
curl http://localhost:8080/api/v1/rule-types
```

## Производительность

- **Индексирование**: Используется IVFFlat индекс для векторного поиска
- **Пул соединений**: pgx connection pool для оптимальной работы с БД
- **Параллелизм**: Отдельные горутины для HTTP и gRPC серверов

## Лицензия

MIT License

## Контакты

Для вопросов и предложений создавайте Issues в репозитории.
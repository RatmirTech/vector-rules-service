# Примеры запросов к Vector Rules Service

## Настройка окружения

```bash
# Базовые URL для тестирования
HTTP_BASE="http://localhost:8080/api/v1"
GRPC_HOST="localhost:9090"

# Установите grpcurl для gRPC запросов
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

## HTTP API Примеры

### Проверка здоровья сервиса
```bash
curl http://localhost:8080/health
```

### Работа с типами правил

#### Создание типов правил
```bash
# Тип валидации
curl -X POST $HTTP_BASE/rule-types \
  -H "Content-Type: application/json" \
  -d '{"name": "validation"}'

# Тип трансформации  
curl -X POST $HTTP_BASE/rule-types \
  -H "Content-Type: application/json" \
  -d '{"name": "transformation"}'

# Тип фильтрации
curl -X POST $HTTP_BASE/rule-types \
  -H "Content-Type: application/json" \
  -d '{"name": "filtering"}'

# Бизнес-логика
curl -X POST $HTTP_BASE/rule-types \
  -H "Content-Type: application/json" \
  -d '{"name": "business_logic"}'
```

#### Получение списка типов правил
```bash
curl $HTTP_BASE/rule-types
```

#### Получение конкретного типа правил
```bash
curl $HTTP_BASE/rule-types/1
```

### Работа с правилами

#### Создание правил валидации
```bash
# Email валидация
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "validation",
    "content": {
      "description": "Validate email format",
      "field": "email",
      "pattern": "^[\\w\\.-]+@[\\w\\.-]+\\.[a-zA-Z]{2,}$",
      "required": true,
      "error_message": "Invalid email format"
    }
  }'

# Телефон валидация
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "validation", 
    "content": {
      "description": "Validate phone number format",
      "field": "phone",
      "pattern": "^\\+?[1-9]\\d{1,14}$",
      "required": false,
      "error_message": "Invalid phone number format"
    }
  }'

# Возраст валидация
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "validation",
    "content": {
      "description": "Validate user age",
      "field": "age", 
      "min_value": 18,
      "max_value": 120,
      "required": true,
      "error_message": "Age must be between 18 and 120"
    }
  }'
```

#### Создание правил трансформации
```bash
# Нормализация имени
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "transformation",
    "content": {
      "description": "Normalize user name to title case",
      "field": "name",
      "operation": "title_case",
      "trim_whitespace": true,
      "remove_extra_spaces": true
    }
  }'

# Форматирование даты
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "transformation",
    "content": {
      "description": "Convert date to ISO format",
      "field": "birth_date",
      "from_format": "DD/MM/YYYY",
      "to_format": "YYYY-MM-DD",
      "timezone": "UTC"
    }
  }'
```

#### Создание правил фильтрации
```bash
# Фильтрация контента
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "filtering",
    "content": {
      "description": "Filter inappropriate content",
      "field": "message",
      "blacklist": ["spam", "offensive", "inappropriate"],
      "action": "block",
      "severity": "high"
    }
  }'

# Географическая фильтрация
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "filtering", 
    "content": {
      "description": "Geographic access filtering",
      "field": "country_code",
      "allowed_countries": ["US", "EU", "CA"],
      "action": "allow",
      "fallback": "deny"
    }
  }'
```

#### Создание бизнес-правил
```bash
# Правило скидки
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "business_logic",
    "content": {
      "description": "Volume discount calculation", 
      "conditions": {
        "min_amount": 1000,
        "customer_type": "premium"
      },
      "actions": {
        "discount_percent": 15,
        "free_shipping": true
      },
      "priority": 1
    }
  }'

# Правило лимитов
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "business_logic",
    "content": {
      "description": "Daily transaction limit",
      "conditions": {
        "account_type": "basic"
      },
      "limits": {
        "max_transactions": 10,
        "max_amount": 5000,
        "period": "daily"
      },
      "actions": {
        "block_excess": true,
        "notify_user": true
      }
    }
  }'
```

#### Получение правил
```bash
# Все правила
curl $HTTP_BASE/rules

# Правила по типу
curl "$HTTP_BASE/rules?type=validation&limit=5"

# Конкретное правило  
curl $HTTP_BASE/rules/1

# Правила с пагинацией
curl "$HTTP_BASE/rules?limit=3&offset=0"
```

#### Обновление правила
```bash
curl -X PUT $HTTP_BASE/rules/1 \
  -H "Content-Type: application/json" \
  -d '{
    "type": "validation",
    "content": {
      "description": "Enhanced email validation with domain check",
      "field": "email", 
      "pattern": "^[\\w\\.-]+@[\\w\\.-]+\\.[a-zA-Z]{2,}$",
      "required": true,
      "domain_blacklist": ["tempmail.com", "10minutemail.com"],
      "error_message": "Invalid email format or blocked domain"
    }
  }'
```

#### Удаление правила
```bash
curl -X DELETE $HTTP_BASE/rules/1
```

## gRPC API Примеры

### Векторный поиск правил

#### Поиск валидационных правил
```bash
grpcurl -plaintext \
  -d '{
    "n": 3,
    "type": "validation",
    "queries": ["email validation", "check email format", "validate email address"]
  }' \
  $GRPC_HOST rule.v1.RuleRetrievalService/Retrieve
```

#### Поиск правил трансформации
```bash
grpcurl -plaintext \
  -d '{
    "n": 2,
    "type": "transformation", 
    "queries": ["normalize name", "format user name", "clean text data"]
  }' \
  $GRPC_HOST rule.v1.RuleRetrievalService/Retrieve
```

#### Поиск бизнес-правил
```bash
grpcurl -plaintext \
  -d '{
    "n": 5,
    "queries": ["discount calculation", "pricing rules", "customer benefits"]
  }' \
  $GRPC_HOST rule.v1.RuleRetrievalService/Retrieve
```

#### Общий поиск без фильтра по типу  
```bash
grpcurl -plaintext \
  -d '{
    "n": 10,
    "queries": ["user data processing", "content moderation"]
  }' \
  $GRPC_HOST rule.v1.RuleRetrievalService/Retrieve
```

### Пример ответа gRPC
```json
{
  "rules": [
    {
      "id": "1",
      "type": "validation",
      "content": {
        "description": "Validate email format",
        "field": "email",
        "pattern": "^[\\w\\.-]+@[\\w\\.-]+\\.[a-zA-Z]{2,}$",
        "required": true,
        "error_message": "Invalid email format"
      },
      "score": 0.8945612,
      "createdAt": "2024-01-15T10:30:00Z",
      "updatedAt": "2024-01-15T10:30:00Z"
    },
    {
      "id": "2", 
      "type": "validation",
      "content": {
        "description": "Validate phone number format",
        "field": "phone",
        "pattern": "^\\+?[1-9]\\d{1,14}$",
        "required": false
      },
      "score": 0.7234891,
      "createdAt": "2024-01-15T10:35:00Z", 
      "updatedAt": "2024-01-15T10:35:00Z"
    }
  ]
}
```

## Комплексные сценарии

### Настройка системы валидации
```bash
# 1. Создать тип правил
curl -X POST $HTTP_BASE/rule-types \
  -H "Content-Type: application/json" \
  -d '{"name": "user_validation"}'

# 2. Добавить правила валидации пользователя
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user_validation",
    "content": {
      "description": "Complete user profile validation",
      "fields": {
        "email": {"pattern": "^[\\w\\.-]+@[\\w\\.-]+\\.[a-zA-Z]{2,}$", "required": true},
        "age": {"min": 18, "max": 120, "required": true}, 
        "phone": {"pattern": "^\\+?[1-9]\\d{1,14}$", "required": false}
      },
      "cross_validation": {
        "email_age_consistency": true,
        "phone_country_match": true
      }
    }
  }'

# 3. Найти похожие правила валидации
grpcurl -plaintext \
  -d '{
    "n": 5,
    "type": "user_validation",
    "queries": ["user profile validation", "validate user data", "check user information"]
  }' \
  $GRPC_HOST rule.v1.RuleRetrievalService/Retrieve
```

### A/B тестирование правил
```bash
# Версия A правила
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "business_logic",
    "content": {
      "description": "Discount rule version A",
      "version": "A",
      "conditions": {"min_amount": 500},
      "discount_percent": 10
    }
  }'

# Версия B правила  
curl -X POST $HTTP_BASE/rules \
  -H "Content-Type: application/json" \
  -d '{
    "type": "business_logic", 
    "content": {
      "description": "Discount rule version B",
      "version": "B",
      "conditions": {"min_amount": 300},
      "discount_percent": 8
    }
  }'

# Поиск правил скидок
grpcurl -plaintext \
  -d '{
    "n": 2,
    "queries": ["discount rules", "pricing strategy"]
  }' \
  $GRPC_HOST rule.v1.RuleRetrievalService/Retrieve
```

## Отладка и мониторинг

### Проверка состояния сервиса
```bash
# Здоровье HTTP API
curl -w "HTTP Status: %{http_code}\nResponse Time: %{time_total}s\n" \
  http://localhost:8080/health

# Тест gRPC connectivity  
grpcurl -plaintext $GRPC_HOST list

# Статистика по правилам
curl "$HTTP_BASE/rules?limit=1" | jq '.rules | length'
curl "$HTTP_BASE/rule-types?limit=100" | jq '.rule_types | length'
```

### Производительность векторного поиска
```bash
# Бенчмарк поиска
time grpcurl -plaintext \
  -d '{
    "n": 50,
    "queries": ["comprehensive rule search test query with multiple terms"]
  }' \
  $GRPC_HOST rule.v1.RuleRetrievalService/Retrieve
```
# Delayed notifier API

Простое API для отправки отложенных уведомлений через RabbitMQ с плагином rabbitmq-delayed-message-exchange.

## Быстрый старт

### Предварительные требования

- Docker

- Go

### Запуск проекта

#### 1. Запустите все сервисы:

```bash
docker run --name delayed_redis -d -p 6379:6379 redis:8.4
docker run --name delayed_rabbitmq -d -p 15672:5672 heidiks/rabbitmq-delayed-message-exchange:4.2.0-management
```

#### 2. Клонируйте репозиторий и запустите:

```bash
git clone github.com/xGerfy/delayed-notifications
cd delayed-notifications
go run ./cmd/main.go
```

#### 3. API будет доступно по адресу: http://localhost:8080

## API Endpoints

### -POST

http://localhost:8080/notify — создание уведомлений с датой и временем отправки

### -GET

http://localhost:8080/notify/{id} — получение статуса уведомления. В поле ID нужно вставить id созданного уведомления

### -DELETE

http://localhost:8080/notify/{id} — отмена запланированного уведомления. В поле ID нужно вставить id созданного уведомления

## Примеры использования

### Создание уведомления

#### В поле "2026-02-20T15:30:00.0+04:00" - нужно указать время когда нужно отправить уведомление!

```bash
curl -X POST http://localhost:8080/notify \
  -H "Content-Type: application/json" \
  -d '{"message": "test message", "send_at": "2026-02-20T15:30:00.0+04:00"}'
```

### Получение уведомления

#### В поле {id} - нужно указать id уведомления!

```bash
curl -X GET http://localhost:8080/notify/{id} \
  -H "Content-Type: application/json"
```

### Удаление уведомления

```bash
curl -X DELETE http://localhost:8080/notify/{id}
```

## Структура проекта

```bash
delayed-notifications/
├── cmd/
│    └── main.go           # Точка входа приложения
├── internal/
│    ├── config/           # Конфигурация
│    ├── entities/         # Модели данных
│    ├── rabbit/           # Слой работы с брокером сообщений (rabbitmq)
│    │     ├── consumer    # Обработчик входящих сообщений
│    │     └── publisher   # Обработчик исходящих сообщений
│    ├── repo/             # Слой работы с БД (redis)
│    ├── router/           # Роутер
│    │     └── handler     # HTTP обработчики
│    └── service/          # Бизнес-логика
├── go.mod                 # Зависимости Go
└── README.md              # Этот файл
```

# PZ18: gRPC микросервисы

Студент: Анастасиади Д.Е. Группа: ПИМО-01-25

Подробная документация API: [docs/pz18_api.md](docs/pz18_api.md)

Диаграмма последовательности: [docs/pz18_diagram.md](docs/pz18_diagram.md)

## Архитектура

- **Auth Service**: Предоставляет gRPC-сервер для метода `Verify`, а также HTTP для совместимости.
- **Tasks Service**: Использует gRPC-клиент для проверки токенов с дедлайном.
## Структура проекта

```
PZ18/
  proto/
    auth.proto
    authpb/
      auth.pb.go
  services/
    auth/
      cmd/auth/main.go
      internal/grpcserver/server.go
      internal/http/handlers.go
    tasks/
      cmd/tasks/main.go
      internal/client/authclient/authclient.go
      internal/http/handlers.go
      internal/http/middleware.go
      internal/service/tasks.go
  shared/
    middleware/
      requestid.go
      logging.go
  docs/
    pz18_api.md
    pz18_diagram.md
  misc/
  README.md
  go.mod
```

## Установка зависимостей

```bash
go get github.com/golang/protobuf@v1.5.4 google.golang.org/grpc@v1.79.3 google.golang.org/protobuf@v1.36.10
go mod tidy
```

## Генерация кода на основании proto-файла

```bash
protoc --go_out=. --go-grpc_out=. proto/auth.proto
```

## Запуск сервисов

### Запуск Auth Service

```bash
cd services/auth
$env:AUTH_GRPC_PORT = "50051"
$env:AUTH_PORT = "8081"
go run ./cmd/auth
```

### Запуск Tasks Service

```bash
cd services/tasks
$env:TASKS_PORT = "8082"
$env:AUTH_GRPC_ADDR = "localhost:50051"
go run ./cmd/tasks
```

## Быстрая проверка работоспособности

### 1. Получить токен (HTTP)

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8081/v1/auth/login" `
  -ContentType "application/json" `
  -Body '{"username":"student","password":"student"}'
```

Ожидаемый ответ: `{"access_token":"demo-token","token_type":"Bearer"}`
![](misc/1.%20Получить%20токен.png)

### 2. Проверить токен через HTTP (Auth)

```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8081/v1/auth/verify" `
  -Headers @{ "Authorization" = "Bearer demo-token" }
```

Ожидаемый ответ: `{"valid":true,"subject":"student"}`
![](misc/2.%20Проверить%20токен%20через%20HTTP.png)

### 3. Создать задачу (Tasks с gRPC verify)

```powershell
Invoke-RestMethod -Method Post -Uri "http://localhost:8082/v1/tasks" `
  -ContentType "application/json" `
  -Headers @{ "Authorization" = "Bearer demo-token" } `
  -Body '{"title":"Test gRPC","description":"Verify via gRPC","due_date":"2026-03-20"}'
```

Ожидаемый ответ: JSON с задачей, ID: "t_1"
![](misc/3.%20Создать%20задачу.png)

### 4. Получить список задач

```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8082/v1/tasks" `
  -Headers @{ "Authorization" = "Bearer demo-token" }
```

Ожидаемый ответ: Массив задач, включая созданную.
![](misc/4.%20Получить%20список%20задач.png)

### 5. Тест недоступности Auth (остановить Auth сервис)

Остановите Auth сервис (Ctrl+C), затем:

```powershell
try {
  Invoke-RestMethod -Method Get -Uri "http://localhost:8082/v1/tasks" `
    -Headers @{ "Authorization" = "Bearer demo-token" }
} catch {
  $_.Exception.Response.StatusCode
}
```

Ожидаемый ответ: 503 (Service Unavailable)
![](misc/5.%20Тест%20недоступности%20Auth.png)

### 6. Тест невалидного токена

```powershell
try {
  Invoke-RestMethod -Method Get -Uri "http://localhost:8082/v1/tasks" `
    -Headers @{ "Authorization" = "Bearer bad-token" }
} catch {
  $_.Exception.Response.StatusCode
}
```

Ожидаемый ответ: 401 (Unauthorized)
![](misc/6.%20Тест%20невалидного%20токена.png)

### 7. Логи в терминалах

- Auth: Логи запуска gRPC и HTTP серверов.
- Tasks: Логи "calling grpc verify" при проверке токена.

![](misc/7.%20Логи.png)

## Описание ошибок и маппинг на HTTP статусы

### Auth Service (gRPC)

- **Invalid token**: gRPC код `codes.Unauthenticated` (пустой или не "demo-token").
- **Internal errors**: gRPC код `codes.Internal` (неожиданные ошибки).

### Tasks Service (HTTP API)

При проверке токена через gRPC, ошибки мапятся на HTTP статусы:

- **gRPC `codes.Unauthenticated`**: HTTP 401 Unauthorized (неверный токен).
- **gRPC `codes.Unavailable` или `codes.DeadlineExceeded`**: HTTP 503 Service Unavailable (Auth сервис недоступен, timeout).
- **Другие gRPC ошибки**: HTTP 502 Bad Gateway (ошибки связи или другие проблемы).

Если токен отсутствует в заголовке Authorization, сразу 401.

## Границы

- Auth: gRPC `Verify` метод, HTTP для логина/verify.
- Tasks: HTTP API для задач, gRPC клиент для auth с дедлайном.
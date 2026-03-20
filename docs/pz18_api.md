# Документация API для PZ18

## gRPC Контракт

### proto/auth.proto

```proto
syntax = "proto3";
package auth;

option go_package = "proto/authpb";

service AuthService {
  rpc Verify(VerifyRequest) returns (VerifyResponse);
}

message VerifyRequest {
  string token = 1;
}

message VerifyResponse {
  bool valid = 1;
  string subject = 2;
}
```

## Auth Service

### gRPC: AuthService.Verify

**Запрос:**

```json
{
  "token": "demo-token"
}
```

**Ответ (успех):**

```json
{
  "valid": true,
  "subject": "student"
}
```

**Ошибки:**

- `codes.Unauthenticated`: Неверный токен
- `codes.Internal`: Внутренние ошибки

### HTTP: POST /v1/auth/login (совместимость)

**Запрос:**

```json
{
  "username": "student",
  "password": "student"
}
```

**Ответ 200:**

```json
{
  "access_token": "demo-token",
  "token_type": "Bearer"
}
```

### HTTP: GET /v1/auth/verify (совместимость)

**Заголовки:**

- Authorization: Bearer <token>

**Ответ 200:**

```json
{
  "valid": true,
  "subject": "student"
}
```

**Ответ 401:**

```json
{
  "valid": false,
  "error": "unauthorized"
}
```

## Tasks Service

Все HTTP эндпоинты требуют заголовок Authorization. Проверка происходит через gRPC с дедлайном 2 секунды.

### POST /v1/tasks

Создать задачу.

**Запрос:**

```json
{
  "title": "Test gRPC",
  "description": "Verify via gRPC",
  "due_date": "2026-03-20"
}
```

**Ответ 201:**

```json
{
  "id": "t_1",
  "title": "Test gRPC",
  "description": "Verify via gRPC",
  "due_date": "2026-03-20",
  "done": false
}
```

**Ошибки:**

- 401: Неверный токен
- 503: Auth сервис недоступен
- 502: Другие ошибки gRPC

### GET /v1/tasks

Список задач.

**Ответ 200:**

```json
[
  {"id":"t_1","title":"Test gRPC","done":false}
]
```

### GET /v1/tasks/{id}

Получить задачу.

**Ответ 200:**

```json
{
  "id": "t_1",
  "title": "Test gRPC",
  "done": false
}
```

### PATCH /v1/tasks/{id}

Обновить задачу.

**Запрос:**

```json
{
  "title": "Updated",
  "done": true
}
```

**Ответ 200:** Обновленная задача.

### DELETE /v1/tasks/{id}

Удалить задачу.

**Ответ 204:** Успешно удалено.

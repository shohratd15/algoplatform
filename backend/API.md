# API Documentation

## Base URL

```bash
http://localhost:8080
```

## 1️⃣ Health Check

| Endpoint | Method | Auth | Description                        |
| -------- | ------ | ---- | ---------------------------------- |
| `/ping`  | GET    | No   | Проверка работоспособности сервиса |

**Response (200 OK):**

```json
"Service is running! DB connection successful."
```

## 2️⃣ Users

### 2.1 Register

| Endpoint    | Method | Auth | Description              |
| ----------- | ------ | ---- | ------------------------ |
| `/register` | POST   | No   | Регистрация пользователя |

**Body:**

```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response (201 Created):** Пустое тело

**Errors:**

- 400 Bad Request — неверные данные
- 500 Internal Server Error — ошибка сервера

### 2.2 Login

| Endpoint | Method | Auth | Description        |
| -------- | ------ | ---- | ------------------ |
| `/login` | POST   | No   | Логин пользователя |

**Body:**

```json
{
  "email": "string",
  "password": "string"
}
```

**Response (200 OK):**

```json
{
  "token": "jwt_token",
  "refresh_token": "jwt_refresh_token",
  "role": "admin|user"
}
```

### 2.3 Refresh Token

| Endpoint       | Method | Auth | Description              |
| -------------- | ------ | ---- | ------------------------ |
| `/api/refresh` | POST   | No   | Обновление access токена |

**Body:**

```json
{
  "refresh_token": "jwt_refresh_token"
}
```

**Response (200 OK):**

```json
{
  "token": "new_access_token",
  "refresh_token": "new_refresh_token",
  "role": "admin|user"
}
```

**Errors:**

- 401 Unauthorized — неверный логин/пароль
- 500 Internal Server Error — ошибка генерации токена

## 3️⃣ Problems

> Protected endpoints: require JWT token (RequireUser), POST & DELETE requires RequireAdmin

### 3.1 Create Problem

| Endpoint              | Method | Auth  | Description           |
| --------------------- | ------ | ----- | --------------------- |
| `/api/admin/problems` | POST   | Admin | Создание новой задачи |

**Body:**

```json
{
  "slug": "two-sum",
  "difficulty": "easy",
  "statements": [
    {
      "language": "ru",
      "title": "Сумма",
      "statement": "Найти сумму..."
    },
    {
      "language": "tm",
      "title": "Iki sanyn jemi",
      "statement": "Iki sanyn jemini...."
    },
    {
      "language": "en",
      "title": "Two Sum",
      "statement": "Find two numbers..."
    }
  ],
  "tests": [
    {
      "id": 1,
      "input_data": "1 2",
      "expected_output": "3",
      "is_sample": true
    },
    {
      "id": 2,
      "input_data": "5 6",
      "expected_output": "11",
      "is_sample": true
    },
    ....
  ]
}
```

**Response (201 Created):** Пустое тело

### 3.2 List Problems

| Endpoint        | Method | Auth | Description            |
| --------------- | ------ | ---- | ---------------------- |
| `/api/problems` | GET    | Yes  | Получение списка задач |

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "slug": "two-sum",
    "difficulty": "easy",
    "created_at": "2025-10-09T00:00:00Z"
  }
]
```

### 3.3 Get Problem Details

| Endpoint               | Method | Auth | Description                    |
| ---------------------- | ------ | ---- | ------------------------------ |
| `/api/problems/detail` | GET    | Yes  | Получение деталей задачи по id |

**Query params:** id — идентификатор задачи

**Response (200 OK):**

```json
{
  "problem": {
    "id": 1,
    "slug": "two-sum",
    "difficulty": "easy"
  },
  "statements": [
    {
      "language": "eng",
      "title": "Two Sum",
      "statement": "Find two numbers..."
    }
  ],
  "tests": [
    {
      "id": 1,
      "input_data": "1 2 3",
      "expected_output": "3",
      "is_sample": true
    }
  ]
}
```

### 3.4 Delete Problem

| Endpoint              | Method | Auth  | Description           |
| --------------------- | ------ | ----- | --------------------- |
| `/api/admin/problems` | DELETE | Admin | Удаление задачи по id |

**Query params:** id — идентификатор задачи

**Response (204 No Content)**
.

## 4️⃣ Submissions

> Protected endpoints: require JWT token (RequireUser)

### 4.1 Create Submission

| Endpoint           | Method | Auth | Description             |
| ------------------ | ------ | ---- | ----------------------- |
| `/api/submissions` | POST   | Yes  | Создание новой отправки |

**Body:**

```json
{
  "user_id": 1,
  "problem_id": 1,
  "language_id": 1,
  "source_code": "print(sum([1,2,3]))"
}
```

**Response (201 Created):**

```json
{
  "id": 123
}
```

### 4.2 Get Submission

| Endpoint           | Method | Auth | Description              |
| ------------------ | ------ | ---- | ------------------------ |
| `/api/submissions` | GET    | Yes  | Получение отправки по id |

**Query params:** id — идентификатор отправки

**Response (200 OK):**

```json
{
  "id": 123,
  "user_id": 1,
  "problem_id": 1,
  "language_id": 1,
  "source_code": "print(sum([1,2,3]))",
  "status": "queued",
  "created_at": "2025-10-09T00:00:00Z",
  "updated_at": "2025-10-09T00:00:00Z"
}
```

## 5️⃣ Errors

| Status Code | Description                                |
| ----------- | ------------------------------------------ |
| 400         | Bad Request — неверные данные              |
| 401         | Unauthorized — не авторизован              |
| 403         | Forbidden — доступ запрещён (только Admin) |
| 404         | Not Found — объект не найден               |
| 500         | Internal Server Error — ошибка сервера     |

## 6️⃣ Summary Table of Endpoints

| Endpoint               | Method | Auth  | Description          |
| ---------------------- | ------ | ----- | -------------------- |
| `/ping`                | GET    | No    | Health check         |
| `/api/register`        | POST   | No    | Register user        |
| `/api/login`           | POST   | No    | Login user           |
| `/api/refresh`         | POST   | No    | Refresh tokens       |
| `/api/admin/problems`  | POST   | Admin | Create problem       |
| `/api/admin/problems`  | PUT    | Admin | Update problem       |
| `/api/problems`        | GET    | Yes   | List problems        |
| `/api/problems/detail` | GET    | Yes   | Get problem details  |
| `/api/admin/problems`  | DELETE | Admin | Delete problem       |
| `/api/submissions`     | POST   | Yes   | Create submission    |
| `/api/submissions`     | GET    | Yes   | Get submission by id |

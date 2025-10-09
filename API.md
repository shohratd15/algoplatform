# AlgoPlatform API Documentation

## 📘 Overview

This document describes the REST API endpoints for **AlgoPlatform** — an educational platform for algorithm practice and automatic solution judging.

---

## 🔧 General Information

**Base URL:**

```
http://localhost:8080/api/v1
```

**Content Type:**

```
application/json
```

**Authentication:**  
All protected routes require a JWT token in the header:

```
Authorization: Bearer <your_token>
```

---

## 🔐 Authentication

### 1. Register a New User

**POST** `/auth/register`

#### Request Body

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword"
}
```

#### Response (201)

```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "created_at": "2025-10-08T12:00:00Z"
}
```

#### Error (400)

```json
{
  "error": "email already exists"
}
```

---

### 2. Login

**POST** `/auth/login`

#### Request Body

```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```

#### Response (200)

```json
{
  "token": "<jwt_token>",
  "expires_in": 3600
}
```

#### Error (401)

```json
{
  "error": "invalid email or password"
}
```

---

## 📘 Problems

### 1. Get All Problems

**GET** `/problems`

#### Response (200)

```json
[
  {
    "id": 1,
    "title": "Two Sum",
    "difficulty": "Easy",
    "description": "Find two numbers that sum up to a target.",
    "created_at": "2025-10-08T12:00:00Z"
  }
]
```

---

### 2. Get Problem by ID

**GET** `/problems/{id}`

#### Example

```
GET /problems/1
```

#### Response (200)

```json
{
  "id": 1,
  "title": "Two Sum",
  "difficulty": "Easy",
  "description": "Find two numbers that sum up to a target.",
  "input_format": "n, target, array",
  "output_format": "array of indices"
}
```

#### Error (404)

```json
{
  "error": "problem not found"
}
```

---

### 3. Create Problem _(Admin only)_

**POST** `/problems`

#### Request Body

```json
{
  "title": "Fibonacci Sequence",
  "difficulty": "Medium",
  "description": "Find n-th Fibonacci number.",
  "input_format": "integer n",
  "output_format": "integer"
}
```

#### Response (201)

```json
{
  "id": 5,
  "title": "Fibonacci Sequence",
  "difficulty": "Medium"
}
```

---

### 4. Update Problem _(Admin only)_

**PUT** `/problems/{id}`

#### Request Body

```json
{
  "difficulty": "Hard",
  "description": "Find n-th Fibonacci number efficiently using dynamic programming."
}
```

#### Response (200)

```json
{
  "message": "problem updated successfully"
}
```

---

### 5. Delete Problem _(Admin only)_

**DELETE** `/problems/{id}`

#### Response (200)

```json
{
  "message": "problem deleted"
}
```

---

## 🧩 Submissions

### 1. Submit Solution

**POST** `/submissions`

#### Request Body

```json
{
  "problem_id": 1,
  "language": "python",
  "code": "def solve(): print('Hello World')"
}
```

#### Response (202)

```json
{
  "submission_id": 101,
  "status": "pending"
}
```

---

### 2. Get Submission Status

**GET** `/submissions/{id}`

#### Example

```
GET /submissions/101
```

#### Response (200)

```json
{
  "id": 101,
  "problem_id": 1,
  "status": "accepted",
  "runtime": 0.123,
  "memory": 2048,
  "created_at": "2025-10-08T12:30:00Z"
}
```

---

### 3. Get All Submissions for User

**GET** `/submissions`

#### Response (200)

```json
[
  {
    "id": 101,
    "problem_id": 1,
    "status": "accepted",
    "runtime": 0.123
  },
  {
    "id": 102,
    "problem_id": 2,
    "status": "wrong answer"
  }
]
```

---

## ⚠️ Error Responses

| HTTP Code | Description           | Example                                   |
| --------- | --------------------- | ----------------------------------------- |
| 400       | Bad Request           | `{ "error": "invalid input" }`            |
| 401       | Unauthorized          | `{ "error": "missing or invalid token" }` |
| 403       | Forbidden             | `{ "error": "permission denied" }`        |
| 404       | Not Found             | `{ "error": "resource not found" }`       |
| 500       | Internal Server Error | `{ "error": "internal server error" }`    |

---

## 🧠 Notes

- All timestamps are returned in **UTC ISO8601** format.
- Submissions are processed asynchronously by the **Judge Worker** service.
- Problem test cases are stored in the database (`problems_tests` table).

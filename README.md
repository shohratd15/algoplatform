# 🧠 AlgoPlatform

## 📘 Overview

**AlgoPlatform** is an educational platform designed to help users learn algorithms through interactive coding challenges and automated evaluation.
It allows students to solve problems, submit code, and receive instant feedback with detailed results.

## 🚀 Features

- User authentication (registration & login)
- Problem management (create, list, update, delete)
- Code submission and judging system
- Asynchronous worker for solution evaluation
- PostgreSQL persistence layer
- Configurable environment variables
- Dockerized development setup
- Clean architecture with domain-driven design principles

## 🧩 Tech Stack

| Layer            | Technology             |
| ---------------- | ---------------------- |
| Language         | Go 1.23+               |
| Framework        | net/http               |
| Database         | PostgreSQL             |
| Logger           | Zap                    |
| Containerization | Docker, Docker Compose |
| Build Tools      | Makefile               |
| Configuration    | godotenv (.env)        |

## 📂 Project Structure

```bash
algoplatform/
├── cmd/
│ └── main.go                   # Application entry point
├── internal/
│ ├── config/                   # Environment and configuration setup
│ ├── controller/http/          # HTTP routing, middleware, and handlers
│ │ ├── handlers/
│ │ │ ├── auth_handlers.go
│ │ │ ├── problem_handlers.go
│ │ │ └── submission_handlers.go
│ │ ├── middleware.go
│ │ └── router.go
│ ├── domain/                   # Domain models and interfaces
│ ├── errors/                   # Centralized error definitions
│ ├── repo/postgres/            # PostgreSQL repository implementations
│ ├── usecase/                  # Core business logic
│ ├── server/                   # HTTP server initialization
│ └── worker/                   # Asynchronous submission judging
├── init.sql                    # Database initialization script
├── Dockerfile                  # Docker image definition
├── docker-compose.yml          # Compose file (API + PostgreSQL)
├── Makefile                    # Simplified CLI commands
├── .env                        # Environment variables
├── go.mod / go.sum             # Dependencies
├── API.md                      # API Documentation
└── README.md                   # Project documentation
```

## ⚙️ Configuration

Create an .env file in the root directory:

```env
APP_PORT=8080
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=algoplatform
JWT_SECRET=your_secret_key
```

## 🐳 Docker Setup

Run the entire environment using Docker:

```bash
# Build and start containers
docker-compose up --build
```

This will start:

- app — Go API service
- db — PostgreSQL database instance

## 🛠️ Manual Setup

```bash
# 1. Clone the repository
git clone https://github.com/shohratd15/algoplatform.git
cd algoplatform

# 2. Install dependencies
go mod tidy

# 3. Set environment variables
cp .env.example .env

# 4. Run PostgreSQL (if not using Docker)
psql -U postgres -f init.sql

# 5. Run the server
go run cmd/main.go

```

## 🧾 Makefile Commands

| Command            | Description             |
| ------------------ | ----------------------- |
| `make build`       | Build the Go binary     |
| `make run`         | Run the application     |
| `make test`        | Run unit tests          |
| `make docker-up`   | Start Docker containers |
| `make docker-down` | Stop Docker containers  |

## 🧪 API Documentation

Full API documentation is available in [API.md](API.md), including:

- Authentication endpoints
- Problem management
- Submissions workflow

## 🧑‍💻 Author

Developed by **Shohrat Dovletmuradov**
Educational project for algorithm learning and backend development with Go.

# Этап сборки
FROM golang:1.24.5 AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарь
RUN go build -o algoplatform ./cmd/

# Этап запуска
FROM ubuntu:22.04

# Устанавливаем необходимые пакеты без интерактивных вопросов
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    netcat \
    python3 \
    golang-go \
    nodejs \
    npm \
    g++ \
    openjdk-11-jdk \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Копируем бинарь
COPY --from=builder /app/algoplatform .

# Копируем скрипт ожидания
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

# Экспонируем порт (берём из ENV APP_PORT)
ENV SERVER_PORT=8080
EXPOSE ${SERVER_PORT}

# Запуск приложения через скрипт ожидания
CMD ["./docker-entrypoint.sh"]
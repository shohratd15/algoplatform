# Этап сборки
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарь
RUN go build -o algoplatform ./cmd/

# Этап запуска
FROM alpine:3.18

# Добавляем netcat для проверки доступности порта
RUN apk add --no-cache ca-certificates tzdata netcat-openbsd

WORKDIR /root/

# Копируем бинарь
COPY --from=builder /app/algoplatform .

# Копируем .env если он нужен внутри контейнера
COPY .env .

# Копируем скрипт ожидания
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

# Экспонируем порт (берём из ENV APP_PORT)
ENV SERVER_PORT=8080
EXPOSE ${SERVER_PORT}

# Запуск приложения через скрипт ожидания
CMD ["./docker-entrypoint.sh"]


#!/bin/sh

# Скрипт ожидания, чтобы предотвратить race condition при запуске БД.

# Используем переменные окружения, установленные в docker-compose.yml и .env
DB_HOST=${DB_HOST:-algoplatform_db}
DB_PORT=${DB_PORT:-5432}
TIMEOUT=30

echo "Waiting for PostgreSQL ($DB_HOST:$DB_PORT) to be available..."

# Используем netcat для проверки доступности порта.
# Проверяем в цикле каждую секунду в течение 30 секунд.
for i in $(seq 1 $TIMEOUT); do
  # nc -z проверяет, слушает ли порт (без отправки данных)
  nc -z "$DB_HOST" "$DB_PORT" > /dev/null 2>&1
  
  # Проверяем код завершения netcat: 0 означает успех
  if [ $? -eq 0 ]; then
    echo "PostgreSQL is available after $i seconds."
    # Порт доступен, запускаем ваше Go-приложение
    exec ./algoplatform 
  fi
  sleep 1
done

echo "Error: Timeout exceeded ($TIMEOUT seconds). Could not connect to PostgreSQL."
exit 1
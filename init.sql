-- Создаём пользователя, если он ещё не существует
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'user') THEN
      CREATE ROLE "user" WITH LOGIN PASSWORD 'password';
   END IF;
END
$$;

-- Создаём базу данных, если она ещё не существует
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'algoplatform_db') THEN
      CREATE DATABASE algoplatform_db OWNER user;
   END IF;
END
$$;

-- Подключаемся к базе данных
\connect algoplatform_db

-- Назначаем владельца схемы public
ALTER SCHEMA public OWNER TO user;

-- Выдаём все права пользователю
GRANT ALL PRIVILEGES ON DATABASE algoplatform_db TO user;

-- DB and user are already created by official postgres image
-- via POSTGRES_USER / POSTGRES_DB env vars from docker-compose.
-- Keep only privilege alignment for the selected DB.

ALTER SCHEMA public OWNER TO "user";
GRANT ALL PRIVILEGES ON DATABASE algoplatform_db TO "user";

CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
  id            BIGSERIAL PRIMARY KEY,
  username      TEXT UNIQUE NOT NULL,
  email         CITEXT UNIQUE NOT NULL,
  role_user     TEXT NOT NULL CHECK (role_user IN ('admin','user')),
  password_hash TEXT NOT NULL,
  created_at    TIMESTAMP NOT NULL DEFAULT now()
);

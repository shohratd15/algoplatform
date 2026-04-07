CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
  id            BIGSERIAL PRIMARY KEY,
  username      TEXT UNIQUE NOT NULL,
  email         CITEXT UNIQUE NOT NULL,
  role_user     TEXT NOT NULL CHECK (role_user IN ('admin','user')),
  password_hash TEXT NOT NULL,
  created_at    TIMESTAMP NOT NULL DEFAULT now()
);

INSERT INTO users (username, email, role_user, password_hash) VALUES ('admin', 'admin@example.com', 'admin', '$2a$10$pDLdzBmVGebeubA772zpsuq5FOUUrUccrgE7W1wnkMVMwKoyVXTMe');
INSERT INTO users (username, email, role_user, password_hash) VALUES ('shohrat', 'shohratd15@gmail.com', 'admin', '$2a$10$HDYDeEYvZV3ucR7eYKc30Op65L47bJaSWE3lQEkexDGC82E3q.zD.');
INSERT INTO users (username, email, role_user, password_hash) VALUES ('user', 'user@example.com', 'user', '$2a$10$pDLdzBmVGebeubA772zpsuq5FOUUrUccrgE7W1wnkMVMwKoyVXTMe');
CREATE TABLE problems (
    id          BIGSERIAL PRIMARY KEY,
    slug        TEXT UNIQUE NOT NULL, -- например "two-sum"
    difficulty  TEXT NOT NULL CHECK (difficulty IN ('easy','medium','hard')),
    created_at  TIMESTAMP DEFAULT now()
);

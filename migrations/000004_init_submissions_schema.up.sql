CREATE TABLE submissions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT REFERENCES users(id) ON DELETE CASCADE,
    problem_id  BIGINT REFERENCES problems(id) ON DELETE CASCADE,
    language    VARCHAR(10) NOT NULL,
    source_code TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending','accepted','wrong_answer','runtime_error')),
    created_at  TIMESTAMP DEFAULT now()
);

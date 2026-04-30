CREATE TABLE submissions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT REFERENCES users(id) ON DELETE CASCADE,
    problem_id  BIGINT REFERENCES problems(id) ON DELETE CASCADE,
    language_id INTEGER NOT NULL,
    source_code TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('queued','running','accepted','wrong_answer','runtime_error','error','time_limit','memory_limit')),
    created_at  TIMESTAMP DEFAULT now(),
    updated_at  TIMESTAMP DEFAULT now()
);

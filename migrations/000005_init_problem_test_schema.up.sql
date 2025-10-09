CREATE TABLE problem_tests (
    id              BIGSERIAL PRIMARY KEY,
    problem_id      BIGINT REFERENCES problems(id) ON DELETE CASCADE,
    input_data      TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    is_sample       BOOLEAN NOT NULL DEFAULT FALSE
);
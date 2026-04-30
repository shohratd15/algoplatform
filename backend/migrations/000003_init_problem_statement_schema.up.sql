CREATE TABLE problem_statements (
    problem_id      BIGINT REFERENCES problems(id) ON DELETE CASCADE,
    language        VARCHAR(3) NOT NULL CHECK (language IN ('eng','rus','tkm')),
    title           TEXT NOT NULL,
    statement       TEXT NOT NULL,
    PRIMARY KEY (problem_id, language)
);

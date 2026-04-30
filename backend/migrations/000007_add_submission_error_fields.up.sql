ALTER TABLE submissions
    ADD COLUMN compile_output TEXT,
    ADD COLUMN stderr TEXT,
    ADD COLUMN message TEXT;

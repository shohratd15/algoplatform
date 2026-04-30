ALTER TABLE submissions
    DROP COLUMN IF EXISTS compile_output,
    DROP COLUMN IF EXISTS stderr,
    DROP COLUMN IF EXISTS message;

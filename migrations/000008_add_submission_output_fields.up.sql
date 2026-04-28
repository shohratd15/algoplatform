-- +migrate Up
ALTER TABLE submissions ADD COLUMN stdout TEXT;
ALTER TABLE submissions ADD COLUMN expected_output TEXT;

-- +migrate Down
ALTER TABLE submissions DROP COLUMN expected_output;
ALTER TABLE submissions DROP COLUMN stdout;
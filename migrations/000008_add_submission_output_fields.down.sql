-- +migrate Down
ALTER TABLE submissions DROP COLUMN expected_output;
ALTER TABLE submissions DROP COLUMN stdout;
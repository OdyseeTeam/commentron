-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE comment ADD COLUMN deleted_at TIMESTAMP NULL;
-- +migrate StatementEnd

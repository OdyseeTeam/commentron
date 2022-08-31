-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE comment ADD COLUMN is_protected BOOL NOT NULL DEFAULT false;
-- +migrate StatementEnd

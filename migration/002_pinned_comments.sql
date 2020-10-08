-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE comment ADD COLUMN is_pinned TINYINT(1) NOT NULL DEFAULT 0;
-- +migrate StatementEnd
-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE comment ADD COLUMN is_flagged TINYINT(1) NOT NULL DEFAULT 0;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD INDEX is_flagged_idx (is_flagged);
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE reaction ADD COLUMN is_flagged TINYINT(1) NOT NULL DEFAULT 0;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE reaction ADD INDEX is_flagged_idx (is_flagged);
-- +migrate StatementEnd
-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE blocked_list_invite CHANGE COLUMN accepted accepted TINYINT(1) NULL DEFAULT NULL;
-- +migrate StatementEnd
-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE blocked_list_appeal DROP FOREIGN KEY blocked_list_appeal_ibfk_1;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_list_appeal CHANGE COLUMN blocked_list_id blocked_list_id BIGINT(20) UNSIGNED NULL DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_list_appeal ADD CONSTRAINT blocked_list_appeal_ibfk_1 FOREIGN KEY (blocked_list_id) REFERENCES blocked_list (id);
-- +migrate StatementEnd


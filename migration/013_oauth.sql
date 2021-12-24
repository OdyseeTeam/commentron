-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE channel ADD COLUMN sub VARCHAR(50) DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE channel ADD INDEX sub_idx (sub, claim_id);
-- +migrate StatementEnd
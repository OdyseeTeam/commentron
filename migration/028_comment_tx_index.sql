-- +migrate Up
-- +migrate StatementBegin
ALTER TABLE comment add index idx_tx_id (tx_id);
-- +migrate StatementEnd

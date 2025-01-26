-- +migrate Up
-- +migrate StatementBegin
ALTER TABLE creator_setting
    ADD COLUMN min_usdc_tip_amount_comment BIGINT UNSIGNED DEFAULT NULL,
    ADD COLUMN min_usdc_tip_amount_super_chat BIGINT UNSIGNED DEFAULT NULL;
-- +migrate StatementEnd

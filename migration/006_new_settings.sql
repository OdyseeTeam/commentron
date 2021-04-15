-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE creator_setting CHANGE COLUMN min_tip_ammount_comment min_tip_amount_comment BIGINT(20) UNSIGNED NULL DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting CHANGE COLUMN min_tip_ammount_super_chat min_tip_amount_super_chat BIGINT(20) UNSIGNED NULL DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN slow_mode_min_gap BIGINT(20) UNSIGNED NULL DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD COLUMN amount BIGINT(20) UNSIGNED NULL DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD COLUMN tx_id VARCHAR(70) NULL DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD INDEX idx_amount (amount);
-- +migrate StatementEnd
-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE comment ADD COLUMN is_fiat BOOL NOT NULL DEFAULT false;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD COLUMN currency VARCHAR(25) DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD INDEX idx_comment_fiat_amount (is_fiat, amount, currency), ALGORITHM=INPLACE, LOCK=NONE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD INDEX idx_comment_is_fiat (is_fiat, timestamp), ALGORITHM=INPLACE, LOCK=NONE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD INDEX idx_comment_is_fiat_by_claim (lbry_claim_id, is_fiat, timestamp), ALGORITHM=INPLACE, LOCK=NONE;
-- +migrate StatementEnd
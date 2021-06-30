-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE comment ADD INDEX idx_comment_timestamp (lbry_claim_id,timestamp), ALGORITHM=INPLACE, LOCK=NONE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD INDEX idx_comment_hidden (is_hidden, comment_id, timestamp), ALGORITHM=INPLACE, LOCK=NONE;
-- +migrate StatementEnd
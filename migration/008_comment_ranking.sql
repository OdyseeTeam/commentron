-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE comment ADD COLUMN popularity_score INT DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD COLUMN controversy_score INT DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD INDEX idx_comment_popularity (lbry_claim_id, popularity_score, timestamp), ALGORITHM=INPLACE, LOCK=NONE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment ADD INDEX idx_comment_controversy (lbry_claim_id, controversy_score, timestamp), ALGORITHM=INPLACE, LOCK=NONE;
-- +migrate StatementEnd
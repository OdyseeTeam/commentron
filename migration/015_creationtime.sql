-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE channel ADD COLUMN created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE channel ADD COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE channel ADD INDEX idx_created_at (created_at);
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE channel ADD INDEX idx_updated_at (updated_at);
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN time_since_first_comment BIGINT DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN blocked_words_fuzziness_match BIGINT DEFAULT NULL;
-- +migrate StatementEnd
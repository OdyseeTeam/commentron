-- +migrate Up

-- +migrate StatementBegin
CREATE INDEX idx_comment_timestamp_alone on comment (timestamp);
-- +migrate StatementEnd

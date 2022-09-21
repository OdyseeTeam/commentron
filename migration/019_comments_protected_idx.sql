-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE comment ADD INDEX idx_is_protected (is_protected), ALGORITHM=INPLACE, LOCK=NONE;
-- +migrate StatementEnd

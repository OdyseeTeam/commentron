-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN livestream_chat_members_only BOOL NOT NULL DEFAULT FALSE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN comments_members_only BOOL NOT NULL DEFAULT FALSE;
-- +migrate StatementEnd

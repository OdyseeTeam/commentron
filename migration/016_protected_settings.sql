-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN public_show_protected BOOL NOT NULL DEFAULT FALSE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN private_show_protected BOOL NOT NULL DEFAULT FALSE;
-- +migrate StatementEnd

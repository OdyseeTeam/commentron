-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN chat_overlay BOOL NOT NULL DEFAULT true;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN chat_overlay_position VARCHAR(10) NOT NULL DEFAULT 'Left';
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN chat_remove_comment BIGINT NOT NULL DEFAULT 30;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN  sticker_overlay BOOL NOT NULL DEFAULT true;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN  sticker_overlay_keep BOOL NOT NULL DEFAULT false;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN sticker_overlay_remove BIGINT NOT NULL DEFAULT 10;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN viewercount_overlay BOOL NOT NULL DEFAULT true;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN viewercount_overlay_position VARCHAR(20) NOT NULL DEFAULT 'Top Left';
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN viewercount_chat_bot BOOL NOT NULL DEFAULT false;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN tipgoal_overlay BOOL NOT NULL DEFAULT true;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN tipgoal_amount BIGINT NOT NULL DEFAULT 1000;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN tipgoal_overlay_position VARCHAR(10) NOT NULL DEFAULT true;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN tipgoal_previous_donations BOOL NOT NULL DEFAULT true;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN tipgoal_currency VARCHAR(10) NOT NULL DEFAULT 'LBC';
-- +migrate StatementEnd
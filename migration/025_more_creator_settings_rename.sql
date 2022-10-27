-- +migrate Up
-- +migrate StatementBegin
ALTER TABLE creator_setting
    RENAME COLUMN featured_channels TO channel_sections;
-- +migrate StatementEnd

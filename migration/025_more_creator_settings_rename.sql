-- +migrate Up
-- +migrate StatementBegin
ALTER TABLE creator_setting
    DROP COLUMN featured_channels,
    ADD COLUMN channel_sections JSON DEFAULT NULL COMMENT 'array data for featured channels';
-- +migrate StatementEnd

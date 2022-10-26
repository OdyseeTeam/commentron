-- +migrate Up
-- +migrate StatementBegin
ALTER TABLE creator_setting
    ADD COLUMN featured_channels JSON DEFAULT NULL COMMENT 'array data for featured channels',
    ADD COLUMN homepage_settings JSON DEFAULT NULL COMMENT 'array data for homepage settings';
-- +migrate StatementEnd

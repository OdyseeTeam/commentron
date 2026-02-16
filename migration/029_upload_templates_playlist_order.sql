-- +migrate Up
-- +migrate StatementBegin
ALTER TABLE creator_setting
    ADD COLUMN upload_templates JSON DEFAULT NULL COMMENT 'json data for creator upload templates',
    ADD COLUMN playlist_order JSON DEFAULT NULL COMMENT 'json data for creator playlist order';
-- +migrate StatementEnd

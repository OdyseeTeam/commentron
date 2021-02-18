-- +migrate Up

-- +migrate StatementBegin
CREATE TABLE creator_setting (
    id                          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    creator_channel_id CHAR(40) NOT NULL,
    comments_enabled            BOOL DEFAULT FALSE,
    min_tip_ammount_comment     BIGINT UNSIGNED DEFAULT NULL,
    min_tip_ammount_super_chat  BIGINT UNSIGNED DEFAULT NULL,
    muted_words                 TEXT DEFAULT NULL,
    created_at                  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    INDEX idx_comments_enabled (comments_enabled),
    FOREIGN KEY fk_creator_channel (creator_channel_id) REFERENCES channel (claim_id)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- +migrate StatementEnd


-- +migrate StatementBegin
CREATE TABLE delegated_moderator (
     id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
     mod_channel_id   CHAR(40) NOT NULL,
     creator_channel_id CHAR(40) NOT NULL,
     permissons       BIGINT UNSIGNED NOT NULL DEFAULT 0,
     created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

     PRIMARY KEY (id),
     INDEX idx_permissions (permissons),
     FOREIGN KEY fk_mod_channel (mod_channel_id) REFERENCES channel (claim_id),
     FOREIGN KEY fk_creator_channel (creator_channel_id) REFERENCES channel (claim_id)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE channel ADD COLUMN is_spammer BOOLEAN DEFAULT 0;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE channel ADD INDEX idx_is_spammer (is_spammer);
-- +migrate StatementEnd
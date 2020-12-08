-- +migrate Up

-- +migrate StatementBegin
create table blocked_entry (
    id                      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    blocked_channel_id      CHAR(40) DEFAULT NULL,
    blocked_by_channel_id   CHAR(40) DEFAULT NULL,
    universally_blocked     BOOL DEFAULT FALSE,
    created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    INDEX idx_blockedall (blocked_channel_id,blocked_by_channel_id,universally_blocked),
    INDEX idx_blocked_by (blocked_by_channel_id),
    INDEX idx_universal (universally_blocked),
    FOREIGN KEY fk_blocked (blocked_channel_id) REFERENCES channel (claim_id),
    FOREIGN KEY fk_blocked_by (blocked_by_channel_id) REFERENCES channel (claim_id)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
create table moderator (
id                      BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
mod_channel_id          CHAR(40) DEFAULT NULL,
mod_level               BIGINT NOT NULL DEFAULT 1,
created_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at              DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

PRIMARY KEY (id),
FOREIGN KEY fk_blocked (mod_channel_id) REFERENCES channel (claim_id)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- +migrate StatementEnd
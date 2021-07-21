-- +migrate Up

-- +migrate StatementBegin
CREATE TABLE blocked_list (
 id                    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
 channel_id            CHAR(40) NOT NULL,
 name                  NVARCHAR(255) NOT NULL,
 category              NVARCHAR(255) NOT NULL,
 description           MEDIUMTEXT NOT NULL,
 member_invite_enabled BOOL DEFAULT FALSE,
 strike_one            BIGINT UNSIGNED DEFAULT NULL,
 strike_two            BIGINT UNSIGNED DEFAULT NULL,
 strike_three          BIGINT UNSIGNED DEFAULT NULL,
 invite_expiration     BIGINT UNSIGNED DEFAULT NULL,
 curse_jar_amount      BIGINT UNSIGNED DEFAULT NULL,
 created_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 updated_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

 PRIMARY KEY (id),
 INDEX idx_blocked_list_name (name),
 INDEX idx_category (category),
 INDEX idx_created_at (created_at),
 INDEX idx_updated_at (updated_at),
 FOREIGN KEY fk_channel (channel_id) REFERENCES channel (claim_id)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE blocked_list_invite (
id                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
blocked_list_id    BIGINT UNSIGNED NOT NULL,
inviter_channel_id CHAR(40) NOT NULL,
invited_channel_id CHAR(40) NOT NULL,
accepted           BOOL DEFAULT FALSE,
message            MEDIUMTEXT NOT NULL,
created_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at         DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

PRIMARY KEY (id),
UNIQUE idx_unique_invite (blocked_list_id,inviter_channel_id,invited_channel_id),
INDEX idx_created_at (created_at),
INDEX idx_updated_at (updated_at),
FOREIGN KEY fk_blocked_list (blocked_list_id) REFERENCES blocked_list (id),
FOREIGN KEY fk_channel (inviter_channel_id) REFERENCES channel (claim_id),
FOREIGN KEY fk_invited_channel (invited_channel_id) REFERENCES channel (claim_id)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE blocked_list_appeal (
 id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
 blocked_list_id  BIGINT UNSIGNED NOT NULL,
 blocked_entry_id BIGINT UNSIGNED NOT NULL,
 appeal           MEDIUMTEXT NOT NULL,
 response         MEDIUMTEXT NOT NULL,
 approved         BOOL DEFAULT NULL,
 escalated        BOOL DEFAULT FALSE,
 tx_id            VARCHAR(64) DEFAULT NULL,
 created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
 updated_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

 PRIMARY KEY (id),
 INDEX idx_appeal (blocked_list_id, approved),
 INDEX idx_escalated (blocked_list_id, escalated),
 INDEX idx_created_at (created_at),
 INDEX idx_updated_at (updated_at),
 FOREIGN KEY fk_blocked_list (blocked_list_id) REFERENCES blocked_list (id),
 FOREIGN KEY fk_blocked_entry (blocked_entry_id) REFERENCES blocked_entry (id)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE channel ADD COLUMN blocked_list_invite_id BIGINT UNSIGNED DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE channel ADD FOREIGN KEY fk_invite (blocked_list_invite_id) REFERENCES blocked_list (id) ON DELETE SET NULL ON UPDATE CASCADE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE channel ADD COLUMN blocked_list_id BIGINT UNSIGNED DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE channel ADD FOREIGN KEY fk_blocked_list (blocked_list_id) REFERENCES blocked_list (id) ON DELETE SET NULL ON UPDATE CASCADE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD COLUMN blocked_list_id BIGINT UNSIGNED DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD FOREIGN KEY fk_blocked_list (blocked_list_id) REFERENCES blocked_list (id) ON DELETE SET NULL ON UPDATE CASCADE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry CHANGE COLUMN blocked_by_channel_id creator_channel_id CHAR(40) DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD COLUMN delegated_moderator_channel_id CHAR(40) DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD FOREIGN KEY fk_mod (delegated_moderator_channel_id) REFERENCES channel (claim_id) ON DELETE RESTRICT ON UPDATE CASCADE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD COLUMN reason MEDIUMTEXT DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD COLUMN offending_comment_id CHAR(64) DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD FOREIGN KEY fk_offending (offending_comment_id) REFERENCES comment (comment_id) ON DELETE SET NULL ON UPDATE CASCADE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD INDEX idx_blockedlist (blocked_channel_id,creator_channel_id,blocked_list_id,universally_blocked);
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD COLUMN expiry DATETIME DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD INDEX idx_expiry (expiry);
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD COLUMN strikes INT DEFAULT 1;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE blocked_entry ADD INDEX idx_strikes (strikes);
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN curse_jar_amount BIGINT UNSIGNED DEFAULT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE creator_setting ADD COLUMN is_filters_enabled BOOLEAN DEFAULT NULL;
-- +migrate StatementEnd
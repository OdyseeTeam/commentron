-- +migrate Up

-- +migrate StatementBegin
ALTER DATABASE social
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE channel (
   claim_id VARCHAR(40)  NOT NULL,
   name  CHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
   CONSTRAINT channel_pk PRIMARY KEY (claim_id)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE comment (
    -- should be changed to CHAR(64)
    comment_id   CHAR(64) NOT NULL,
    -- should be changed to CHAR(40)
    lbry_claim_id CHAR(40) NOT NULL,
    -- can be null, so idk if this should be char(40)
    channel_id   CHAR(40) DEFAULT NULL,
    body        TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    parent_id    CHAR(64) DEFAULT NULL,
    signature   CHAR(128) DEFAULT NULL,
    -- 22 chars long is prolly enough
    signingts   VARCHAR(22) DEFAULT NULL,

    timestamp   INTEGER NOT NULL,
    -- there's no way that the timestamp will ever reach 22 characters
    is_hidden    BOOLEAN DEFAULT FALSE,
    CONSTRAINT COMMENT_PRIMARY_KEY PRIMARY KEY (comment_id)
    -- setting null implies comment is top level
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment
    ADD FOREIGN KEY (channel_id) REFERENCES channel (claim_id) ON DELETE CASCADE ON UPDATE CASCADE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE comment
    ADD FOREIGN KEY (parent_id) REFERENCES comment (comment_id) ON UPDATE CASCADE ON DELETE CASCADE;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE INDEX claim_comment_index ON comment (lbry_claim_id, comment_id);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE INDEX channel_comment_index ON comment (channel_id, comment_id);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE comment_opinion (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    comment_id char(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    channel_id char(40) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    signature char(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    signingts varchar(22) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    timestamp int NOT NULL,
    rating tinyint DEFAULT '1',

    PRIMARY KEY (id),
    FOREIGN KEY (comment_id) REFERENCES comment (comment_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (channel_id) REFERENCES channel (claim_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE content_opinion (
    id bigint unsigned NOT NULL AUTO_INCREMENT,
    claim_id char(40) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    channel_id char(40) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    signature char(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    signingts varchar(22) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    timestamp int NOT NULL,
    rating tinyint DEFAULT '1',

    PRIMARY KEY (id),
    FOREIGN KEY (channel_id) REFERENCES channel (claim_id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +migrate StatementEnd
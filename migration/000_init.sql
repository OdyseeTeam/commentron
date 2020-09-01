-- +migrate Up
-- +migrate StatementBegin
CREATE SCHEMA IF NOT EXISTS social;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS social.CHANNEL (
                           claimid varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL,
                           name varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
                           PRIMARY KEY (claimid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS social.COMMENT (
                           commentid char(64) COLLATE utf8mb4_unicode_ci NOT NULL,
                           lbryclaimid char(40) COLLATE utf8mb4_unicode_ci NOT NULL,
                           channelid char(40) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                           body text COLLATE utf8mb4_unicode_ci NOT NULL,
                           parentid char(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                           signature char(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                           signingts varchar(22) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
                           timestamp int(11) NOT NULL,
                           ishidden tinyint(1) DEFAULT '0',
                           PRIMARY KEY (commentid),
                           KEY comment_channel_fk (channelid),
                           KEY comment_parent_fk (parentid),
                           KEY lbryclaimid (lbryclaimid),
                           CONSTRAINT comment_channel_fk FOREIGN KEY (channelid) REFERENCES CHANNEL (claimid) ON DELETE CASCADE ON UPDATE CASCADE,
                           CONSTRAINT comment_parent_fk FOREIGN KEY (parentid) REFERENCES COMMENT (commentid) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +migrate StatementEnd

-- +migrate StatementBegin
USE commentron;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER DATABASE commentron
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

-- +migrate StatementBegin
INSERT INTO commentron.channel (channel.claim_id,channel.name)
SELECT c.claimid, c.name FROM social.CHANNEL c;
-- +migrate StatementEnd

-- +migrate StatementBegin
SET FOREIGN_KEY_CHECKS = 0;
-- +migrate StatementEnd

-- +migrate StatementBegin
INSERT INTO commentron.comment (
    comment.comment_id,
    comment.lbry_claim_id,
    comment.channel_id,
    comment.body,
    comment.parent_id,
    comment.signature,
    comment.signingts,
    comment.timestamp,
    comment.is_hidden)

SELECT c.commentid,
       c.lbryclaimid,
       c.channelid,
       c.body,
       c.parentid,
       c.signature,
       c.signingts,
       c.timestamp,
       c.ishidden
FROM social.COMMENT c
WHERE c.parentid IS NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
INSERT INTO commentron.comment (
    comment.comment_id,
    comment.lbry_claim_id,
    comment.channel_id,
    comment.body,
    comment.parent_id,
    comment.signature,
    comment.signingts,
    comment.timestamp,
    comment.is_hidden)

SELECT c.commentid,
       c.lbryclaimid,
       c.channelid,
       c.body,
       c.parentid,
       c.signature,
       c.signingts,
       c.timestamp,
       c.ishidden
FROM social.COMMENT c
WHERE c.parentid IS NOT NULL;
-- +migrate StatementEnd

-- +migrate StatementBegin
SET FOREIGN_KEY_CHECKS = 1;
-- +migrate StatementEnd
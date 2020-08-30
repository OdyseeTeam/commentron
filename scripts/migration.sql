
INSERT INTO social.channel (channel.claim_id,channel.name)
SELECT c.claimid, c.name FROM social_prod.CHANNEL c;

SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO social.comment (
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
FROM social_prod.COMMENT c
WHERE c.parentid IS NULL;

INSERT INTO social.comment (
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
FROM social_prod.COMMENT c
WHERE c.parentid IS NOT NULL;

SET FOREIGN_KEY_CHECKS = 1;
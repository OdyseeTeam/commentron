
INSERT INTO commentron.channel (channel.claim_id,channel.name)
SELECT c.claimid, c.name FROM social.CHANNEL c;

SET FOREIGN_KEY_CHECKS = 0;

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

SET FOREIGN_KEY_CHECKS = 1;
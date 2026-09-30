-- +migrate Up
-- +migrate StatementBegin
CREATE INDEX idx_protected_claim_parent_deleted ON comment(is_protected, lbry_claim_id, parent_id, deleted_at);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE INDEX idx_comment_list ON comment(lbry_claim_id, parent_id, is_flagged, deleted_at, is_pinned, popularity_score, timestamp);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE INDEX idx_channel_deleted ON comment(channel_id, deleted_at);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE INDEX idx_timestamp_desc ON comment_classification (timestamp DESC);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE INDEX idx_claim_amount_protected_deleted ON comment(lbry_claim_id, amount, is_protected, deleted_at);
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE INDEX idx_deleted_at ON comment(`deleted_at`);
-- +migrate StatementEnd

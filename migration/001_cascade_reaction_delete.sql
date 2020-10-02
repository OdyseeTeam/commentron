-- +migrate Up

-- +migrate StatementBegin
ALTER TABLE reaction DROP FOREIGN KEY reaction_ibfk_1;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE reaction DROP FOREIGN KEY reaction_ibfk_2;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE reaction DROP FOREIGN KEY reaction_ibfk_3;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE reaction
    ADD CONSTRAINT reaction_ibfk_1
        FOREIGN KEY (channel_id)
            REFERENCES channel (claim_id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE reaction
    ADD CONSTRAINT reaction_ibfk_2
        FOREIGN KEY (comment_id)
            REFERENCES comment (comment_id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
-- +migrate StatementEnd

-- +migrate StatementBegin
ALTER TABLE reaction
    ADD CONSTRAINT reaction_ibfk_3
        FOREIGN KEY (reaction_type_id)
            REFERENCES reaction_type (id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
-- +migrate StatementEnd
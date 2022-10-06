-- +migrate Up

-- +migrate StatementBegin
CREATE TABLE comment_classification (

    -- The `comment.comment_id`
    comment_id CHAR(64) NOT NULL PRIMARY KEY,

	-- This is the primary feature. We want to tune this.
	toxicity FLOAT NOT NULL,

	-- Specific forms of abuse let us adapt to other models in the
	-- future without having to rebuild this table too often.
	severe_toxicity FLOAT NOT NULL,
	obscene FLOAT NOT NULL,
	identity_attack FLOAT NOT NULL,
	insult FLOAT NOT NULL,
	threat FLOAT NOT NULL,
	sexual_explicit FLOAT NOT NULL,

	-- I don't have these yet but can add them in the future.
	nazi FLOAT NOT NULL,
	doxx FLOAT NOT NULL,

	-- Has a Odysee moderator reviewed this classification.
	is_reviewed BOOLEAN DEFAULT false,

	-- If they have reviewed it (not null) did they agree or disagree.
	reviewer_approved BOOLEAN,

    -- Note: The at column should replicate the comment's timestamp.
    -- It's worth denormalizing a bit to get an index on exactly what we want
    -- in the moderator interface.
    timestamp INT NOT NULL,

	-- Useful for auditing, worth the byte cost for simplicity
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	model_ident VARCHAR(20),

    CONSTRAINT comment_classification_comment_fk
       FOREIGN KEY (comment_id) REFERENCES comment (comment_id)
          ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;


-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE INDEX idx_comment_classification_for_mods on comment_classification (
    is_reviewed, timestamp, toxicity DESC
);
-- +migrate StatementEnd

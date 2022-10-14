-- +migrate Up

-- +migrate StatementBegin

CREATE TABLE claim_to_channel (
    claim_id   VARCHAR(40) PRIMARY KEY,
    channel_id VARCHAR(40) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;


-- +migrate StatementEnd

-- +migrate StatementEnd

CREATE TABLE channel_algo_callbacks (
    channel_id VARCHAR(40) NOT NULL,
    watcher_id INT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (channel_id, watcher_id)
) ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- +migrate StatementBegin

-- +migrate StatementBegin

CREATE INDEX idx_claim_to_channel_channel_id ON claim_to_channel (channel_id);

-- +migrate StatementEnd

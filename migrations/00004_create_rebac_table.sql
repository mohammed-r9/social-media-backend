-- +goose UP
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS follows (
    follower_id UUID NOT NULL REFERENCES users(id),
    following_id UUID NOT NULL REFERENCES users(id),
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id)
);

CREATE TABLE blocks (
    blocker_id UUID NOT NULL REFERENCES users(id),
    blocked_id UUID NOT NULL REFERENCES users(id),
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (blocker_id, blocked_id)
);

CREATE INDEX idx_follows_follower_id
ON follows(follower_id);

CREATE INDEX idx_follows_following_id
ON follows(following_id);

CREATE INDEX idx_follows_following_created_at
ON follows(following_id, created_at DESC);

CREATE INDEX idx_blocks_blocker_id
ON blocks(blocker_id);

CREATE INDEX idx_blocks_blocked_id
ON blocks(blocked_id);
-- +goose StatementEnd

-- +goose DOWN
-- +goose StatementBegin
DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS blocks;
-- +goose StatementEnd

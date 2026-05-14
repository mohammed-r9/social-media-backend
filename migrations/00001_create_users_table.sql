-- +goose UP
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL CHECK (length(trim(email)) > 0),
	phone TEXT UNIQUE DEFAULT NULL,
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    is_suspended BOOL NOT NULL DEFAULT false,
    password_hash BYTEA NOT NULL CHECK (octet_length(password_hash) > 0),
    verified_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users; 
-- +goose StatementEnd

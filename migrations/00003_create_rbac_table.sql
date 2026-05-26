-- +goose UP
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
-- +goose StatementEnd

-- +goose DOWN
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP INDEX IF EXISTS idx_user_roles_role_id;
DROP INDEX IF EXISTS idx_role_permissions_role_id;
DROP INDEX IF EXISTS idx_role_permissions_permission_id;
-- +goose StatementEnd

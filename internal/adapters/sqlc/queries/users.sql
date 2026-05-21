-- name: CreateUser :one
INSERT INTO users (
	id, email, name, password_hash, username 
) VALUES ( @id, @email, @name, @password_hash, @username )
RETURNING id, email, name, username, password_hash, verified_at, created_at, updated_at, is_suspended, phone;

-- name: GetUserByEmail :one
SELECT id, email, name, username, password_hash, verified_at, created_at, updated_at, is_suspended, phone 
FROM users
WHERE email = @email;

-- name: UpdateUserPassword :execrows
UPDATE users
	SET password_hash = @password_hash
	WHERE id = @id;

-- name: VerifyUserEmail :execrows
UPDATE users
	SET verified_at = CURRENT_TIMESTAMP
	WHERE id = @id;

-- name: GetUserByID :one
SELECT id, email, name, username, password_hash, verified_at, created_at, updated_at, is_suspended, phone 
FROM users
WHERE id = @id;

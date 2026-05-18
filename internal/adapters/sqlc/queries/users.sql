-- name: CreateUser :one
INSERT INTO users (
	id, email, name, password_hash 
) VALUES ( @id, @email, @name, @password_hash )
RETURNING id, email, name, password_hash, verified_at, created_at, updated_at, is_suspended, phone;

-- name: GetUserByEmail :one
SELECT id, email, name, password_hash, verified_at, created_at, updated_at, is_suspended, phone 
FROM users
WHERE email = @email;

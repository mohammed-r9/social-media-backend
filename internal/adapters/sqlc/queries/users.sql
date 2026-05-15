-- name: CreateUser :one
INSERT INTO users (
	id, email, name, password_hash 
) VALUES ( @id, @email, @name, @password_hash )
RETURNING id, email, name, created_at;

-- name: CreateSession :one
INSERT INTO sessions (
	id, user_id, refresh_token_hash, csrf_token_hash, device_name, expires_at
) VALUES ( @id, @user_id, @refresh_token_hash, @csrf_token_hash, @device_name, @expires_at )
RETURNING id, user_id, refresh_token_hash, csrf_token_hash, device_name, expires_at, revoked_at;

-- name: GetUserSession :one
SELECT id, user_id, refresh_token_hash, csrf_token_hash, device_name, expires_at, revoked_at 
FROM sessions
WHERE id = @id AND user_id = @user_id;

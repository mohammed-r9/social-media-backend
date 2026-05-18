-- name: CreateSession :one
INSERT INTO sessions (
	id, user_id, refresh_token_hash, csrf_token_hash, device_name, expires_at
) VALUES ( @id, @user_id, @refresh_token_hash, @csrf_token_hash, @device_name, @expires_at )
RETURNING id, user_id, refresh_token_hash, csrf_token_hash, device_name, expires_at, revoked_at;

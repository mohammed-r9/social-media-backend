-- name: CreateSession :one
INSERT INTO sessions (
	id, user_id, refresh_token_hash, csrf_token_hash, device_name, expires_at
) VALUES ( @id, @user_id, @refresh_token_hash, @csrf_token_hash, @device_name, @expires_at )
RETURNING id, user_id, refresh_token_hash, csrf_token_hash, device_name, expires_at, revoked_at;

-- name: GetUserSession :one
SELECT 
    s.id,
    s.user_id,
    s.refresh_token_hash,
    s.csrf_token_hash,
    s.device_name,
    s.expires_at,
    s.revoked_at,
	u.verified_at
FROM sessions s
JOIN users u ON u.id = s.user_id
WHERE s.id = @id;

-- name: RevokeSession :one
UPDATE sessions
	SET revoked_at = CURRENT_TIMESTAMP
	WHERE id = @id
RETURNING id;

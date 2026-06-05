-- name: CreateUser :one
INSERT INTO users (
	id, email, password_hash, username 
) VALUES ( @id, @email, @password_hash, @username )
RETURNING id, email, username, password_hash, verified_at, created_at, updated_at, is_suspended, phone;

-- name: GetUserByEmail :one
SELECT id, email, username, password_hash, verified_at, created_at, updated_at, is_suspended, phone 
FROM users
WHERE email = @email;

-- name: UpdateUserPassword :one
UPDATE users
	SET password_hash = @password_hash
	WHERE id = @id
RETURNING id;

-- name: VerifyUserEmail :one
UPDATE users
	SET verified_at = CURRENT_TIMESTAMP
	WHERE id = @id
RETURNING id;

-- name: GetUserByID :one
SELECT
    u.id,
    u.email,
    u.username,
    u.verified_at,
    u.created_at,
    u.updated_at,
	u.password_hash,
    u.is_suspended,
    u.phone,

    p.display_name,
    p.bio,
    p.avatar_key,
    p.website,
    p.followers_count,
    p.following_count,
    p.posts_count,
    p.created_at AS profile_created_at,
    p.updated_at AS profile_updated_at

FROM users u
LEFT JOIN profiles p ON p.user_id = u.id
WHERE u.id = @id;

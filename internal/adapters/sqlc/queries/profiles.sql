-- name: CreateProfile :one
INSERT INTO profiles (
	user_id, display_name	
) VALUES ( @user_id, @display_name )
	RETURNING user_id;

-- name: UpdateUserProfile :one
UPDATE profiles
SET bio = COALESCE(sqlc.narg('bio'), bio),
    display_name = COALESCE(sqlc.narg('display_name'), display_name),
    website = COALESCE(sqlc.narg('website'), website)
WHERE user_id = @user_id
RETURNING user_id, bio, display_name, avatar_key, website;

-- name: UpdateUserAvatar :one
UPDATE profiles
SET avatar_key = @avatar_key
WHERE user_id = @user_id
RETURNING user_id, avatar_key;

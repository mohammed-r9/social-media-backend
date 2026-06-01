-- name: CreateProfile :one
INSERT INTO profiles (
	user_id, display_name	
) VALUES ( @user_id, @display_name )
	RETURNING user_id;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY display_name;

-- name: FindUserByDisplayName :one
SELECT user_id, password_hash, first_name, last_name, pfp_url FROM users
WHERE display_name = @display_name;

-- name: FindUserByID :one
SELECT display_name, first_name, last_name, pfp_url FROM users
WHERE user_id = @user_id;

-- name: CreateUser :one
INSERT INTO users (
    display_name,
    first_name,
    last_name,
    password_hash,
    user_key,
    pfp_url
) VALUES (
    @display_name, @first_name, @last_name, @password_hash, @user_key, @pfp_url
)
RETURNING user_id;

-- name: UpdateUserPfp :exec
UPDATE users
SET pfp_url = $1
WHERE user_id = $2;

-- name: FindUserKey :one
SELECT user_key FROM users
WHERE user_id = @user_id;
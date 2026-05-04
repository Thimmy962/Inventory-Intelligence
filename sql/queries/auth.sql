-- name: GetLevel :one
SELECT is_manager FROM staffs WHERE id = $1;


-- name: LoginUser :one
SELECT id, username, is_manager
FROM staffs WHERE username = $1 AND pword = $2;



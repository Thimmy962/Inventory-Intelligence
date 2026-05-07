-- name: GetLevel :one
SELECT is_manager FROM staffs WHERE id = $1;


-- name: LoginUser :one
SELECT id, username, is_manager, pword
FROM staffs WHERE username = $1;


-- name: CreateStaff :exec
INSERT INTO staffs (id, first_name, last_name, pword, is_manager, username)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5);



-- name: TryLogin :one
SELECT * FROM staffs
WHERE username = $1 and pword = $2;

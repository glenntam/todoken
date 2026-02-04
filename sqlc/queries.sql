-- name: CreateUser :exec
INSERT INTO users (email, password, created)
VALUES ($1, $2, $3);

-- name: GetUserByEmail :one
SELECT * FROM users WHERE LOWER(email) = LOWER($1);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users SET password = $2 where id = $1;

-- name: InsertPasswordReset :exec
INSERT INTO password_resets (token, user_id, expires)
VALUES ($1, $2, $3);

-- name: GetPasswordResetByToken :one
SELECT * FROM password_resets
WHERE token = $1 AND expires > $2;

-- name: DeletePasswordResetByID :exec
DELETE FROM password_resets WHERE user_id = $1;


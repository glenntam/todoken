-- name: CreateUser :exec
INSERT INTO users (email, password_hash, created_at)
VALUES (?, ?, ?);

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: CreateSession :exec
INSERT INTO sessions (token, user_id, expires_at)
VALUES (?, ?, ?);

-- name: GetSession :one
SELECT * FROM sessions WHERE token = ?;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE token = ?;

-- name: CreateTask :exec
INSERT INTO tasks (user_id, created_at, category, text)
VALUES (?, ?, ?, ?);

-- name: GetTasksByUser :many
SELECT * FROM tasks WHERE user_id = ? ORDER BY created_at;

-- name: UpdateTask :exec
UPDATE tasks SET category = ?, text = ?, completed_at = ? WHERE id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = ?;

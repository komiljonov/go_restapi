-- name: AllUsers :many
SELECT *
FROM users;


-- name: CreateUser :one
INSERT INTO users (name,
                   phone_number,
                   password,
                   birthdate)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: GetByPhoneNumber :one
SELECT *
FROM users
WHERE phone_number = $1;


-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;


-- name: UpdateUser :one
UPDATE users
SET name         = $2,
    birthdate    = $3
where id = $1
RETURNING *;



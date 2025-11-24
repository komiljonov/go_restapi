-- name: createUser :one
INSERT INTO users (
  name,phone_number,password,birthday
) VALUES (
  $1,$2,$3,$4
) RETURNING id,name,phone_number,password,birthday;
-- name: CreatePerson :one
INSERT INTO persons (external_id, email, password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPerson :one
SELECT *
FROM persons
WHERE external_id = $1;

-- name: GetPersonByEmail :one
SELECT *
FROM persons
WHERE email = $1;

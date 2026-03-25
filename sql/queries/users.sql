-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUserEmail :one
UPDATE users
SET 
    email = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET 
    hashed_password = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserChirpyRed :one
UPDATE users
SET 
    is_chirpy_red = true,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
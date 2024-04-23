-- name: FindUserByID :one
SELECT *
FROM users
WHERE id = @id;

-- name: MarkUserAsVerified :one
UPDATE users
SET verified           = true,
    verification_token = null
WHERE id = @id
  AND verified = false
  AND verification_token IS NOT NULL
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (email, password, verified, verification_token)
VALUES (@email, @password, @verified, @verification_token)
RETURNING *;

-- name: FindUserByVerificationToken :one
SELECT *
FROM users
WHERE verification_token = @verification_token;


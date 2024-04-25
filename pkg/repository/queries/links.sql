-- name: GetLinkByID :one
SELECT *
FROM links
WHERE id = $1;

-- name: GetUserLinks :many
SELECT *
FROM links
WHERE user_id = @UserId
ORDER BY created_at DESC
LIMIT @maxResults OFFSET @queryOffset;

-- name: GetLinkBySlug :one
SELECT *
FROM links
WHERE slug = @slug;

-- name: CreateLink :one
INSERT INTO links (user_id, slug, description, expire_at, target_url)
VALUES (@user_id, @slug, @description, @expire_at, @target_url)
RETURNING *;

-- name: UpdateLink :one
UPDATE links
SET description = @description,
    target_url  = @target_url,
    updated_at  = NOW()
WHERE id = @id
  AND user_id = @user_id
RETURNING *;

-- name: DeleteLinkByID :exec
DELETE
FROM links
WHERE id = $1;

-- name: GetExpiredLinks :many
SELECT *
FROM links
WHERE expire_at < NOW()
LIMIT @maxResults FOR UPDATE SKIP LOCKED;

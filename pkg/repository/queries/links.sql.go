// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: links.sql

package queries

import (
	"context"
	"time"
)

const countUserLinks = `-- name: CountUserLinks :one
SELECT COUNT(*)
FROM links
WHERE user_id = $1
`

func (q *Queries) CountUserLinks(ctx context.Context, userid int64) (int64, error) {
	row := q.db.QueryRow(ctx, countUserLinks, userid)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createLink = `-- name: CreateLink :one
INSERT INTO links (user_id, slug, description, expire_at, target_url)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, slug, description, banned, expire_at, target_url, created_at, updated_at
`

type CreateLinkParams struct {
	UserID      int64
	Slug        string
	Description *string
	ExpireAt    *time.Time
	TargetUrl   string
}

func (q *Queries) CreateLink(ctx context.Context, arg CreateLinkParams) (Link, error) {
	row := q.db.QueryRow(ctx, createLink,
		arg.UserID,
		arg.Slug,
		arg.Description,
		arg.ExpireAt,
		arg.TargetUrl,
	)
	var i Link
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Slug,
		&i.Description,
		&i.Banned,
		&i.ExpireAt,
		&i.TargetUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteLinkByID = `-- name: DeleteLinkByID :exec
DELETE
FROM links
WHERE id = $1
`

func (q *Queries) DeleteLinkByID(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteLinkByID, id)
	return err
}

const getExpiredLinks = `-- name: GetExpiredLinks :many
SELECT id, user_id, slug, description, banned, expire_at, target_url, created_at, updated_at
FROM links
WHERE expire_at < NOW()
LIMIT $1 FOR UPDATE SKIP LOCKED
`

func (q *Queries) GetExpiredLinks(ctx context.Context, maxresults int32) ([]Link, error) {
	rows, err := q.db.Query(ctx, getExpiredLinks, maxresults)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Link
	for rows.Next() {
		var i Link
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Slug,
			&i.Description,
			&i.Banned,
			&i.ExpireAt,
			&i.TargetUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLinkByID = `-- name: GetLinkByID :one
SELECT id, user_id, slug, description, banned, expire_at, target_url, created_at, updated_at
FROM links
WHERE id = $1
`

func (q *Queries) GetLinkByID(ctx context.Context, id int64) (Link, error) {
	row := q.db.QueryRow(ctx, getLinkByID, id)
	var i Link
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Slug,
		&i.Description,
		&i.Banned,
		&i.ExpireAt,
		&i.TargetUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLinkBySlug = `-- name: GetLinkBySlug :one
SELECT id, user_id, slug, description, banned, expire_at, target_url, created_at, updated_at
FROM links
WHERE slug = $1
`

func (q *Queries) GetLinkBySlug(ctx context.Context, slug string) (Link, error) {
	row := q.db.QueryRow(ctx, getLinkBySlug, slug)
	var i Link
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Slug,
		&i.Description,
		&i.Banned,
		&i.ExpireAt,
		&i.TargetUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserLinks = `-- name: GetUserLinks :many
SELECT id, user_id, slug, description, banned, expire_at, target_url, created_at, updated_at
FROM links
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $3 OFFSET $2
`

type GetUserLinksParams struct {
	Userid      int64
	Queryoffset int32
	Maxresults  int32
}

func (q *Queries) GetUserLinks(ctx context.Context, arg GetUserLinksParams) ([]Link, error) {
	rows, err := q.db.Query(ctx, getUserLinks, arg.Userid, arg.Queryoffset, arg.Maxresults)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Link
	for rows.Next() {
		var i Link
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Slug,
			&i.Description,
			&i.Banned,
			&i.ExpireAt,
			&i.TargetUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateLink = `-- name: UpdateLink :one
UPDATE links
SET description = $1,
    target_url  = $2,
    expire_at   = $3,
    updated_at  = NOW()
WHERE id = $4
  AND user_id = $5
RETURNING id, user_id, slug, description, banned, expire_at, target_url, created_at, updated_at
`

type UpdateLinkParams struct {
	Description *string
	TargetUrl   string
	ExpireAt    *time.Time
	ID          int64
	UserID      int64
}

func (q *Queries) UpdateLink(ctx context.Context, arg UpdateLinkParams) (Link, error) {
	row := q.db.QueryRow(ctx, updateLink,
		arg.Description,
		arg.TargetUrl,
		arg.ExpireAt,
		arg.ID,
		arg.UserID,
	)
	var i Link
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Slug,
		&i.Description,
		&i.Banned,
		&i.ExpireAt,
		&i.TargetUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

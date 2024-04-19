package repository

import (
	"context"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/db"
	"tikkin/pkg/model"
)

type LinksRepository struct {
	db *db.DB
}

func NewLinksRepository(db *db.DB) LinksRepository {
	return LinksRepository{db: db}
}

func (l *LinksRepository) GetLink(id int) (*model.Link, error) {
	link := model.Link{}
	err := l.db.Pool.QueryRow(context.Background(),
		"SELECT id, user_id, slug, description, banned, expire_at, target_url, created_at, updated_at FROM links WHERE id = $1", id).
		Scan(&link.ID, &link.UserId, &link.Slug, &link.Description, &link.Banned,
			&link.ExpireAt, &link.TargetUrl, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		log.Err(err).Msg("Failed to find link")
		return nil, err
	}

	return &link, nil
}

func (l *LinksRepository) GetUserLinks(id int64, page int) ([]model.Link, error) {
	rows, err := l.db.Pool.Query(context.Background(),
		"SELECT id, slug, description, banned, expire_at, target_url, created_at, updated_at FROM links WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		id, 20, page*20)
	if err != nil {
		log.Err(err).Msg("Failed to get links")
		return nil, err
	}

	links := make([]model.Link, 0)
	for rows.Next() {
		link := model.Link{}
		err = rows.Scan(&link.ID, &link.Slug, &link.Description, &link.Banned,
			&link.ExpireAt, &link.TargetUrl, &link.CreatedAt, &link.UpdatedAt)
		if err != nil {
			log.Err(err).Msg("Failed to scan link")
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func (l *LinksRepository) GetLinkBySlug(slug string) (*model.Link, error) {
	link := model.Link{}
	err := l.db.Pool.QueryRow(context.Background(),
		"SELECT id, user_id, slug, description, banned, expire_at, target_url, created_at, updated_at FROM links WHERE slug = $1", slug).
		Scan(&link.ID, &link.UserId, &link.Slug, &link.Description, &link.Banned,
			&link.ExpireAt, &link.TargetUrl, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		log.Err(err).Msg("Failed to find link")
		return nil, err
	}

	return &link, nil
}

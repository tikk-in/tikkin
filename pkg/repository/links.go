package repository

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/config"
	"tikkin/pkg/db"
	"tikkin/pkg/model"
	"tikkin/pkg/utils"
)

type LinksRepository struct {
	db     *db.DB
	config *config.Config
}

func NewLinksRepository(db *db.DB, config *config.Config) LinksRepository {
	return LinksRepository{db: db, config: config}
}

func (l *LinksRepository) GetLink(id int64) (*model.Link, error) {
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

func (l *LinksRepository) CreateLink(link model.Link) (*model.Link, error) {
	// Create a new link
	log.Info().
		Str("slug", link.Slug).
		Str("description", link.Description).
		Str("target_url", link.TargetUrl).
		Msg("Creating link")

	if link.Slug == "" {
		link.Slug = utils.GenerateSlug(l.config.Links.Length)
	}

	if utils.Contains(utils.BlockedSlugs, link.Slug) {
		log.Warn().Str("slug", link.Slug).Msg("Blocked slug")
		return nil, errors.New("slug_denied")
	}

	linkId := int64(0)
	err := l.db.Pool.QueryRow(context.Background(),
		"INSERT INTO links (user_id, slug, description, expire_at, target_url) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		link.UserId, link.Slug, link.Description, nil, link.TargetUrl).Scan(&linkId)

	if err != nil {
		log.Err(err).Msg("Failed to create link")
		return nil, err
	}

	return l.GetLink(linkId)
}

func (l *LinksRepository) DeleteLink(id int64) error {
	_, err := l.db.Pool.Exec(context.Background(), "DELETE FROM links WHERE id = $1", id)
	if err != nil {
		log.Err(err).Msg("Failed to delete link")
		return err
	}

	return nil
}

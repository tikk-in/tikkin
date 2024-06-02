package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/model"
	"tikkin/pkg/repository/queries"
	"tikkin/pkg/utils"
)

func (r *Repository) GetLinkByID(id int64) (*model.Link, error) {

	linkEntity, err := r.db.Queries(context.Background()).GetLinkByID(context.Background(), id)
	if err != nil {
		log.Err(err).Msg("Failed to find link")
		return nil, err
	}

	return linkEntity.ToModel(), nil
}

func (r *Repository) GetUserLinks(userId int64, page, pageSize int32) ([]model.Link, error) {

	params := queries.GetUserLinksParams{
		Userid:      userId,
		Maxresults:  pageSize,
		Queryoffset: page * pageSize,
	}

	results, err := r.db.Queries(context.Background()).GetUserLinks(context.Background(), params)
	if err != nil {
		log.Err(err).Msg("Failed to find user links")
		return nil, err
	}

	var links []model.Link
	for _, result := range results {
		links = append(links, *result.ToModel())
	}

	return links, nil
}

func (r *Repository) GetLinkBySlug(slug string) (*model.Link, error) {
	link, err := r.db.Queries(context.Background()).GetLinkBySlug(context.Background(), slug)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, nil
		}
		log.Err(err).Msg("Failed to find link")
		return nil, err
	}
	return link.ToModel(), nil
}

func (r *Repository) CreateLink(link model.Link) (*model.Link, error) {
	// Create a new link
	var expireAt = ""
	if link.ExpireAt != nil {
		expireAt = link.ExpireAt.String()
	}

	log.Info().
		Str("slug", link.Slug).
		Str("description", link.Description).
		Str("target_url", link.TargetUrl).
		Str("expire_at", expireAt).
		Msg("Creating link")

	if link.Slug == "" {
		link.Slug = utils.GenerateSlug(r.config.Links.Length)
	}

	if utils.Contains(utils.BlockedSlugs, link.Slug) {
		log.Warn().Str("slug", link.Slug).Msg("Blocked slug")
		return nil, errors.New("slug_denied")
	}

	res, err := r.db.Queries(context.Background()).CreateLink(context.Background(),
		queries.CreateLinkParams{
			UserID:      link.UserId,
			Slug:        link.Slug,
			Description: &link.Description,
			ExpireAt:    link.ExpireAt,
			TargetUrl:   link.TargetUrl,
		})

	if err != nil {
		log.Err(err).Msg("Failed to create link")
		return nil, err
	}

	return res.ToModel(), nil
}

func (r *Repository) DeleteLink(ctx context.Context, id int64) error {
	log.Info().Int64("id", id).Msg("Deleting link")
	err := r.db.Queries(ctx).DeleteLinkByID(context.Background(), id)
	if err != nil {
		log.Err(err).Msg("Failed to delete link")
		return err
	}

	log.Info().Int64("id", id).Msg("Link deleted")
	return nil
}

func (r *Repository) UpdateLink(id int64, link model.Link) (*model.Link, error) {

	updatedLink, err := r.db.Queries(context.Background()).UpdateLink(context.Background(), queries.UpdateLinkParams{
		ID:          id,
		UserID:      link.UserId,
		Description: &link.Description,
		TargetUrl:   link.TargetUrl,
	})

	if err != nil {
		log.Err(err).Msg("Failed to update link")
		return nil, err
	}

	return updatedLink.ToModel(), nil
}

func (r *Repository) GetExpiredLinks(ctx context.Context) ([]model.Link, error) {
	result, err := r.db.Queries(ctx).GetExpiredLinks(ctx, 10)
	if err != nil {
		log.Err(err).Msg("Failed to get expired links")
		return nil, err
	}

	var links []model.Link
	for _, link := range result {
		links = append(links, *link.ToModel())
	}
	return links, nil
}

func (r *Repository) CountUserLinks(userId int64) (int64, error) {
	count, err := r.db.Queries(context.Background()).CountUserLinks(context.Background(), userId)
	if err != nil {
		log.Err(err).Msg("Failed to count user links")
		return 0, err
	}
	return count, nil
}

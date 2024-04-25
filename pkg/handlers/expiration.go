package handlers

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/config"
	"tikkin/pkg/db"
	"tikkin/pkg/repository"
	"time"
)

type ExpirationHandler struct {
	db             *db.DB
	config         *config.Config
	linkRepository *repository.LinksRepository
}

func NewExpirationHandler(db *db.DB, config *config.Config, linksRepository repository.LinksRepository) ExpirationHandler {
	return ExpirationHandler{
		db:             db,
		config:         config,
		linkRepository: &linksRepository,
	}
}

var NoLinksToExpireErr = errors.New("no_links_to_expire")

func (e *ExpirationHandler) deleteLinksBatch() error {
	ctx := context.Background()

	return e.db.WithTx(ctx, func(ctx context.Context) error {
		links, err := e.linkRepository.GetExpiredLinks(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get expired links")
			return err
		}
		if len(links) == 0 {
			return NoLinksToExpireErr
		}
		for _, link := range links {
			err = e.linkRepository.DeleteLink(ctx, link.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to delete expired link. Continuing...")
			}
		}
		return nil
	})
}

func (e *ExpirationHandler) ExpirationLoop() {

	log.Info().Msg("Starting link expiration loop...")

	for {
		err := e.deleteLinksBatch()
		if err != nil {
			if !errors.Is(err, NoLinksToExpireErr) {
				log.Error().Err(err).Msg("Failed to delete links batch")
			}
			time.Sleep(1 * time.Minute)
		}
	}

}

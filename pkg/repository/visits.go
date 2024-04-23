package repository

import (
	"context"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/db"
	"tikkin/pkg/model"
	"tikkin/pkg/repository/queries"
)

type VisitsRepository struct {
	db *db.DB
}

func NewVisitsRepository(db *db.DB) VisitsRepository {
	return VisitsRepository{db: db}
}

func (l *VisitsRepository) InsertVisit(visit model.Visits) {

	_, err := l.db.Queries.InsertVisit(context.Background(), queries.InsertVisitParams{
		ID:          visit.ID,
		LinkID:      visit.LinkID,
		UserAgent:   visit.UserAgent,
		Referrer:    visit.Referrer,
		CountryCode: visit.CountryCode,
	})

	if err != nil {
		log.Err(err).Msg("Failed to insert visit")
		return
	}
}

func (l *VisitsRepository) CountVisits(link model.Link) int64 {
	count, err := l.db.Queries.CountVisitsByLinkID(context.Background(), link.ID)
	if err != nil {
		log.Err(err).Msg("Failed to count visits")
		return 0
	}
	return count
}

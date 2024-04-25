package repository

import (
	"context"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/model"
	"tikkin/pkg/repository/queries"
)

func (r *Repository) InsertVisit(visit model.Visits) {

	_, err := r.db.Queries(context.Background()).InsertVisit(context.Background(), queries.InsertVisitParams{
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

func (r *Repository) CountVisits(link model.Link) int64 {
	count, err := r.db.Queries(context.Background()).CountVisitsByLinkID(context.Background(), link.ID)
	if err != nil {
		log.Err(err).Msg("Failed to count visits")
		return 0
	}
	return count
}

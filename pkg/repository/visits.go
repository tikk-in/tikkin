package repository

import (
	"context"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/db"
	"tikkin/pkg/model"
)

type VisitsRepository struct {
	db *db.DB
}

func NewVisitsRepository(db *db.DB) VisitsRepository {
	return VisitsRepository{db: db}
}

func (l *VisitsRepository) InsertVisit(visit model.Visits) {
	_, err := l.db.Pool.Exec(context.Background(),
		"INSERT INTO visits (id, link_id, user_agent, referrer, country_code) VALUES ($1, $2, $3, $4, $5)", visit.ID, visit.LinkID, visit.UserAgent, visit.Referrer, visit.CountryCode)
	if err != nil {
		log.Err(err).Msg("Failed to increment visits")
		return
	}
}

func (l *VisitsRepository) CountVisits(link model.Link) int64 {
	var count int64
	err := l.db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM visits WHERE link_id = $1", link.ID).Scan(&count)
	if err != nil {
		log.Err(err).Msg("Failed to count visits")
		return 0
	}
	return count
}

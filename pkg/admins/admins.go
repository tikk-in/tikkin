package admins

import (
	"context"
	"github.com/rs/zerolog/log"
	"tikkin/pkg"
	"tikkin/pkg/config"
	"tikkin/pkg/model"
	"tikkin/pkg/utils"
)

func EnsureAdmin(cfg *config.Config, db *pkg.DB, adminPassword string) {

	log.Info().Msg("Ensuring admin...")

	connection, err := db.Pool.Acquire(context.Background())
	if err != nil {
		log.Panic().Err(err).Msg("Failed to acquire connection")
	}
	defer connection.Release()

	admin := model.User{}

	row := db.Pool.QueryRow(context.Background(), "SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1", cfg.Admin.Email)
	err = row.Scan(&admin.ID, &admin.Email, &admin.Password, &admin.CreatedAt, &admin.UpdatedAt)
	if err != nil {
		if err.Error() != "no rows in result set" {
			log.Panic().Err(err).Msg("Failed to query for admin")
		}
	}

	if admin.ID != 0 {
		log.Info().Str("email", cfg.Admin.Email).Msg("Admin found. Skipping creation...")
		return
	}

	if adminPassword == "" {
		log.Info().Msg("Admin not found and no password provided, skipping admin creation...")
		return
	}

	log.Info().Str("email", cfg.Admin.Email).Msg("Admin not found, creating...")

	hash, err := utils.HashPassword(adminPassword)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to hash password")
	}

	_, err = db.Pool.Exec(context.Background(), "INSERT INTO users (email, password) VALUES ($1, $2)",
		cfg.Admin.Email, hash)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create admin")
	}
	log.Info().Str("email", cfg.Admin.Email).Msg("Admin created with provided password")
}

package admins

import (
	"context"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/config"
	"tikkin/pkg/db"
	"tikkin/pkg/model"
	"tikkin/pkg/utils"
)

func EnsureAdmin(cfg *config.Config, db *db.DB) {

	if cfg.Admin.Email == "" || cfg.Admin.Password == "" {
		log.Info().Msg("No admin email provided, skipping admin creation...")
		return
	}

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
		log.Info().Str("email", cfg.Admin.Email).Msg("Admin found. Checking credentials...")
		if !utils.DoesPasswordHashMatch(cfg.Admin.Password, admin.Password) {
			log.Warn().Msg("Admin password does not match provided password. Updating password...")
			hash, err := utils.HashPassword(cfg.Admin.Password)
			if err != nil {
				log.Panic().Err(err).Msg("Failed to hash password")
			}
			_, err = db.Pool.Exec(context.Background(), "UPDATE users SET password = $1 WHERE id = $2", hash, admin.ID)
			if err != nil {
				log.Panic().Err(err).Msg("Failed to update admin password")
			}
			log.Info().Str("email", cfg.Admin.Email).Msg("Admin password updated")
		}
		return
	}

	if cfg.Admin.Password == "" {
		log.Info().Msg("Admin not found and no password provided, skipping admin creation...")
		return
	}

	log.Info().Str("email", cfg.Admin.Email).Msg("Admin not found, creating...")

	hash, err := utils.HashPassword(cfg.Admin.Password)
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

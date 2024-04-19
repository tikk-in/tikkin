package main

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"strconv"
	"tikkin/pkg/admins"
	"tikkin/pkg/config"
	"tikkin/pkg/db"
	"tikkin/pkg/email"
	"tikkin/pkg/handlers"
)

func main() {

	log.Info().Msg("Starting Tikkin")

	// load config
	cfgFlags, err := config.ParseFlags()
	if err != nil {
		log.Fatal().Err(err)
	}
	cfg, err := config.LoadConfig(cfgFlags.ConfigPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	cfg.Email.SMTP.Password = cfgFlags.SMTPPassword

	db := db.NewDB(*cfg)

	admins.EnsureAdmin(cfg, db, cfgFlags.AdminPassword)

	emailHandler := email.NewEmailHandler(cfg)
	userHandler := handlers.NewUserHandler(db, cfg, &emailHandler)
	linkHandler := handlers.NewLinkHandler(db, cfg)
	redirectHandler := handlers.NewRedirectHandler(linkHandler)

	app := fiber.New()

	loginHandler := handlers.NewLoginHandler(cfg, db, userHandler)

	// Unauthenticated routes
	app.Post("/api/v1/auth/login", loginHandler.HandleLogin)
	app.Post("/api/v1/auth/signup", loginHandler.HandleRegister)
	app.Get("/api/v1/auth/verify/:token", userHandler.HandleVerification)
	app.Get("/:slug", redirectHandler.HandleRedirect)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.Server.Jwt.Secret)},
	}))

	// Authenticated routes
	app.Post("/api/v1/links", linkHandler.HandleCreateLink)
	app.Put("/api/v1/links/:id", linkHandler.HandleUpdateLink)
	app.Delete("/api/v1/links/:id", linkHandler.HandleDeleteLink)
	app.Get("/api/v1/links", linkHandler.HandleGetLinks)

	app.Listen(":" + strconv.Itoa(cfg.Server.Port))
}

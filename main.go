package main

import (
	"flag"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"strconv"
	"tikkin/pkg"
	"tikkin/pkg/admins"
	"tikkin/pkg/config"
)

func main() {

	log.Info().Msg("Starting Tikkin")

	adminPassword := ""
	flag.StringVar(&adminPassword, "admin-password", "", "Admin password")
	flag.Parse()

	// load config
	cfgPath, err := config.ParseFlags()
	if err != nil {
		log.Fatal().Err(err)
	}
	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatal().Err(err)
	}

	db := pkg.NewDB(*cfg)

	admins.EnsureAdmin(cfg, db, adminPassword)

	userHandler := pkg.NewUserHandler(db)
	linkHandler := pkg.NewLinkHandler(db, cfg)
	redirectHandler := pkg.NewRedirectHandler(linkHandler)

	app := fiber.New()

	loginHandler := pkg.NewLoginHandler(cfg, db, userHandler)

	// Unauthenticated routes
	app.Post("/api/v1/auth/login", loginHandler.HandleLogin)
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

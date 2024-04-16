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
	"tikkin/pkg/email"
)

func main() {

	log.Info().Msg("Starting Tikkin")

	adminPassword := ""
	smtpPassword := ""
	flag.StringVar(&adminPassword, "admin-password", "", "Admin password")
	flag.StringVar(&smtpPassword, "smtp-password", "", "SMTP password")
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

	cfg.Email.SMTP.Password = smtpPassword

	db := pkg.NewDB(*cfg)

	admins.EnsureAdmin(cfg, db, adminPassword)

	emailHandler := email.NewEmailHandler(cfg)
	userHandler := pkg.NewUserHandler(db, cfg, &emailHandler)
	linkHandler := pkg.NewLinkHandler(db, cfg)
	redirectHandler := pkg.NewRedirectHandler(linkHandler)

	app := fiber.New()

	loginHandler := pkg.NewLoginHandler(cfg, db, userHandler)

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

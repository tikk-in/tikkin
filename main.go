package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"
	"strconv"
	"tikkin/pkg/admins"
	"tikkin/pkg/config"
	"tikkin/pkg/db"
	"tikkin/pkg/email"
	"tikkin/pkg/handlers"
	"tikkin/pkg/repository"

	_ "tikkin/docs"
)

func enableDocs(app *fiber.App, cfg *config.Config) {

	siteUrl := cfg.Site.URL
	if cfg.Docs.OverrideSiteUrl != "" {
		siteUrl = cfg.Docs.OverrideSiteUrl
	}

	app.Get("/api-docs/*", swagger.New(swagger.Config{ // custom
		URL:          siteUrl + "/doc.json",
		DeepLinking:  false,
		DocExpansion: "list",
		CustomStyle:  `.swagger-ui .topbar { display: none }`,
	}))

	app.Get("/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("docs/swagger.json")
	})
}

// @title Tikkin API
// @version 1.0
// @description This is the Tikkin API documentation.
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	log.Info().Msg("Starting Tikkin...")

	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	db := db.NewDB(*cfg)
	admins.EnsureAdmin(cfg, db)

	repo := repository.NewRepository(db, cfg)
	emailHandler := email.NewEmailHandler(cfg)
	userHandler := handlers.NewUserHandler(db, cfg, &emailHandler, repo)
	linkHandler := handlers.NewLinkHandler(db, cfg, repo)
	redirectHandler := handlers.NewRedirectHandler(linkHandler, repo)

	metricsApp := fiber.New(fiber.Config{})
	prometheus := fiberprometheus.New("tikkin")
	prometheus.RegisterAt(metricsApp, "/metrics")
	go metricsApp.Listen(":3001")

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(prometheus.Middleware)

	if cfg.Docs.Enabled {
		enableDocs(app, cfg)
	}

	loginHandler := handlers.NewLoginHandler(cfg, db, userHandler)

	// Unauthenticated routes
	app.Post("/api/v1/auth/login", loginHandler.HandleLogin)
	app.Post("/api/v1/auth/signup", loginHandler.HandleRegister)
	app.Get("/api/v1/users/verify/:token", userHandler.HandleVerification)
	app.Get("/:slug", redirectHandler.HandleRedirect)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.Server.Jwt.Secret)},
	}))

	// Authenticated routes
	app.Post("/api/v1/links", linkHandler.HandleCreateLink)
	app.Put("/api/v1/links/:id", linkHandler.HandleUpdateLink)
	app.Delete("/api/v1/links/:id", linkHandler.HandleDeleteLink)
	app.Get("/api/v1/links", linkHandler.HandleGetLinks)

	if cfg.Links.DeleteExpired {
		expirationHandler := handlers.NewExpirationHandler(db, cfg, repo)
		go expirationHandler.ExpirationLoop()
	}

	app.Listen(":" + strconv.Itoa(cfg.Server.Port))
}

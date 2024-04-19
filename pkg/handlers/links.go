package handlers

import (
	"context"
	"errors"
	"github.com/aidarkhanov/nanoid/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/config"
	"tikkin/pkg/db"
	"tikkin/pkg/dto"
	"tikkin/pkg/model"
	"tikkin/pkg/repository"
	"tikkin/pkg/utils"
)

const DefaultAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var BLOCKED_SLUGS = []string{"admin", "api", "auth", "login", "logout", "register", "links", "users", "not_found"}

type LinkHandler struct {
	db               *db.DB
	Config           *config.Config
	repository       *repository.LinksRepository
	visitsRepository *repository.VisitsRepository
}

func NewLinkHandler(db *db.DB, config *config.Config) LinkHandler {
	repo := repository.NewLinksRepository(db)
	visitRepo := repository.NewVisitsRepository(db)
	return LinkHandler{db: db, Config: config, repository: &repo, visitsRepository: &visitRepo}
}

func (l *LinkHandler) generateSlug() string {
	result, err := nanoid.GenerateString(DefaultAlphabet, l.Config.Links.Length)
	if err != nil {
		log.Err(err).Msg("Failed to generate slug. Retrying...")
		return l.generateSlug()
	}
	if utils.Contains(BLOCKED_SLUGS, result) {
		log.Warn().Str("slug", result).Msg("Generated blocked slug. Retrying...")
		return l.generateSlug()
	}
	return result
}

func (l *LinkHandler) HandleCreateLink(c *fiber.Ctx) error {
	link := new(model.Link)

	if err := c.BodyParser(link); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if len(link.Description) > 1000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "description_too_long",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	log.Info().Msg("User: " + user.Claims.(jwt.MapClaims)["email"].(string))

	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)
	link.UserId = int64(userId)

	existingLink, err := l.repository.GetLinkBySlug(link.Slug)
	if err != nil {
		log.Err(err).Msg("Failed to check if slug exists")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal_error",
		})
	}
	if existingLink != nil {
		log.Warn().Str("slug", link.Slug).Msg("Slug already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "slug_exists",
		})
	}

	link, err = l.createLink(*link)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create link",
		})
	}

	return c.JSON(link)
}

func (l *LinkHandler) createLink(link model.Link) (*model.Link, error) {
	// Create a new link
	log.Info().
		Str("slug", link.Slug).
		Str("description", link.Description).
		Str("target_url", link.TargetUrl).
		Msg("Creating link")

	if link.Slug == "" {
		link.Slug = l.generateSlug()
	}

	if utils.Contains(BLOCKED_SLUGS, link.Slug) {
		log.Warn().Str("slug", link.Slug).Msg("Blocked slug")
		return nil, errors.New("slug_denied")
	}

	linkId := 0
	err := l.db.Pool.QueryRow(context.Background(),
		"INSERT INTO links (user_id, slug, description, expire_at, target_url) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		link.UserId, link.Slug, link.Description, nil, link.TargetUrl).Scan(&linkId)

	if err != nil {
		log.Err(err).Msg("Failed to create link")
		return nil, err
	}

	return l.repository.GetLink(linkId)

	// Save the link to the database
}

func (l *LinkHandler) HandleUpdateLink(ctx *fiber.Ctx) error {
	return ctx.SendString("Update link")
}

func (l *LinkHandler) HandleDeleteLink(ctx *fiber.Ctx) error {
	return ctx.SendString("Delete link")
}

func (l *LinkHandler) HandleGetLinks(ctx *fiber.Ctx) error {
	page := ctx.QueryInt("page")
	if page < 0 {
		page = 0
	}
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)
	if userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	links, err := l.repository.GetUserLinks(int64(userId), page)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get links",
		})
	}

	dtos := make([]dto.LinkDTO, len(links))
	for i, link := range links {

		visits := l.visitsRepository.CountVisits(link)
		dtos[i] = dto.LinkDTO{
			ID:          link.ID,
			Slug:        link.Slug,
			Description: link.Description,
			TargetUrl:   link.TargetUrl,
			CreatedAt:   link.CreatedAt,
			UpdatedAt:   link.UpdatedAt,
			ExpireAt:    link.ExpireAt,
			Visits:      visits,
		}
	}

	return ctx.JSON(dtos)
}

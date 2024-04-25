package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"net/url"
	"strconv"
	"tikkin/pkg/config"
	"tikkin/pkg/db"
	"tikkin/pkg/dto"
	"tikkin/pkg/model"
	"tikkin/pkg/repository"
	"time"
)

type LinkHandler struct {
	db               *db.DB
	Config           *config.Config
	repository       *repository.LinksRepository
	visitsRepository *repository.VisitsRepository
}

func NewLinkHandler(db *db.DB, config *config.Config, linksRepository repository.LinksRepository) LinkHandler {
	visitRepo := repository.NewVisitsRepository(db)
	return LinkHandler{db: db, Config: config, repository: &linksRepository, visitsRepository: &visitRepo}
}

func validateNewLink(link *model.Link) error {
	if len(link.Description) > 1000 {
		return fiber.NewError(fiber.StatusBadRequest, "description_too_long")
	}
	if len(link.Slug) > 25 {
		return fiber.NewError(fiber.StatusBadRequest, "slug_too_long")
	}
	if len(link.TargetUrl) > 1000 {
		return fiber.NewError(fiber.StatusBadRequest, "target_url_too_long")
	}
	if link.ExpireAt != nil && link.ExpireAt.Before(time.Now()) {
		return fiber.NewError(fiber.StatusBadRequest, "expire_at_invalid")
	}
	// validate valid URL
	_, err := url.ParseRequestURI(link.TargetUrl)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "target_url_invalid")
	}
	u, err := url.Parse(link.TargetUrl)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fiber.NewError(fiber.StatusBadRequest, "target_url_invalid")
	}
	return nil
}

// HandleCreateLink creates a new link
// @Summary Create a new link
// @Description Create a new link
// @Tags links
// @Accept json
// @Produce json
// @Param link body model.Link true "Link"
// @Success 200 {object} model.Link
// @Router /api/v1/links [post]
// @Security ApiKeyAuth
func (l *LinkHandler) HandleCreateLink(c *fiber.Ctx) error {
	link := new(model.Link)

	if err := c.BodyParser(link); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := validateNewLink(link); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
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

	link, err = l.repository.CreateLink(*link)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create link",
		})
	}

	return c.JSON(link)
}

// HandleUpdateLink updates a link
// @Summary Update a link
// @Description Update a link
// @Tags links
// @Accept json
// @Produce json
// @Param id path int true "Link ID"
// @Param link body model.Link true "Link"
// @Router /api/v1/links/{id} [put]
// @Success 200 {object} model.Link
func (l *LinkHandler) HandleUpdateLink(ctx *fiber.Ctx) error {

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)
	if userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	linkIdStr := ctx.Params("id")
	if linkIdStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id_required",
		})
	}

	linkId, err := strconv.ParseInt(linkIdStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id_invalid",
		})
	}

	body := new(model.Link)
	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := validateNewLink(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	body.UserId = int64(userId)

	link, err := l.repository.UpdateLink(linkId, *body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update link",
		})
	}

	return ctx.JSON(link)
}

// HandleDeleteLink deletes a link
// @Summary Delete a link
// @Description Delete a link
// @Tags links
// @Accept json
// @Produce json
// @Param id path int true "Link ID"
// @Success 200
// @Router /api/v1/links/{id} [delete]
// @Security ApiKeyAuth
func (l *LinkHandler) HandleDeleteLink(ctx *fiber.Ctx) error {

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)
	if userId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	linkIdStr := ctx.Params("id")
	if linkIdStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id_required",
		})
	}

	linkId, err := strconv.ParseInt(linkIdStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id_invalid",
		})
	}

	link, err := l.repository.GetLinkByID(linkId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get link",
		})
	}

	if link.UserId != int64(userId) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	err = l.repository.DeleteLink(linkId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete link",
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// HandleGetLinks returns all links for the authenticated user
// @Summary Get all links
// @Description Get all links for the authenticated user
// @Tags links
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {array} dto.LinkDTO
// @Router /api/v1/links [get]
// @Security ApiKeyAuth
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

	links, err := l.repository.GetUserLinks(int64(userId), int32(page))
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
			CreatedAt:   &link.CreatedAt,
			UpdatedAt:   &link.UpdatedAt,
			ExpireAt:    link.ExpireAt,
			Visits:      visits,
		}
	}

	return ctx.JSON(dtos)
}

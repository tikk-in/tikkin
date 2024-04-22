package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	"strings"
	"tikkin/pkg/model"
	"tikkin/pkg/repository"
)

type RedirectHandler struct {
	LinkHandler      *LinkHandler
	repository       repository.LinksRepository
	visitsRepository repository.VisitsRepository
}

func NewRedirectHandler(linkHandler LinkHandler) RedirectHandler {
	repo := repository.NewLinksRepository(linkHandler.db)
	visitsRepository := repository.NewVisitsRepository(linkHandler.db)
	return RedirectHandler{LinkHandler: &linkHandler, repository: repo, visitsRepository: visitsRepository}
}

func (r *RedirectHandler) handleVisit(link *model.Link, headers map[string][]string, realIP string) {

	userAgent := headers["User-Agent"]
	referrer := headers["Referer"]

	userAgentStr := strings.Join(userAgent, ",")
	referrerStr := strings.Join(referrer, ",")

	visit := model.Visits{
		ID:          uuid.NewString(),
		LinkID:      link.ID,
		UserAgent:   &userAgentStr,
		Referrer:    &referrerStr,
		CountryCode: nil,
	}

	r.visitsRepository.InsertVisit(visit)
}

func (r *RedirectHandler) HandleRedirect(c *fiber.Ctx) error {

	slug := utils.CopyString(c.Params("slug"))

	headers := c.GetReqHeaders()
	realIP := c.IP()
	targetLink, err := r.repository.GetLinkBySlug(slug)
	if err != nil || targetLink == nil {
		return c.Redirect("/not_found")
	}
	go r.handleVisit(targetLink, headers, realIP)
	return c.Redirect(targetLink.TargetUrl)
}

func (r *RedirectHandler) getLinkBySlug(slug string) string {
	link, err := r.LinkHandler.repository.GetLinkBySlug(slug)
	if err != nil || link == nil {
		return "not_found"
	}
	return link.TargetUrl
}

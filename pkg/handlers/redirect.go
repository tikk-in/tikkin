package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	"strings"
	"tikkin/pkg/model"
	"tikkin/pkg/repository"
	tikkinutils "tikkin/pkg/utils"
	"time"
)

type RedirectHandler struct {
	LinkHandler *LinkHandler
	repository  repository.Repository
}

func NewRedirectHandler(linkHandler LinkHandler, repository repository.Repository) RedirectHandler {
	return RedirectHandler{LinkHandler: &linkHandler, repository: repository}
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

	r.repository.InsertVisit(visit)
}

// HandleRedirect handles link redirection
func (r *RedirectHandler) HandleRedirect(c *fiber.Ctx) error {

	if tikkinutils.IsInvalidPath(c.Path()) {
		return c.SendStatus(fiber.StatusNotFound)
	}

	slug := utils.CopyString(c.Params("slug"))

	headers := c.GetReqHeaders()
	realIP := c.IP()
	targetLink, err := r.repository.GetLinkBySlug(slug)
	if err != nil || targetLink == nil {
		return c.Redirect("/not_found")
	}

	if targetLink.ExpireAt != nil && targetLink.ExpireAt.Before(time.Now()) {
		return c.Redirect("/not_found")
	}

	go r.handleVisit(targetLink, headers, realIP)
	return c.Redirect(targetLink.TargetUrl)
}

func (r *RedirectHandler) getLinkBySlug(slug string) string {
	link, err := r.repository.GetLinkBySlug(slug)
	if err != nil || link == nil {
		return "not_found"
	}
	return link.TargetUrl
}

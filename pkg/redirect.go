package pkg

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type RedirectHandler struct {
	LinkHandler *LinkHandler
}

func NewRedirectHandler(linkHandler LinkHandler) RedirectHandler {
	return RedirectHandler{LinkHandler: &linkHandler}
}

func (r *RedirectHandler) handleVisit(slug string, headers map[string][]string, realIP string) {

}

func (r *RedirectHandler) HandleRedirect(c *fiber.Ctx) error {

	slug := utils.CopyString(c.Params("slug"))

	headers := c.GetReqHeaders()
	realIP := c.IP()
	target := r.getLinkBySlug(slug)
	go r.handleVisit(target, headers, realIP)

	return c.Redirect(target)
}

func (r *RedirectHandler) getLinkBySlug(slug string) string {
	link, err := r.LinkHandler.GetLinkBySlug(slug)
	if err != nil || link == nil {
		return "not_found"
	}
	return link.TargetUrl
}

package utils

import (
	"github.com/aidarkhanov/nanoid/v2"
	"github.com/rs/zerolog/log"
)

const defaultAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var BlockedSlugs = []string{"admin", "api", "auth", "login", "logout", "register", "links", "users", "not_found"}

func GenerateSlug(length int) string {
	result, err := nanoid.GenerateString(defaultAlphabet, length)
	if err != nil {
		log.Err(err).Msg("Failed to generate slug. Retrying...")
		return GenerateSlug(length)
	}
	if Contains(BlockedSlugs, result) {
		log.Warn().Str("slug", result).Msg("Generated blocked slug. Retrying...")
		return GenerateSlug(length)
	}
	return result
}

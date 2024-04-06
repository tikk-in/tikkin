package pkg

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"tikkin/pkg/config"
	"tikkin/pkg/utils"
	"time"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	Config      *config.Config
	db          *DB
	UserHandler UserHandler
}

func NewLoginHandler(cfg *config.Config, db *DB, userHandler UserHandler) *AuthHandler {
	return &AuthHandler{Config: cfg, db: db, UserHandler: userHandler}
}

func (l *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	login := new(Login)
	err := c.BodyParser(login)
	if err != nil {
		log.Err(err).Msg("Cannot parse JSON")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := l.UserHandler.FindUserByEmail(login.Email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		log.Err(err).Msg("Failed to find user")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if !utils.DoesPasswordHashMatch(login.Password, user.Password) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"admin":   true,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(l.Config.Server.Jwt.Secret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

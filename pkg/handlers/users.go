package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"tikkin/pkg/config"
	"tikkin/pkg/db"
	"tikkin/pkg/email"
	"tikkin/pkg/model"
	"tikkin/pkg/repository"
)

type UserHandler struct {
	db           *db.DB
	config       *config.Config
	emailHandler *email.EmailHandler
	repository   repository.Repository
}

func NewUserHandler(db *db.DB,
	config *config.Config,
	emailHandler *email.EmailHandler,
	repository repository.Repository) UserHandler {
	return UserHandler{db: db, config: config, emailHandler: emailHandler, repository: repository}
}

func (u *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	return c.SendString("Get user")
}

func (u *UserHandler) FindUserByEmail(email string) (*model.User, error) {
	row := u.db.Pool.QueryRow(context.Background(), "SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1", email)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// HandleVerification handles user verification
// @Summary Verify user after registration
// @Accept json
// @Produce json
// @Tags users
// @Param token path string true "Verification token"
// @Success 200
// @Router /api/v1/users/verify/{token} [get]
func (u *UserHandler) HandleVerification(ctx *fiber.Ctx) error {
	token := ctx.Params("token")
	user, err := u.repository.FindUserByVerificationToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	_, err = u.repository.MarkUserAsVerified(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify user"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User verified"})
}

func (u *UserHandler) SignUpUser(user model.User) (*model.User, error) {
	createdUser, err := u.repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	if u.config.Email.Enabled {
		err = u.emailHandler.SendVerificationEmail(*createdUser)
		if err != nil {
			return nil, err
		}
	}

	return createdUser, nil
}

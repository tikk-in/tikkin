package pkg

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"tikkin/pkg/config"
	"tikkin/pkg/email"
	"tikkin/pkg/model"
)

type UserHandler struct {
	db           *DB
	config       *config.Config
	emailHandler *email.EmailHandler
}

func NewUserHandler(db *DB, config *config.Config, emailHandler *email.EmailHandler) UserHandler {
	return UserHandler{db: db, config: config, emailHandler: emailHandler}
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

func (u *UserHandler) CreateUser(user model.User) (*model.User, error) {

	if user.VerificationToken == "" {
		verificationToken, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		user.VerificationToken = verificationToken.String()
	}

	row := u.db.Pool.QueryRow(context.Background(), "INSERT INTO users (email, password, verified, verification_token) VALUES ($1, $2, $3, $4) RETURNING id, email, password, verified, verification_token, created_at, updated_at",
		user.Email, user.Password, user.Verified, user.VerificationToken)

	createdUser := model.User{}
	err := row.Scan(&createdUser.ID, &createdUser.Email, &createdUser.Password,
		&createdUser.Verified, &createdUser.VerificationToken,
		&createdUser.CreatedAt, &createdUser.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if u.config.Email.Enabled {
		err = u.emailHandler.SendVerificationEmail(createdUser)
		if err != nil {
			return nil, err
		}
	}

	return &createdUser, nil
}

func (u *UserHandler) HandleVerification(ctx *fiber.Ctx) error {
	token := ctx.Params("token")
	user, err := u.FindUserByVerificationToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	_, err = u.MarkUserAsVerified(user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify user"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User verified"})
}

func (u *UserHandler) FindUserByVerificationToken(token string) (*model.User, error) {
	row := u.db.Pool.QueryRow(context.Background(), "SELECT id, email, password, verified, verification_token, created_at, updated_at FROM users WHERE verification_token = $1", token)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Verified, &user.VerificationToken, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserHandler) MarkUserAsVerified(user *model.User) (*model.User, error) {
	if user.Verified {
		return user, nil
	}
	if user.ID == 0 {
		return nil, errors.New("user.not.found")
	}

	res, err := u.db.Pool.Exec(context.Background(), "UPDATE users SET verified = $1, verification_token = null WHERE id = $2", true, user.ID)
	if err != nil {
		return nil, err
	}
	if res.RowsAffected() == 0 {
		return nil, errors.New("user.not.found")
	}

	return u.FindUserByID(user.ID)
}

func (u *UserHandler) FindUserByID(id int64) (*model.User, error) {
	row := u.db.Pool.QueryRow(context.Background(), "SELECT id, email, password, verified, created_at, updated_at FROM users WHERE id = $1", id)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

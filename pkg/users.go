package pkg

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"tikkin/pkg/model"
)

type UserHandler struct {
	db *DB
}

func NewUserHandler(db *DB) UserHandler {
	return UserHandler{db: db}
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

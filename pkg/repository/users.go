package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"tikkin/pkg/db"
	"tikkin/pkg/model"
)

type UsersRepository struct {
	db *db.DB
}

func NewUsersRepository(db *db.DB) UsersRepository {
	return UsersRepository{db: db}
}

func (u *UsersRepository) FindUserByVerificationToken(token string) (*model.User, error) {
	row := u.db.Pool.QueryRow(context.Background(), "SELECT id, email, password, verified, verification_token, created_at, updated_at FROM users WHERE verification_token = $1", token)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Verified, &user.VerificationToken, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersRepository) MarkUserAsVerified(user *model.User) (*model.User, error) {
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

func (u *UsersRepository) FindUserByID(id int64) (*model.User, error) {
	row := u.db.Pool.QueryRow(context.Background(), "SELECT id, email, password, verified, created_at, updated_at FROM users WHERE id = $1", id)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UsersRepository) CreateUser(user model.User) (*model.User, error) {

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

	return &createdUser, nil
}

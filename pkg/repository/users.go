package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"tikkin/pkg/model"
	"tikkin/pkg/repository/queries"
)

type UserRepository interface {
	FindUserByVerificationToken(token string) (*model.User, error)
}

func (r *Repository) FindUserByVerificationToken(token string) (*model.User, error) {
	if token == "" {
		return nil, errors.New("token.empty")
	}
	user, err := r.db.Queries(context.Background()).FindUserByVerificationToken(context.Background(), &token)
	if err != nil {
		return nil, err
	}

	pUser := model.User(user)
	return &pUser, nil
}

func (r *Repository) MarkUserAsVerified(user *model.User) (*model.User, error) {
	if user.Verified {
		return user, nil
	}
	if user.ID == 0 {
		return nil, errors.New("user.not.found")
	}

	updatedUser, err := r.db.Queries(context.Background()).MarkUserAsVerified(context.Background(), user.ID)
	if err != nil {
		return nil, err
	}
	pUser := model.User(updatedUser)
	return &pUser, nil
}

func (r *Repository) FindUserByID(id int64) (*model.User, error) {

	user, err := r.db.Queries(context.Background()).FindUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	cUser := model.User(user)
	return &cUser, nil
}

func (r *Repository) CreateUser(user model.User) (*model.User, error) {

	if user.VerificationToken == nil || *user.VerificationToken == "" {
		verificationToken, err := uuid.NewRandom()
		if err != nil {
			return nil, err
		}
		tokenStr := verificationToken.String()
		user.VerificationToken = &tokenStr
	}

	result, err := r.db.Queries(context.Background()).CreateUser(context.Background(),
		queries.CreateUserParams{
			Email:             user.Email,
			Password:          user.Password,
			Verified:          user.Verified,
			VerificationToken: user.VerificationToken,
		})
	if err != nil {
		return nil, err
	}

	createdUser := model.User(result)
	return &createdUser, nil
}

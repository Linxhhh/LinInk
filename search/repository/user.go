package repository

import (
	"context"

	"github.com/Linxhhh/LinInk/search/domain"
	"github.com/Linxhhh/LinInk/search/repository/dao"
)

type UserRepository interface {
	PutUser(ctx context.Context, user domain.User) error
	SearchUser(ctx context.Context, expression string) ([]domain.User, error)
}

type userRepository struct {
	userDAO dao.UserElasticDAO
}

func NewUserRepository(userDAO dao.UserElasticDAO) UserRepository {
	return &userRepository{
		userDAO: userDAO,
	}
}

func (repo *userRepository) PutUser(ctx context.Context, user domain.User) error {
	return repo.userDAO.InputUser(ctx, dao.User{
		Id:           user.Id,
		NickName:     user.NickName,
		Introduction: user.Introduction,
	})
}

func (repo *userRepository) SearchUser(ctx context.Context, expression string) ([]domain.User, error) {
	users, err := repo.userDAO.SearchUser(ctx, expression)
	if err != nil {
		return nil, err
	}

	domainUsers := make([]domain.User, 0, len(users))
	for _, user := range users {
		domainUser := domain.User{
			Id:           user.Id,
			NickName:     user.NickName,
			Introduction: user.Introduction,
		}
		domainUsers = append(domainUsers, domainUser)
	}
	return domainUsers, nil
}

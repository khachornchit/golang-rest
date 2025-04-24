package ports

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang-rest/internal/core/domain"
)

type UserRepositoryInterface interface {
	CreateUser(user *domain.User) error
	GetAllUsers() ([]domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUserLoginByEmail(email string) (*domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	UpdateUserByID(id string, updates bson.M) (*domain.User, error)
	DeleteUserByID(id string) error
}

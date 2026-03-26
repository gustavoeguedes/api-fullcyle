package database

import "github.com/gustavoeguedes/api-fullcycle/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

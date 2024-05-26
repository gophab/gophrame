package service

import "github.com/gophab/gophrame/service/dto"

type UserService interface {
	CreateUser(user *dto.User) (*dto.User, error)
	GetById(id string) (*dto.User, error)
	LoadAllPolicy() error
}

func GetUserService() UserService {
	return _services.UserService
}

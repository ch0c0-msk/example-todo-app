package service

import (
	todo "github.com/ch0c0-msk/example-todo-app"
	repostiory "github.com/ch0c0-msk/example-todo-app/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repostiory.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(*repos),
	}
}

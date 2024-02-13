package repository

import (
	todo "github.com/ch0c0-msk/example-todo-app"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUserByUsernameAndPassword(username string, password string) (todo.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSql(db),
	}
}

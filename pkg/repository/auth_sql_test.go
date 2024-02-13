package repository

import (
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	todo "github.com/ch0c0-msk/example-todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAuthSql_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err.Error())
	}

	dbX := sqlx.NewDb(db, "postgres")
	defer dbX.Close()

	auth := NewAuthSql(dbX)

	testCases := []struct {
		name    string
		mock    func()
		input   todo.User
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("Test", "test", "password").WillReturnRows(rows)
			},
			input: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "password",
			},
			want: 1,
		},
		{
			name: "Empty Fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO users").
					WithArgs("Test", "test", "").WillReturnRows(rows)
			},
			input: todo.User{
				Name:     "Test",
				Username: "test",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			got, err := auth.CreateUser(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

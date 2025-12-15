package postgresql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func TestStorage_LoginUser(t *testing.T) {
	type args struct {
		login model.Login
	}
	tests := []struct {
		name    string
		args    args
		query   string
		wantErr bool
	}{
		{
			name: "simple test #1",
			args: args{
				login: model.Login{Login: "1", Password: ""},
			},
			query: `
				SELECT login, password FROM users
				WHERE login = $1 
			`,
		},
		{
			name: "simple test #2",
			args: args{
				login: model.Login{Login: "1", Password: ""},
			},
			query: `
				SELECT login, password FROM users
				WHERE login = $2 
			`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			s := &Storage{
				db: db,
			}

			rows := sqlmock.NewRows([]string{"login", "password"}).
				AddRow("1", "2")

			mock.ExpectBegin()
			mock.ExpectQuery(tt.query).
				WithArgs(tt.args.login.Login).
				WillReturnRows(rows)
			mock.ExpectCommit()

			_, err := s.LoginUser(context.Background(), tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.LoginUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

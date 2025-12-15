package postgresql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStorage_GetBalance(t *testing.T) {
	type args struct {
		login string
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
				login: "1",
			},
			query: `
				SELECT current, withdrawn FROM users
				WHERE login = $1 
			`,
		},
		{
			name: "simple test #2",
			args: args{
				login: "1",
			},
			query: `
				SELECT current, withdrawn FROM users
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

			rows := sqlmock.NewRows([]string{"current", "withdrawn"}).
				AddRow(1, 2)

			mock.ExpectBegin()
			mock.ExpectQuery(tt.query).
				WithArgs(tt.args.login).
				WillReturnRows(rows)
			mock.ExpectCommit()

			_, err := s.GetBalance(context.Background(), tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorage_ListWithdrawals(t *testing.T) {
	type args struct {
		login string
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
				login: "1",
			},
			query: `
				SELECT number, sum, processed_at FROM orders
				WHERE login = $1 AND sum > 0
			`,
		},
		{
			name: "simple test #2",
			args: args{
				login: "1",
			},
			query: `
				SELECT number, sum, processed_at FROM orders
				WHERE login = $2 AND sum > 0
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

			rows := sqlmock.NewRows([]string{"number", "sum", "processed_at"}).
				AddRow("1", 2, time.Now())

			mock.ExpectBegin()
			mock.ExpectQuery(tt.query).
				WithArgs(tt.args.login).
				WillReturnRows(rows)
			mock.ExpectCommit()

			_, err := s.ListWithdrawals(context.Background(), tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ListWithdrawals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

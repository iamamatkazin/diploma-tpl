package postgresql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iamamatkazin/diploma-tpl/internal/gophermart/model"
)

func TestStorage_LoadOrder(t *testing.T) {
	type args struct {
		ctx   context.Context
		login string
		order string
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
				order: "1",
			},
			query: `
				INSERT INTO orders (number, login, status, uploaded_at)
				VALUES ($1, $2, $3, $4)
			`,
		},
		{
			name: "simple test #2",
			args: args{
				login: "1",
				order: "1",
			},
			query: `
				INSERT INTO orders (number, login, status, uploaded_at)
				VALUES ($1)
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

			mock.ExpectBegin()
			mock.ExpectExec(tt.query).
				WithArgs(tt.args.order, tt.args.login, model.New, sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			_, err := s.LoadOrder(context.Background(), tt.args.login, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.LoadOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorage_ListOrders(t *testing.T) {
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
				SELECT number, status, accrual, uploaded_at FROM orders
				WHERE login = $1 
			`,
		},
		{
			name: "simple test #2",
			args: args{
				login: "1",
			},
			query: `
				SELECT number, status, accrual, uploaded_at FROM orders
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

			rows := sqlmock.NewRows([]string{"number", "status", "accrual", "uploaded_at"}).
				AddRow("1", "new", 100.50, time.Now())

			mock.ExpectBegin()
			mock.ExpectQuery(tt.query).
				WithArgs(tt.args.login).
				WillReturnRows(rows)
			mock.ExpectCommit()

			_, err := s.ListOrders(context.Background(), tt.args.login)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

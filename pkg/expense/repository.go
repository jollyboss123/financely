package expense

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/jollyboss123/finance-tracker/pkg/server/message"
	s "github.com/shopspring/decimal"
)

type Expense interface {
	Create(ctx context.Context, request *CreateRequest) (int, error)
	List(ctx context.Context) ([]*Schema, error)
	Read(ctx context.Context, expenseID int) (*Schema, error)
	Update(ctx context.Context, request *UpdateRequest) error
	Delete(ctx context.Context, expenseID int) error
	Total(ctx context.Context) (s.Decimal, error)
}

type expenseRepository struct {
	db *sqlx.DB
}

const (
	InsertIntoExpenses = "insert into expenses (title, amount) values ($1, $2) returning id"
	SelectFromExpenses = "select * from expenses order by created_at desc"
	SelectExpenseByID  = "select * from expenses where id = $1"
	UpdateExpense      = "update expenses set title = $1, amount = $2 where id = $3 returning id"
	DeleteExpenseByID  = "delete from expenses where id = $1 returning id"
	TotalExpenses      = "select sum(amount) from expenses"
)

var (
	ErrFetchingExpenses = errors.New("error fetching expenses")
)

func New(db *sqlx.DB) *expenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) Create(ctx context.Context, request *CreateRequest) (expenseId int, err error) {
	if err = r.db.QueryRowContext(ctx, InsertIntoExpenses, request.Title, request.Amount).Scan(&expenseId); err != nil {
		return 0, errors.New("repository.Expense.Create")
	}

	return expenseId, nil
}

func (r *expenseRepository) List(ctx context.Context) ([]*Schema, error) {
	var expenses []*Schema
	err := r.db.SelectContext(ctx, &expenses, SelectFromExpenses)
	if err != nil {
		return nil, ErrFetchingExpenses
	}

	return expenses, nil
}

func (r *expenseRepository) Read(ctx context.Context, expenseID int) (*Schema, error) {
	var expense Schema
	err := r.db.GetContext(ctx, &expense, SelectExpenseByID, expenseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, message.ErrBadRequest
		}
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepository) Update(ctx context.Context, request *UpdateRequest) error {
	var returnedID int

	err := r.db.QueryRowContext(ctx, UpdateExpense, request.Title, request.Amount, request.ID).Scan(&returnedID)
	if err != nil {
		return err
	}

	return nil
}

func (r *expenseRepository) Delete(ctx context.Context, expenseID int) error {
	var returnedID int

	_, err := r.Read(ctx, expenseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return message.ErrNoRecord
		}
	}

	err = r.db.QueryRowContext(ctx, DeleteExpenseByID, expenseID).Scan(&returnedID)
	if err != nil {
		return err
	}

	return nil
}

func (r *expenseRepository) Total(ctx context.Context) (s.Decimal, error) {
	var total s.Decimal
	err := r.db.GetContext(ctx, &total, TotalExpenses)
	if err != nil {
		return s.NewFromInt(0), err
	}
	return total, nil
}

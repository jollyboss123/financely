package expense

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jollyboss123/finance-tracker/pkg/server/message"
	s "github.com/shopspring/decimal"
)

type Expense interface {
	Create(ctx context.Context, request *CreateRequest) (int, error)
	List(ctx context.Context, filter *Filter) ([]*Schema, error)
	Read(ctx context.Context, expenseID int) (*Schema, error)
	Update(ctx context.Context, request *UpdateRequest) error
	Delete(ctx context.Context, expenseID int) error
	Total(ctx context.Context, filter *Filter) (s.Decimal, error)
	Search(ctx context.Context, filter *Filter) ([]*Schema, error)
}

type expenseRepository struct {
	db *sqlx.DB
}

const (
	InsertIntoExpenses         = "insert into expenses (title, amount, transaction_date) values ($1, $2, $3) returning id"
	SelectFromExpenses         = "select * from expenses order by transaction_date desc"
	SelectFromExpensesPaginate = "select * from expenses order by transaction_date desc limit $1 offset $2"
	SelectExpenseByID          = "select * from expenses where id = $1"
	UpdateExpense              = "update expenses set title = $1, amount = $2 where id = $3 returning id"
	DeleteExpenseByID          = "delete from expenses where id = $1 returning id"
	TotalExpensesDynamic       = "select sum(amount) from expenses where "
	SearchExpensesDynamic      = "select * from expenses where title like '%' || $1 || '%' "
)

var (
	ErrFetchingExpenses = errors.New("error fetching expenses")
)

func New(db *sqlx.DB) *expenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) Create(ctx context.Context, request *CreateRequest) (expenseId int, err error) {
	if err = r.db.QueryRowContext(ctx, InsertIntoExpenses, request.Title, request.Amount, request.TransactionDate).Scan(&expenseId); err != nil {
		return 0, errors.New("repository.Expense.Create")
	}

	return expenseId, nil
}

func (r *expenseRepository) List(ctx context.Context, filter *Filter) ([]*Schema, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	var expenses []*Schema
	if filter.Pagination.DisablePaging {
		err := r.db.SelectContext(ctx, &expenses, SelectFromExpenses)
		if err != nil {
			return nil, ErrFetchingExpenses
		}
	} else {
		err := r.db.SelectContext(ctx, &expenses, SelectFromExpensesPaginate, filter.Pagination.Limit, filter.Pagination.Offset)
		if err != nil {
			return nil, ErrFetchingExpenses
		}
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

func (r *expenseRepository) Total(ctx context.Context, filter *Filter) (s.Decimal, error) {
	if filter == nil {
		return s.NewFromInt(0), errors.New("filter cannot be nil")
	}

	baseQuery := TotalExpensesDynamic
	var sqlParams []interface{}
	paramCounter := 0

	if filter.Year != "" {
		paramCounter++
		baseQuery += fmt.Sprintf("extract(year from transaction_date) = $%d ", paramCounter)
		sqlParams = append(sqlParams, filter.Year)
	}

	if filter.Month != "" {
		paramCounter++
		baseQuery += fmt.Sprintf("and extract(month from transaction_date) = $%d ", paramCounter)
		sqlParams = append(sqlParams, filter.Month)
	}

	if filter.Day != "" {
		paramCounter++
		baseQuery += fmt.Sprintf("and extract(day from transaction_date) = $%d", paramCounter)
		sqlParams = append(sqlParams, filter.Day)
	}

	var total s.Decimal
	err := r.db.GetContext(ctx, &total, baseQuery, sqlParams...)
	if err != nil {
		return s.NewFromInt(0), err
	}
	return total, nil
}

func (r *expenseRepository) Search(ctx context.Context, filter *Filter) ([]*Schema, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	baseQuery := SearchExpensesDynamic
	sqlParams := []interface{}{filter.Title}
	paramCounter := 1

	if filter.Year != "" {
		paramCounter++
		baseQuery += fmt.Sprintf("and extract(year from transaction_date) = $%d ", paramCounter)
		sqlParams = append(sqlParams, filter.Year)
	}

	if filter.Month != "" {
		paramCounter++
		baseQuery += fmt.Sprintf("and extract(month from transaction_date) = $%d ", paramCounter)
		sqlParams = append(sqlParams, filter.Month)
	}

	if filter.Day != "" {
		paramCounter++
		baseQuery += fmt.Sprintf("and extract(day from transaction_date) = $%d ", paramCounter)
		sqlParams = append(sqlParams, filter.Day)
	}

	paramCounter++
	baseQuery += fmt.Sprintf("order by transaction_date desc limit $%d ", paramCounter)
	sqlParams = append(sqlParams, filter.Pagination.Limit)

	paramCounter++
	baseQuery += fmt.Sprintf("offset $%d", paramCounter)
	sqlParams = append(sqlParams, filter.Pagination.Offset)

	var expenses []*Schema
	err := r.db.SelectContext(ctx, &expenses, baseQuery, sqlParams...)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

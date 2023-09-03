package expense

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jollyboss123/finance-tracker/pkg/server/message"
	"github.com/shopspring/decimal"
	"strings"
)

type Expense interface {
	Create(ctx context.Context, request *CreateRequest) (uuid.UUID, error)
	List(ctx context.Context, filter *Filter) ([]*Schema, error)
	Read(ctx context.Context, expenseID uuid.UUID) (*Schema, error)
	Update(ctx context.Context, request *UpdateRequest) error
	Delete(ctx context.Context, expenseID uuid.UUID) error
	Total(ctx context.Context, filter *Filter) (int64, error)
	Search(ctx context.Context, filter *Filter) ([]*Schema, error)
	Average(ctx context.Context, filter *Filter) (int64, error)
}

type expenseRepository struct {
	db *sqlx.DB
}

const (
	InsertIntoExpenses = `insert into expenses 
(title, amount_ud, currency_id_ud, currency_code_ud, amount_base, currency_id_base, currency_code_base, transaction_date) 
values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`
	SelectFromExpenses         = "select * from expenses order by transaction_date desc"
	SelectFromExpensesPaginate = "select * from expenses order by transaction_date desc limit $1 offset $2"
	SelectExpenseByID          = "select * from expenses where id = $1"
	UpdateExpense              = `update expenses 
set title = $1, amount_ud = $2, currency_id_ud = $3, currency_code_ud = $4, 
amount_base = $5, currency_id_base = $6, currency_code_base = $7, transaction_date = $8 
where id = $9 
returning id`
	DeleteExpenseByID      = "delete from expenses where id = $1 returning id"
	TotalExpensesDynamic   = "select coalesce(sum(amount_base), 0) from expenses where "
	SearchExpensesDynamic  = "select * from expenses where title like '%' || $1 || '%' "
	AverageExpensesDynamic = "select coalesce(avg(amount_base), 0) from expenses where "
)

var (
	ErrFetchingExpenses = errors.New("error fetching expenses")
)

func New(db *sqlx.DB) *expenseRepository {
	return &expenseRepository{db: db}
}

func (er *expenseRepository) Create(ctx context.Context, request *CreateRequest) (expenseId uuid.UUID, err error) {
	if err = er.db.QueryRowContext(ctx, InsertIntoExpenses,
		request.Title,
		request.Amount,
		request.CurrencyID,
		request.CurrencyCode,
		request.BaseAmount,
		request.BaseCurrencyID,
		request.BaseCurrencyCode,
		request.TransactionDate).Scan(&expenseId); err != nil {
		return uuid.Nil, errors.New("repository.Expense.Create")
	}

	return expenseId, nil
}

func (er *expenseRepository) List(ctx context.Context, filter *Filter) ([]*Schema, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	var expenses []*Schema
	if filter.Pagination.DisablePaging {
		err := er.db.SelectContext(ctx, &expenses, SelectFromExpenses)
		if err != nil {
			return nil, ErrFetchingExpenses
		}
	} else {
		err := er.db.SelectContext(ctx, &expenses, SelectFromExpensesPaginate, filter.Pagination.Limit, filter.Pagination.Offset)
		if err != nil {
			return nil, ErrFetchingExpenses
		}
	}
	return expenses, nil
}

func (er *expenseRepository) Read(ctx context.Context, expenseID uuid.UUID) (*Schema, error) {
	var expense Schema
	err := er.db.GetContext(ctx, &expense, SelectExpenseByID, expenseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, message.ErrBadRequest
		}
		return nil, err
	}
	return &expense, nil
}

func (er *expenseRepository) Update(ctx context.Context, request *UpdateRequest) error {
	var returnedID uuid.UUID

	err := er.db.QueryRowContext(ctx, UpdateExpense,
		request.Title,
		request.Amount,
		request.CurrencyID,
		request.CurrencyCode,
		request.BaseAmount,
		request.BaseCurrencyID,
		request.BaseCurrencyCode,
		request.TransactionDate,
		request.ID).Scan(&returnedID)
	if err != nil {
		return err
	}

	return nil
}

func (er *expenseRepository) Delete(ctx context.Context, expenseID uuid.UUID) error {
	var returnedID uuid.UUID

	_, err := er.Read(ctx, expenseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return message.ErrNoRecord
		}
		return err
	}

	err = er.db.QueryRowContext(ctx, DeleteExpenseByID, expenseID).Scan(&returnedID)
	if err != nil {
		return err
	}

	return nil
}

func (er *expenseRepository) Total(ctx context.Context, filter *Filter) (int64, error) {
	if filter == nil {
		return 0, errors.New("filter cannot be nil")
	}

	var clauses []string
	var sqlParams []interface{}
	paramCounter := 0

	if filter.Year == "" && filter.Month == "" && filter.Day == "" {
		clauses = append(clauses, "transaction_date::date = current_date")
	} else {
		dateFilters := []struct {
			Field  string
			Column string
		}{
			{filter.Year, "year"},
			{filter.Month, "month"},
			{filter.Day, "day"},
		}

		for _, dateFilter := range dateFilters {
			if dateFilter.Field != "" {
				paramCounter++
				clause := fmt.Sprintf("extract(%s from transaction_date) = $%d", dateFilter.Column, paramCounter)
				clauses = append(clauses, clause)
				sqlParams = append(sqlParams, dateFilter.Field)
			}
		}
	}

	baseQuery := TotalExpensesDynamic + strings.Join(clauses, " and ")

	var total int64
	err := er.db.GetContext(ctx, &total, baseQuery, sqlParams...)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (er *expenseRepository) Search(ctx context.Context, filter *Filter) ([]*Schema, error) {
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
	err := er.db.SelectContext(ctx, &expenses, baseQuery, sqlParams...)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func (er *expenseRepository) Average(ctx context.Context, filter *Filter) (int64, error) {
	if filter == nil {
		return 0, errors.New("filter cannot be nil")
	}

	var clauses []string
	var sqlParams []interface{}
	paramCounter := 0

	if filter.Year == "" && filter.Month == "" && filter.Day == "" {
		clauses = append(clauses, "transaction_date::date = current_date")
	} else {
		dateFilters := []struct {
			Field  string
			Column string
		}{
			{filter.Year, "year"},
			{filter.Month, "month"},
			{filter.Day, "day"},
		}

		for _, dateFilter := range dateFilters {
			if dateFilter.Field != "" {
				paramCounter++
				clause := fmt.Sprintf("extract(%s from transaction_date) = $%d", dateFilter.Column, paramCounter)
				clauses = append(clauses, clause)
				sqlParams = append(sqlParams, dateFilter.Field)
			}
		}
	}

	baseQuery := AverageExpensesDynamic + strings.Join(clauses, " and ")

	var avg decimal.Decimal
	err := er.db.GetContext(ctx, &avg, baseQuery, sqlParams...)
	if err != nil {
		return 0, err
	}
	return avg.IntPart(), nil
}

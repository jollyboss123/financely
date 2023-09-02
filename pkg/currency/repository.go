package currency

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jollyboss123/finance-tracker/pkg/server/message"
)

type Currency interface {
	Create(ctx context.Context, request *CreateRequest) (uuid.UUID, error)
	List(ctx context.Context, filter *Filter) ([]*Schema, error)
	Read(ctx context.Context, currencyID uuid.UUID) (*Schema, error)
	ReadByCode(ctx context.Context, code string) (uuid.UUID, error)
	Update(ctx context.Context, request *UpdateRequest) error
	Delete(ctx context.Context, currencyID uuid.UUID) error
	Search(ctx context.Context, filter *Filter) ([]*Schema, error)
}

type currencyRepository struct {
	db *sqlx.DB
}

const (
	InsertIntoCurrency         = "insert into currency (code, numeric_code, fraction, grapheme, template, decimal, thousand) values ($1, $2, $3, $4, $5, $6, $7) returning id"
	SelectFromCurrency         = "select * from currency order by code asc"
	SelectFromCurrencyPaginate = "select * from currency order by code asc limit $1 offset $2"
	SelectCurrencyByID         = "select * from currency where id = $1"
	SelectCurrencyByCode       = "select id from currency where code = $1 limit 1"
	UpdateCurrency             = "update currency set code = $1, numeric_code = $2, fraction = $3, grapheme = $4, template = $5, decimal = $6, thousand = $7 where id = $8 returning id"
	DeleteCurrencyByID         = "delete from currency where id = $1 returning id"
	SearchCurrencyPaginate     = "select * from currency where code like '%' || '%' || $1 || '%' || '%' order by code asc limit $2 offset $3"
)

var (
	ErrFetchingCurrency = errors.New("error fetching currencies")
)

func New(db *sqlx.DB) *currencyRepository {
	return &currencyRepository{db: db}
}

func (cr *currencyRepository) Create(ctx context.Context, request *CreateRequest) (currencyID uuid.UUID, err error) {
	if err = cr.db.QueryRowContext(ctx, InsertIntoCurrency,
		request.Code,
		request.NumericCode,
		request.Fraction,
		request.Grapheme,
		request.Template,
		request.Decimal,
		request.Thousand).Scan(&currencyID); err != nil {
		return uuid.Nil, errors.New("repository.Currency.Create")
	}

	return currencyID, nil
}

func (cr *currencyRepository) List(ctx context.Context, filter *Filter) ([]*Schema, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	var currencies []*Schema
	if filter.Pagination.DisablePaging {
		err := cr.db.SelectContext(ctx, &currencies, SelectFromCurrency)
		if err != nil {
			return nil, ErrFetchingCurrency
		}
	} else {
		err := cr.db.SelectContext(ctx, &currencies, SelectFromCurrencyPaginate)
		if err != nil {
			return nil, ErrFetchingCurrency
		}
	}
	return currencies, nil
}

func (cr *currencyRepository) Read(ctx context.Context, currencyID uuid.UUID) (*Schema, error) {
	var currency Schema
	err := cr.db.GetContext(ctx, &currency, SelectCurrencyByID, currencyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, message.ErrBadRequest
		}
		return nil, err
	}
	return &currency, nil
}

func (cr *currencyRepository) ReadByCode(ctx context.Context, code string) (uuid.UUID, error) {
	var cID uuid.UUID
	err := cr.db.QueryRowContext(ctx, SelectCurrencyByCode, code).Scan(&cID)
	if err != nil {
		return uuid.Nil, errors.New("repository.Currency.Read")
	}
	return cID, nil
}

func (cr *currencyRepository) Update(ctx context.Context, request *UpdateRequest) error {
	var returnedID uuid.UUID

	err := cr.db.QueryRowContext(ctx, UpdateCurrency,
		request.Code,
		request.NumericCode,
		request.Fraction,
		request.Grapheme,
		request.Template,
		request.Decimal,
		request.Thousand,
		request.ID).Scan(&returnedID)

	if err != nil {
		return err
	}

	return nil
}

func (cr *currencyRepository) Delete(ctx context.Context, currencyID uuid.UUID) error {
	var returnedID uuid.UUID

	_, err := cr.Read(ctx, currencyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return message.ErrNoRecord
		}
		return err
	}

	err = cr.db.QueryRowContext(ctx, DeleteCurrencyByID, currencyID).Scan(&returnedID)
	if err != nil {
		return err
	}

	return nil
}

func (cr *currencyRepository) Search(ctx context.Context, filter *Filter) ([]*Schema, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	var currencies []*Schema
	err := cr.db.SelectContext(ctx, &currencies, SearchCurrencyPaginate,
		filter.Code,
		filter.Pagination.Limit,
		filter.Pagination.Offset)
	if err != nil {
		return nil, err
	}
	return currencies, nil
}

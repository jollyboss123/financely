package rate

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Rate interface {
	Read(ctx context.Context, dest string) (float64, error)
	Create(ctx context.Context, base, dest string, rate float64) error
}

type rateRepository struct {
	db *sqlx.DB
}

const (
	//	SelectFromRate = `select e.rate
	//from exchange e
	//join currency bc on e.currency_id = bc.id
	//join currency dc on e.currency_id = dc.id
	//where bc.code = $1
	//and dc.code = $2`
	SelectFromRate = `select e.rate
from exchange e
join currency c on e.to_currency_id = c.id
where c.code = $1`
	CreateRate = `insert into exchange (from_currency_id, to_currency_id, rate)
select
    (select id from currency where code = $1),
    (select id from currency where code = $2),
    $3
on conflict (from_currency_id, to_currency_id)
do update set rate = excluded.rate returning rate`
)

func New(db *sqlx.DB) *rateRepository {
	return &rateRepository{db: db}
}

//	func (rr *rateRepository) Read(ctx context.Context, base, dest string) (float64, error) {
//		var r float64
//		err := rr.db.QueryRowContext(ctx, SelectFromRate, base, dest).Scan(&r)
//		if err != nil {
//			return 0, err
//		}
//		return r, nil
//	}
func (rr *rateRepository) Read(ctx context.Context, dest string) (float64, error) {
	var r float64
	err := rr.db.QueryRowContext(ctx, SelectFromRate, dest).Scan(&r)
	if err != nil {
		return 0, err
	}
	return r, nil
}

func (rr *rateRepository) Create(ctx context.Context, base, dest string, rate float64) error {
	var r float64
	if err := rr.db.QueryRowContext(ctx, CreateRate, base, dest, rate).Scan(&r); err != nil {
		return err
	}
	return nil
}

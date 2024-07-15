package repository

import (
	"context"

	"github.com/fahmyabida/labbaika/internal/app/domain"
	pkgErrors "github.com/fahmyabida/labbaika/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type LedgerRepository struct {
	DB *gorm.DB
}

func NewLedgerRepository(db *gorm.DB) domain.ILedgerRepo {
	return &LedgerRepository{
		DB: db,
	}
}

func (r LedgerRepository) Create(ctx context.Context, data *domain.Ledger) (err error) {
	dbResult := r.DB.WithContext(ctx).Create(data)
	if dbResult.Error != nil {
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		postgresError, ok := dbResult.Error.(*pgconn.PgError)
		if ok && postgresError.Code == "23505" {
			return pkgErrors.DuplicateLedgerError(pkgErrors.ErrDuplicateLedger)
		}
		return dbResult.Error
	}

	return
}

func (r LedgerRepository) FindByID(ctx context.Context, id string) (order domain.Ledger, err error) {
	dbResult := r.DB.Model(&order).Where("id = ?", id).Find(&order)
	if dbResult.RowsAffected == 0 {
		return order, pkgErrors.LedgerNotFoundError(pkgErrors.ErrLedgerNotFound)
	}
	return order, dbResult.Error
}

func (r LedgerRepository) GetLedgers(ctx context.Context, parameters interface{}) (result []domain.Ledger, err error) {
	dbResult := r.DB.Model(&result).Where("param = ?", parameters).Find(&result)
	if dbResult.RowsAffected == 0 {
		return result, pkgErrors.LedgerNotFoundError(pkgErrors.ErrLedgerNotFound)
	} else if err = dbResult.Error; err != nil {
		return result, err
	}
	return result, nil
}

func (r LedgerRepository) Update(ctx context.Context, order *domain.Ledger) error {
	dbResult := r.DB.Model(order).Where("id = ?", order.ID).Updates(order)
	if err := dbResult.Error; err != nil {
		return err
	}
	return nil
}

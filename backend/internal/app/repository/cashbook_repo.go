package repository

import (
	"context"
	"time"

	"github.com/fahmyabida/labbaika/internal/app/domain"
	pkgErrors "github.com/fahmyabida/labbaika/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type CashbookRepository struct {
	DB *gorm.DB
}

func NewCashbookRepository(db *gorm.DB) domain.ICashbookRepo {
	return &CashbookRepository{
		DB: db,
	}
}

func (r CashbookRepository) Create(ctx context.Context, data *domain.Cashbook) (err error) {
	dbResult := r.DB.WithContext(ctx).Create(data)
	if dbResult.Error != nil {
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		postgresError, ok := dbResult.Error.(*pgconn.PgError)
		if ok && postgresError.Code == "23505" {
			return pkgErrors.DuplicateCashbookError(pkgErrors.ErrDuplicateCashbook)
		}
		return dbResult.Error
	}

	return
}

func (r CashbookRepository) FindByID(ctx context.Context, id string) (entry domain.Cashbook, err error) {
	dbResult := r.DB.Model(&entry).Where("id = ?", id).Find(&entry)
	if dbResult.RowsAffected == 0 {
		return entry, pkgErrors.CashbookNotFoundError(pkgErrors.ErrCashbookNotFound)
	}
	return entry, dbResult.Error
}

// GetCashbooks lists entries in [startDate, endDate). endDate is expected to
// already be an exclusive upper bound (e.g. "day after the requested end
// date") so a whole calendar day is matched regardless of time-of-day —
// callers translating a human "YYYY-MM-DD" end date should bump it by one
// day before calling this.
func (r CashbookRepository) GetCashbooks(ctx context.Context, startDate, endDate time.Time) (entries []domain.Cashbook, err error) {
	query := r.DB.WithContext(ctx).Where("deleted_at IS NULL")
	if !startDate.IsZero() {
		query = query.Where("created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("created_at < ?", endDate)
	}

	// No rows matching the filter is a normal, valid list result — not an error.
	if err = query.Order("created_at DESC").Find(&entries).Error; err != nil {
		return entries, err
	}
	return entries, nil
}

// GetSummary mirrors GetCashbooks' date-window semantics (endDate exclusive).
func (r CashbookRepository) GetSummary(ctx context.Context, startDate, endDate time.Time) (summary []domain.CategorySummary, err error) {
	query := r.DB.WithContext(ctx).Model(&domain.Cashbook{}).
		Select("category, type, COALESCE(SUM(amount), 0) as total").
		Where("deleted_at IS NULL").
		Group("category, type")

	if !startDate.IsZero() {
		query = query.Where("created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("created_at < ?", endDate)
	}

	if err = query.Scan(&summary).Error; err != nil {
		return nil, err
	}
	return summary, nil
}

func (r CashbookRepository) Update(ctx context.Context, data *domain.Cashbook) error {
	dbResult := r.DB.Model(data).Where("id = ?", data.ID).Updates(data)
	if err := dbResult.Error; err != nil {
		return err
	}
	return nil
}

// Delete soft-deletes by setting deleted_at — DeletedAt is a plain *time.Time
// (not gorm.DeletedAt), so GORM's automatic soft-delete hook doesn't apply here.
func (r CashbookRepository) Delete(ctx context.Context, id string) error {
	now := time.Now()
	dbResult := r.DB.WithContext(ctx).Model(&domain.Cashbook{}).Where("id = ?", id).Update("deleted_at", &now)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	if dbResult.RowsAffected == 0 {
		return pkgErrors.CashbookNotFoundError(pkgErrors.ErrCashbookNotFound)
	}
	return nil
}

// UpdateReceipt sets just the receipt_path column — a single-column update so
// it can't accidentally null out other fields the way a zero-valued struct
// passed to GORM's Updates() could.
func (r CashbookRepository) UpdateReceipt(ctx context.Context, id string, receiptURL string) error {
	dbResult := r.DB.WithContext(ctx).Model(&domain.Cashbook{}).Where("id = ?", id).Update("receipt_path", receiptURL)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	if dbResult.RowsAffected == 0 {
		return pkgErrors.CashbookNotFoundError(pkgErrors.ErrCashbookNotFound)
	}
	return nil
}

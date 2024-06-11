package repository

import (
	"context"

	"github.com/fahmyabida/labbaika/internal/app/domain"
	pkgErrors "github.com/fahmyabida/labbaika/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type MenuRepository struct {
	DB *gorm.DB
}

func NewMenuRepository(db *gorm.DB) domain.IMenuRepo {
	return &MenuRepository{
		DB: db,
	}
}

func (r MenuRepository) Create(ctx context.Context, data *domain.Menu) (err error) {
	dbResult := r.DB.WithContext(ctx).Create(data)
	if dbResult.Error != nil {
		// https://www.postgresql.org/docs/current/errcodes-appendix.html
		postgresError, ok := dbResult.Error.(*pgconn.PgError)
		if ok && postgresError.Code == "23505" {
			return pkgErrors.DuplicateMenuError(pkgErrors.ErrDuplicateMenu)
		}
		return dbResult.Error
	}

	return
}

func (r MenuRepository) FindByID(ctx context.Context, id string) (order domain.Menu, err error) {
	dbResult := r.DB.Model(&order).Where("id = ?", id).Find(&order)
	if dbResult.RowsAffected == 0 {
		return order, pkgErrors.MenuNotFoundError(pkgErrors.ErrMenuNotFound)
	}
	return order, dbResult.Error
}

func (r MenuRepository) GetMenus(ctx context.Context, parameters interface{}) (orders []domain.Menu, err error) {
	dbResult := r.DB.Model(&orders).Where("param = ?", parameters).Find(&orders)
	if dbResult.RowsAffected == 0 {
		return orders, pkgErrors.MenuNotFoundError(pkgErrors.ErrMenuNotFound)
	} else if err = dbResult.Error; err != nil {
		return orders, err
	}
	return orders, nil
}

func (r MenuRepository) Update(ctx context.Context, order *domain.Menu) error {
	dbResult := r.DB.Model(order).Where("id = ?", order.ID).Updates(order)
	if err := dbResult.Error; err != nil {
		return err
	}
	return nil
}

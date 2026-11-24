package repository

import (
	"context"
	"errors"

	"github.com/fahmyabida/labbaika/internal/app/domain"
	pkgErrors "github.com/fahmyabida/labbaika/pkg/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.ICategoryRepo {
	return &CategoryRepository{DB: db}
}

func (r CategoryRepository) Create(ctx context.Context, data *domain.Category) error {
	dbResult := r.DB.WithContext(ctx).Create(data)
	if dbResult.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(dbResult.Error, &pgErr) && pgErr.Code == "23505" {
			return pkgErrors.DuplicateCategoryError(pkgErrors.ErrDuplicateCategory)
		}
		return dbResult.Error
	}
	return nil
}

func (r CategoryRepository) GetCategories(ctx context.Context) (categories []domain.Category, err error) {
	if err = r.DB.WithContext(ctx).Order("name").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r CategoryRepository) FindByID(ctx context.Context, id string) (category domain.Category, err error) {
	dbResult := r.DB.WithContext(ctx).Where("id = ?", id).Find(&category)
	if dbResult.Error != nil {
		return category, dbResult.Error
	}
	if dbResult.RowsAffected == 0 {
		return category, pkgErrors.CategoryNotFoundError(pkgErrors.ErrCategoryNotFound)
	}
	return category, nil
}

func (r CategoryRepository) FindByName(ctx context.Context, name string) (category domain.Category, err error) {
	dbResult := r.DB.WithContext(ctx).Where("LOWER(name) = LOWER(?)", name).Find(&category)
	if dbResult.Error != nil {
		return category, dbResult.Error
	}
	if dbResult.RowsAffected == 0 {
		return category, pkgErrors.CategoryNotFoundError(pkgErrors.ErrCategoryNotFound)
	}
	return category, nil
}

// Update renames the category and cascades the rename to any cashbook rows
// already using the old name. Category is stored on Cashbook by name, not a
// foreign key (see package doc) — without this cascade, a rename silently
// splits one category into two for any report/sum done by name (e.g. the
// dashboard), since past rows would keep the stale name forever.
func (r CategoryRepository) Update(ctx context.Context, data *domain.Category) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.Category
		if err := tx.Where("id = ?", data.ID).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return pkgErrors.CategoryNotFoundError(pkgErrors.ErrCategoryNotFound)
			}
			return err
		}

		dbResult := tx.Model(data).Where("id = ?", data.ID).Updates(data)
		if dbResult.Error != nil {
			var pgErr *pgconn.PgError
			if errors.As(dbResult.Error, &pgErr) && pgErr.Code == "23505" {
				return pkgErrors.DuplicateCategoryError(pkgErrors.ErrDuplicateCategory)
			}
			return dbResult.Error
		}
		if dbResult.RowsAffected == 0 {
			return pkgErrors.CategoryNotFoundError(pkgErrors.ErrCategoryNotFound)
		}

		if existing.Name != data.Name {
			if err := tx.Model(&domain.Cashbook{}).
				Where("LOWER(category) = LOWER(?)", existing.Name).
				Update("category", data.Name).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete is a hard delete — Category has no soft-delete column. Note this
// does not touch any Cashbook rows already using this category's name (the
// link is by name, not a foreign key), so past entries keep their text as-is.
func (r CategoryRepository) Delete(ctx context.Context, id string) error {
	dbResult := r.DB.WithContext(ctx).Where("id = ?", id).Delete(&domain.Category{})
	if dbResult.Error != nil {
		return dbResult.Error
	}
	if dbResult.RowsAffected == 0 {
		return pkgErrors.CategoryNotFoundError(pkgErrors.ErrCategoryNotFound)
	}
	return nil
}

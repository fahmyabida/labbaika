package domain

import (
	"context"
	"time"
)

// Category is master data for Cashbook.Category. Cashbook.Category must match
// an existing Category.Name (case-insensitive) — categories are managed
// explicitly through this CRUD, not created ad hoc from the cashbook form.
type Category struct {
	ID        string     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string     `json:"name" gorm:"type:text;not null"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
}

type ICategoryUsecase interface {
	Create(c context.Context, data *Category) (response interface{}, err error)
	Get(c context.Context) (response interface{}, err error)
	GetById(c context.Context, id string) (response interface{}, err error)
	Update(c context.Context, data *Category) (response interface{}, err error)
	Delete(c context.Context, id string) (response interface{}, err error)
}

type ICategoryRepo interface {
	Create(ctx context.Context, data *Category) error
	GetCategories(ctx context.Context) (categories []Category, err error)
	FindByID(ctx context.Context, id string) (Category, error)
	// FindByName looks up a category case-insensitively — used to validate
	// Cashbook.Category against the master list.
	FindByName(ctx context.Context, name string) (Category, error)
	Update(ctx context.Context, data *Category) error
	Delete(ctx context.Context, id string) error
}

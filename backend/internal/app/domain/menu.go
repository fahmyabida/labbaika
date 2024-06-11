package domain

import (
	"context"
	"time"
)

type Menu struct {
	ID          string     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string     `json:"name" gorm:"type:text"`
	Price       float64    `json:"price"`
	Description string     `json:"description" gorm:"type:text"`
	Status      string     `json:"status" gorm:"type:text"`
	CreatedAt   *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
	CreatedBy   string     `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
	UpdatedBy   string     `json:"updated_by"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
}

type IMenuUsecase interface {
	Create(c context.Context) (response interface{}, err error)
	Get(c context.Context) (response interface{}, err error)
	GetById(c context.Context) (response interface{}, err error)
	Update(c context.Context) (response interface{}, err error)
	Delete(c context.Context) (response interface{}, err error)
}

type IMenuRepo interface {
	Create(ctx context.Context, data *Menu) (err error)
	FindByID(ctx context.Context, id string) (order Menu, err error)
	GetMenus(ctx context.Context, parameters interface{}) (orders []Menu, err error)
	Update(ctx context.Context, order *Menu) error
}

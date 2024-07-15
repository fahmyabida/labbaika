package domain

import (
	"context"
	"time"
)

type LedgerType string

const (
	CREDIT LedgerType = "CREDIT"
	DEBIT  LedgerType = "DEBIT"
)

type Ledger struct {
	ID          string     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Type        LedgerType `json:"type" gorm:"type:text"`
	Amount      float64    `json:"amount" gorm:"type:numeric"`
	Category    string     `json:"category" gorm:"type:text"`
	Description string     `json:"description" gorm:"type:text"`
	CreatedAt   *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
	CreatedBy   string     `json:"created_by"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
	UpdatedBy   string     `json:"updated_by"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
}

type ILedgerUsecase interface {
	Create(c context.Context) (response interface{}, err error)
	Get(c context.Context) (response interface{}, err error)
	GetById(c context.Context) (response interface{}, err error)
	Update(c context.Context) (response interface{}, err error)
	Delete(c context.Context) (response interface{}, err error)
}

type ILedgerRepo interface {
	Create(ctx context.Context, data *Ledger) (err error)
	FindByID(ctx context.Context, id string) (data Ledger, err error)
	GetLedgers(ctx context.Context, parameters interface{}) (datas []Ledger, err error)
	Update(ctx context.Context, data *Ledger) error
}

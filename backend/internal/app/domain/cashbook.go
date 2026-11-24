package domain

import (
	"context"
	"time"
)

type CashbookType string

const (
	IN  CashbookType = "IN"
	OUT CashbookType = "OUT"
)

type Cashbook struct {
	ID          string       `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Type        CashbookType `json:"type" gorm:"type:text"`
	Amount      float64      `json:"amount" gorm:"type:numeric"`
	Category    string       `json:"category" gorm:"type:text"`
	Description string       `json:"description" gorm:"type:text"`
	CreatedAt   *time.Time   `json:"created_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
	CreatedBy   string       `json:"created_by"`
	UpdatedAt   *time.Time   `json:"updated_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
	UpdatedBy   string       `json:"updated_by"`
	DeletedAt   *time.Time   `json:"deleted_at,omitempty" gorm:"type:timestamptz;NOT NULL;default:null"`
	// ReceiptURL is a server-relative path (e.g. "/receipts/<file>") ready to
	// append to the frontend's backend base URL, or "" if nothing's attached.
	ReceiptURL string `json:"receipt_url,omitempty" gorm:"column:receipt_path;type:text"`
}

// CategorySummary is one (category, type) bucket's total — the shape the
// dashboard bar chart consumes.
type CategorySummary struct {
	Category string       `json:"category"`
	Type     CashbookType `json:"type"`
	Total    float64      `json:"total"`
}

type ICashbookUsecase interface {
	Create(c context.Context, data *Cashbook) (response interface{}, err error)
	Get(c context.Context, startDate, endDate time.Time) (response interface{}, err error)
	GetById(c context.Context, id string) (response interface{}, err error)
	Update(c context.Context, data *Cashbook) (response interface{}, err error)
	Delete(c context.Context, id string) (response interface{}, err error)
	// UploadReceipt attaches an already-saved file at receiptURL to the
	// cashbook entry. Saving/converting the file itself is an HTTP-layer
	// concern (handler), not this usecase's.
	UploadReceipt(c context.Context, id string, receiptURL string) (response interface{}, err error)
	GetSummary(c context.Context, startDate, endDate time.Time) (response interface{}, err error)
}

type ICashbookRepo interface {
	Create(ctx context.Context, data *Cashbook) (err error)
	FindByID(ctx context.Context, id string) (data Cashbook, err error)
	// GetCashbooks lists entries, optionally bounded by created_at. A zero
	// time.Time on either side means that side is unbounded. Excludes
	// soft-deleted rows.
	GetCashbooks(ctx context.Context, startDate, endDate time.Time) (datas []Cashbook, err error)
	Update(ctx context.Context, data *Cashbook) error
	Delete(ctx context.Context, id string) error
	UpdateReceipt(ctx context.Context, id string, receiptURL string) error
	// GetSummary sums amount grouped by (category, type) over the same
	// [startDate, endDate) window semantics as GetCashbooks, excluding
	// soft-deleted rows.
	GetSummary(ctx context.Context, startDate, endDate time.Time) ([]CategorySummary, error)
}

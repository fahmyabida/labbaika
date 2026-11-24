package usecase

import (
	"context"
	"time"

	"github.com/fahmyabida/labbaika/internal/app/domain"
	"github.com/fahmyabida/labbaika/internal/logger"
	pkgErrors "github.com/fahmyabida/labbaika/pkg/errors"
)

type CashbookUsecase struct {
	CashbookRepository domain.ICashbookRepo
	CategoryRepository domain.ICategoryRepo
}

func NewCashbookUsecase(cashbookRepository domain.ICashbookRepo, categoryRepository domain.ICategoryRepo) *CashbookUsecase {
	return &CashbookUsecase{cashbookRepository, categoryRepository}
}

// validate enforces required fields before anything touches the DB: amount
// must be a real positive value and category must be set (categories are
// mandatory, not optional, now that there's a managed master list to pick
// from).
func validateCashbook(data *domain.Cashbook) error {
	if data.Amount <= 0 {
		return pkgErrors.InvalidInputError("amount must be greater than zero")
	}
	if data.Category == "" {
		return pkgErrors.InvalidInputError("category is required")
	}
	return nil
}

// resolveCategory validates data.Category against the master Category list
// (case-insensitive) and normalizes it to the stored casing. Categories are
// managed explicitly via the Category CRUD — an unknown name is rejected
// here, not auto-created. A blank category is left as-is — it's optional.
func (u *CashbookUsecase) resolveCategory(c context.Context, data *domain.Cashbook) error {
	if data.Category == "" {
		return nil
	}
	category, err := u.CategoryRepository.FindByName(c, data.Category)
	if err != nil {
		logger.Error(c, "category not in master list", err, logger.Fields{"category": data.Category})
		return err
	}
	data.Category = category.Name
	return nil
}

func (u *CashbookUsecase) Create(c context.Context, data *domain.Cashbook) (response interface{}, err error) {
	if err = validateCashbook(data); err != nil {
		return nil, err
	}
	if err = u.resolveCategory(c, data); err != nil {
		return nil, err
	}
	if err = u.CashbookRepository.Create(c, data); err != nil {
		logger.Error(c, "failed to create cashbook", err)
		return nil, err
	}
	return data, nil
}

func (u *CashbookUsecase) Get(c context.Context, startDate, endDate time.Time) (response interface{}, err error) {
	response, err = u.CashbookRepository.GetCashbooks(c, startDate, endDate)
	if err != nil {
		logger.Error(c, "failed to get cashbooks", err, logger.Fields{"start_date": startDate, "end_date": endDate})
	}
	return response, err
}

func (u *CashbookUsecase) GetById(c context.Context, id string) (response interface{}, err error) {
	response, err = u.CashbookRepository.FindByID(c, id)
	if err != nil {
		logger.Error(c, "failed to get cashbook by id", err, logger.Fields{"id": id})
	}
	return response, err
}

func (u *CashbookUsecase) Update(c context.Context, data *domain.Cashbook) (response interface{}, err error) {
	if err = validateCashbook(data); err != nil {
		return nil, err
	}
	if err = u.resolveCategory(c, data); err != nil {
		return nil, err
	}
	if err = u.CashbookRepository.Update(c, data); err != nil {
		logger.Error(c, "failed to update cashbook", err, logger.Fields{"id": data.ID})
		return nil, err
	}
	return data, nil
}

func (u *CashbookUsecase) Delete(c context.Context, id string) (response interface{}, err error) {
	if err = u.CashbookRepository.Delete(c, id); err != nil {
		logger.Error(c, "failed to delete cashbook", err, logger.Fields{"id": id})
		return nil, err
	}
	return nil, nil
}

func (u *CashbookUsecase) UploadReceipt(c context.Context, id string, receiptURL string) (response interface{}, err error) {
	if err = u.CashbookRepository.UpdateReceipt(c, id, receiptURL); err != nil {
		logger.Error(c, "failed to attach receipt", err, logger.Fields{"id": id})
		return nil, err
	}
	return u.GetById(c, id)
}

func (u *CashbookUsecase) GetSummary(c context.Context, startDate, endDate time.Time) (response interface{}, err error) {
	response, err = u.CashbookRepository.GetSummary(c, startDate, endDate)
	if err != nil {
		logger.Error(c, "failed to get cashbook summary", err, logger.Fields{"start_date": startDate, "end_date": endDate})
	}
	return response, err
}

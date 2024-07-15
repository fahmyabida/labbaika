package usecase

import (
	"context"

	"github.com/fahmyabida/labbaika/internal/app/domain"
)

type LedgerUsecase struct {
	LedgerRepository domain.ILedgerRepo
}

func NewLedgerUsecase(LedgerRepository domain.ILedgerRepo) *LedgerUsecase {
	return &LedgerUsecase{LedgerRepository}
}

func (u *LedgerUsecase) Create(c context.Context) (response interface{}, err error) {

	return nil, nil
}

func (u *LedgerUsecase) Get(c context.Context) (response interface{}, err error) {
	return nil, nil
}

func (u *LedgerUsecase) GetById(c context.Context) (response interface{}, err error) {
	return nil, nil
}

func (u *LedgerUsecase) Update(c context.Context) (response interface{}, err error) {
	return nil, nil
}

func (u *LedgerUsecase) Delete(c context.Context) (response interface{}, err error) {
	return nil, nil
}

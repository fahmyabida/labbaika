package usecase

import (
	"context"

	"github.com/fahmyabida/labbaika/internal/app/domain"
)

type MenuUsecase struct {
	MenuRepository domain.IMenuRepo
}

func NewMenuUsecase(MenuRepository domain.IMenuRepo) *MenuUsecase {
	return &MenuUsecase{MenuRepository}
}

func (u *MenuUsecase) Create(c context.Context) (response interface{}, err error) {

	return nil, nil
}

func (u *MenuUsecase) Get(c context.Context) (response interface{}, err error) {
	return nil, nil
}

func (u *MenuUsecase) GetById(c context.Context) (response interface{}, err error) {
	return nil, nil
}

func (u *MenuUsecase) Update(c context.Context) (response interface{}, err error) {
	return nil, nil
}

func (u *MenuUsecase) Delete(c context.Context) (response interface{}, err error) {
	return nil, nil
}

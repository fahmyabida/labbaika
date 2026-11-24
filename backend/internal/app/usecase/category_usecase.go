package usecase

import (
	"context"

	"github.com/fahmyabida/labbaika/internal/app/domain"
	"github.com/fahmyabida/labbaika/internal/logger"
)

type CategoryUsecase struct {
	CategoryRepository domain.ICategoryRepo
}

func NewCategoryUsecase(categoryRepository domain.ICategoryRepo) *CategoryUsecase {
	return &CategoryUsecase{categoryRepository}
}

func (u *CategoryUsecase) Create(c context.Context, data *domain.Category) (response interface{}, err error) {
	if err = u.CategoryRepository.Create(c, data); err != nil {
		logger.Error(c, "failed to create category", err, logger.Fields{"name": data.Name})
		return nil, err
	}
	return data, nil
}

func (u *CategoryUsecase) Get(c context.Context) (response interface{}, err error) {
	response, err = u.CategoryRepository.GetCategories(c)
	if err != nil {
		logger.Error(c, "failed to list categories", err)
	}
	return response, err
}

func (u *CategoryUsecase) GetById(c context.Context, id string) (response interface{}, err error) {
	response, err = u.CategoryRepository.FindByID(c, id)
	if err != nil {
		logger.Error(c, "failed to get category by id", err, logger.Fields{"id": id})
	}
	return response, err
}

func (u *CategoryUsecase) Update(c context.Context, data *domain.Category) (response interface{}, err error) {
	if err = u.CategoryRepository.Update(c, data); err != nil {
		logger.Error(c, "failed to update category", err, logger.Fields{"id": data.ID})
		return nil, err
	}
	return data, nil
}

func (u *CategoryUsecase) Delete(c context.Context, id string) (response interface{}, err error) {
	if err = u.CategoryRepository.Delete(c, id); err != nil {
		logger.Error(c, "failed to delete category", err, logger.Fields{"id": id})
		return nil, err
	}
	return nil, nil
}

package handler

import (
	"net/http"

	"github.com/fahmyabida/labbaika/internal/app/domain"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	CategoryUsecase domain.ICategoryUsecase
}

func InitCategoryHandler(e *echo.Group, categoryUsecase domain.ICategoryUsecase) {
	handler := CategoryHandler{CategoryUsecase: categoryUsecase}

	e.POST("/category", handler.CreateCategory)
	e.GET("/category", handler.GetCategory)
	e.GET("/category/:id", handler.GetCategoryByID)
	e.PUT("/category/:id", handler.UpdateCategory)
	e.DELETE("/category/:id", handler.DeleteCategory)
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	var body domain.Category
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	body.ID = ""

	response, err := h.CategoryUsecase.Create(c.Request().Context(), &body)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, response)
}

func (h *CategoryHandler) GetCategory(c echo.Context) error {
	response, err := h.CategoryUsecase.Get(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CategoryHandler) GetCategoryByID(c echo.Context) error {
	response, err := h.CategoryUsecase.GetById(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	var body domain.Category
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	body.ID = c.Param("id")

	response, err := h.CategoryUsecase.Update(c.Request().Context(), &body)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	if _, err := h.CategoryUsecase.Delete(c.Request().Context(), c.Param("id")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

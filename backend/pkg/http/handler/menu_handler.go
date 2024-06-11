package handler

import (
	"net/http"

	"github.com/fahmyabida/labbaika/internal/app/domain"
	"github.com/labstack/echo/v4"
)

type MenuHandler struct {
	MenuUsecase domain.IMenuUsecase
}

func InitMenuHandler(e *echo.Group, MenuUsecase domain.IMenuUsecase) {
	handler := MenuHandler{MenuUsecase: MenuUsecase}

	e.POST("/menu", handler.CreateMenu)
	e.GET("/menu", handler.GetMenu)
	e.GET("/menu/{:id}", handler.GetMenuByID)
	e.PUT("/menu", handler.UpdateMenu)
	e.DELETE("/menu", handler.DeleteMenu)
}

func (h *MenuHandler) CreateMenu(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "implement me!",
	})
}

func (h *MenuHandler) GetMenu(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "implement me!",
	})
}

func (h *MenuHandler) GetMenuByID(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "implement me!",
	})
}

func (h *MenuHandler) UpdateMenu(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "implement me!",
	})
}

func (h *MenuHandler) DeleteMenu(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "implement me!",
	})
}

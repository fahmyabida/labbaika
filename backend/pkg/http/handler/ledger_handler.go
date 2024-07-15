package handler

import (
	"net/http"

	"github.com/fahmyabida/labbaika/internal/app/domain"
	"github.com/labstack/echo/v4"
)

type LedgerHandler struct {
	LedgerUsecase domain.ILedgerUsecase
}

func InitLedgerHandler(e *echo.Group, LedgerUsecase domain.ILedgerUsecase) {
	handler := LedgerHandler{LedgerUsecase: LedgerUsecase}

	e.POST("/menu", handler.CreateLedger)
	e.GET("/menu", handler.GetLedger)
	e.GET("/menu/{:id}", handler.GetLedgerByID)
	e.PUT("/menu", handler.UpdateLedger)
	e.DELETE("/menu", handler.DeleteLedger)
}

func (h *LedgerHandler) CreateLedger(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "implement me!",
	})
}

func (h *LedgerHandler) GetLedger(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "implement me!",
	})
}

func (h *LedgerHandler) GetLedgerByID(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "implement me!",
	})
}

func (h *LedgerHandler) UpdateLedger(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "implement me!",
	})
}

func (h *LedgerHandler) DeleteLedger(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "implement me!",
	})
}

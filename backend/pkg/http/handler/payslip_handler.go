package handler

import (
	"github.com/fahmyabida/labbaika/internal/app/domain"
	pkgErrors "github.com/fahmyabida/labbaika/pkg/errors"
	"github.com/labstack/echo/v4"
)

type PayslipHandler struct {
	PayslipUsecase domain.IPayslipUsecase
}

func InitPayslipHandler(e *echo.Group, PayslipUsecase domain.IPayslipUsecase) {
	handler := PayslipHandler{PayslipUsecase: PayslipUsecase}

	e.POST("/payslips/convert", handler.ConvertCsvToDocx)
}

func (h *PayslipHandler) ConvertCsvToDocx(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return pkgErrors.FileNotFoundError(pkgErrors.ErrFileNotFound)
	}

	src, err := file.Open()
	if err != nil {
		return pkgErrors.FileNotFoundError(pkgErrors.ErrFileNotFound)
	}
	defer src.Close()

	ctx := c.Request().Context()
	fileName, err := h.PayslipUsecase.ConvertPayslip(ctx, src)
	if err != nil {
		return err
	}
	return c.Attachment(fileName, fileName)
}

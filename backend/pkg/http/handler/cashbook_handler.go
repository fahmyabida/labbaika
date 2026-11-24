package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fahmyabida/labbaika/internal/app/domain"
	"github.com/labstack/echo/v4"
)

type CashbookHandler struct {
	CashbookUsecase domain.ICashbookUsecase
}

// receiptDir is where uploaded receipt files land, served back at /receipts/*
// (registered alongside the echo instance in cashbook-service.go).
const receiptDir = "uploads/receipts"

// maxReceiptSize caps the uploaded file before any webp conversion — the
// whole point of converting is to not let the server's disk fill up with
// unprocessed originals, so an unbounded upload would defeat that.
const maxReceiptSize = 5 << 20 // 5MB

func InitCashbookHandler(e *echo.Group, cashbookUsecase domain.ICashbookUsecase) {
	handler := CashbookHandler{CashbookUsecase: cashbookUsecase}

	e.POST("/cashbook", handler.CreateCashbook)
	e.GET("/cashbook", handler.GetCashbook)
	e.GET("/cashbook/summary", handler.GetCashbookSummary)
	e.GET("/cashbook/:id", handler.GetCashbookByID)
	e.PUT("/cashbook/:id", handler.UpdateCashbook)
	e.DELETE("/cashbook/:id", handler.DeleteCashbook)
	e.POST("/cashbook/:id/receipt", handler.UploadReceipt)
}

func (h *CashbookHandler) CreateCashbook(c echo.Context) error {
	var body domain.Cashbook
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	body.ID = "" // DB generates this

	response, err := h.CashbookUsecase.Create(c.Request().Context(), &body)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, response)
}

func (h *CashbookHandler) GetCashbook(c echo.Context) error {
	startDate, endDate, err := parseDateRangeParams(c)
	if err != nil {
		return err
	}

	response, err := h.CashbookUsecase.Get(c.Request().Context(), startDate, endDate)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CashbookHandler) GetCashbookSummary(c echo.Context) error {
	startDate, endDate, err := parseDateRangeParams(c)
	if err != nil {
		return err
	}

	response, err := h.CashbookUsecase.GetSummary(c.Request().Context(), startDate, endDate)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CashbookHandler) GetCashbookByID(c echo.Context) error {
	response, err := h.CashbookUsecase.GetById(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CashbookHandler) UpdateCashbook(c echo.Context) error {
	var body domain.Cashbook
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	body.ID = c.Param("id")

	response, err := h.CashbookUsecase.Update(c.Request().Context(), &body)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response)
}

func (h *CashbookHandler) DeleteCashbook(c echo.Context) error {
	if _, err := h.CashbookUsecase.Delete(c.Request().Context(), c.Param("id")); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *CashbookHandler) UploadReceipt(c echo.Context) error {
	id := c.Param("id")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing file field")
	}
	if fileHeader.Size > maxReceiptSize {
		return echo.NewHTTPError(http.StatusRequestEntityTooLarge, "file too large, max 5MB")
	}

	receiptURL, err := saveReceiptFile(fileHeader, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save file")
	}

	response, err := h.CashbookUsecase.UploadReceipt(c.Request().Context(), id, receiptURL)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response)
}

// saveReceiptFile writes the uploaded file to disk, best-effort converts it
// to webp, and returns the server-relative URL to store on the cashbook row.
// Conversion is skipped (not failed) when the `cwebp` binary isn't installed
// — the original file is kept as-is in that case.
func saveReceiptFile(fileHeader *multipart.FileHeader, cashbookID string) (string, error) {
	if err := os.MkdirAll(receiptDir, 0o755); err != nil {
		return "", err
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	rawName := fmt.Sprintf("%s_%d%s", cashbookID, time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	rawPath := filepath.Join(receiptDir, rawName)

	dst, err := os.Create(rawPath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return "", err
	}
	dst.Close()

	finalPath := rawPath
	if webpPath, converted := convertToWebP(rawPath); converted {
		os.Remove(rawPath) // keep only the smaller webp copy
		finalPath = webpPath
	}

	return "/receipts/" + filepath.Base(finalPath), nil
}

// convertToWebP shells out to the `cwebp` CLI (no cgo/libwebp binding needed
// in the Go binary itself). Returns ("", false) if cwebp isn't on PATH or the
// conversion fails — caller falls back to keeping the original file.
func convertToWebP(srcPath string) (string, bool) {
	if _, err := exec.LookPath("cwebp"); err != nil {
		return "", false
	}
	dstPath := strings.TrimSuffix(srcPath, filepath.Ext(srcPath)) + ".webp"
	if err := exec.Command("cwebp", "-quiet", "-q", "80", srcPath, "-o", dstPath).Run(); err != nil {
		return "", false
	}
	return dstPath, true
}

// parseDateParam parses "YYYY-MM-DD". Empty string means unbounded (zero time.Time).
// Parsed in time.Local, not UTC — created_at is stored/read as a local-offset
// timestamptz, so a bare UTC parse would misalign "which calendar day" a
// record near midnight actually falls on for whoever's reading these dates.
func parseDateParam(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.ParseInLocation("2006-01-02", value, time.Local)
}

// parseDateRangeParams reads start_date/end_date query params and bumps
// end_date to the start of the next day — it's meant inclusively ("show me
// all of Oct 9") but repo filters use created_at < endDate.
func parseDateRangeParams(c echo.Context) (startDate, endDate time.Time, err error) {
	startDate, err = parseDateParam(c.QueryParam("start_date"))
	if err != nil {
		return startDate, endDate, echo.NewHTTPError(http.StatusBadRequest, "invalid start_date, expected YYYY-MM-DD")
	}
	endDate, err = parseDateParam(c.QueryParam("end_date"))
	if err != nil {
		return startDate, endDate, echo.NewHTTPError(http.StatusBadRequest, "invalid end_date, expected YYYY-MM-DD")
	}
	if !endDate.IsZero() {
		endDate = endDate.AddDate(0, 0, 1)
	}
	return startDate, endDate, nil
}

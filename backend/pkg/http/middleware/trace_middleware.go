package middleware

import (
	"context"
	"time"

	"github.com/fahmyabida/labbaika-payslip/internal/logger"
	"github.com/labstack/echo/v4"
)

func TraceMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			traceHandler(c)
			return next(c)
		}
	}
}

func traceHandler(c echo.Context) {
	traceID := c.Request().Header.Get(string(logger.XTraceID))
	if traceID == "" {
		traceID = "lbk-" + time.Now().Format("20060102150405999999999")
		c.Request().Header.Add(string(logger.XTraceID), traceID)
	}
	c.Response().Header().Add(string(logger.XTraceID), traceID)
	ctx := context.WithValue(c.Request().Context(), logger.XTraceID, traceID)
	c.SetRequest(c.Request().WithContext(ctx))
}

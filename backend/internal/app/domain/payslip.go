package domain

import (
	"context"
	"io"
)

type IPayslipUsecase interface {
	ConvertPayslip(c context.Context, file io.Reader) (fileName string, err error)
}

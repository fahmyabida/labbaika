package app

import (
	"fmt"
	"os"

	"github.com/fahmyabida/labbaika-payslip/internal/app/domain"
	"github.com/fahmyabida/labbaika-payslip/internal/app/usecase"
	"github.com/fahmyabida/labbaika-payslip/internal/logger"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "labbaika-service",
		Short: "labbaika-service is application to manage Labbaika service",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	PayslipUsecase domain.IPayslipUsecase
)

func init() {
	initLogger()
	initApp()
}

func initApp() {
	PayslipUsecase = usecase.NewPayslipUsecase()
}

func initLogger() {
	logger.InitStaticLogger("labbaika-svc")
}

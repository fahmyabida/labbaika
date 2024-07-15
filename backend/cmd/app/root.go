package app

import (
	"fmt"
	"log"
	"os"

	"github.com/fahmyabida/labbaika/cmd/config"
	"github.com/fahmyabida/labbaika/internal/app/domain"
	"github.com/fahmyabida/labbaika/internal/app/repository"
	"github.com/fahmyabida/labbaika/internal/app/usecase"
	"github.com/fahmyabida/labbaika/internal/logger"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
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
	database *gorm.DB

	// repository
	LedgerRepository domain.ILedgerRepo
	MenuRepository   domain.IMenuRepo

	// usecase
	LedgerUsecase  domain.ILedgerUsecase
	MenuUsecase    domain.IMenuUsecase
	PayslipUsecase domain.IPayslipUsecase
)

func init() {
	initLogger()
	if err := config.InitEnv(); err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(func() {
		// initDatabase()
		initApp()
	})
}

func initApp() {
	MenuRepository = repository.NewMenuRepository(database)

	LedgerUsecase = usecase.NewLedgerUsecase(LedgerRepository)
	MenuUsecase = usecase.NewMenuUsecase(MenuRepository)
	PayslipUsecase = usecase.NewPayslipUsecase()
}

func initDatabase() {
	rw, ro := config.LoadForPostgres()
	database = config.InitDB(rw, ro)
}

func initLogger() {
	logger.InitStaticLogger("labbaika-svc")
}

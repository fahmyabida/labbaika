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
	CashbookRepository domain.ICashbookRepo
	CategoryRepository domain.ICategoryRepo

	// usecase
	CashbookUsecase domain.ICashbookUsecase
	CategoryUsecase domain.ICategoryUsecase
	PayslipUsecase  domain.IPayslipUsecase
)

func init() {
	initLogger()
	if err := config.InitEnv(); err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(initApp)
}

// noDBCommands are services that must boot without Postgres.
var noDBCommands = map[string]bool{"payslip-svc": true}

func initApp() {
	PayslipUsecase = usecase.NewPayslipUsecase()
	if len(os.Args) > 1 && noDBCommands[os.Args[1]] {
		return
	}

	initDatabase()
	CashbookRepository = repository.NewCashbookRepository(database)
	CategoryRepository = repository.NewCategoryRepository(database)

	CashbookUsecase = usecase.NewCashbookUsecase(CashbookRepository, CategoryRepository)
	CategoryUsecase = usecase.NewCategoryUsecase(CategoryRepository)
}

func initDatabase() {
	rw, ro := config.LoadForPostgres()
	database = config.InitDB(rw, ro)
}

func initLogger() {
	logger.InitStaticLogger("labbaika-svc")
}

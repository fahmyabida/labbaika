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

	MenuRepository domain.IMenuRepo

	PayslipUsecase domain.IPayslipUsecase
	MenuUsecase    domain.IMenuUsecase
)

func init() {
	initLogger()
	if err := config.InitEnv(); err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(func() {
		initDatabase()
		initApp()
	})
}

func initApp() {
	MenuRepository = repository.NewMenuRepository(database)

	PayslipUsecase = usecase.NewPayslipUsecase()
	MenuUsecase = usecase.NewMenuUsecase(MenuRepository)
}

func initDatabase() {
	rw, ro := config.LoadForPostgres()
	database = config.InitDB(rw, ro)
}

func initLogger() {
	logger.InitStaticLogger("labbaika-svc")
}

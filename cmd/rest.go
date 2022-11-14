package cmd

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var restServer = &cobra.Command{
	Use:   "rest",
	Short: "This is rest server for labbaika-payslip",
	Run:   RunRestServer,
}

func RunRestServer(cmd *cobra.Command, args []string) {
	e := echo.New()
	e.POST("/payslips/convert", payslip.ConvertPayslipCsvToDocxfunc)
	e.Logger.Fatal(e.Start(":8686"))
}

func Execute() {
	if err := restServer.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"

	httpHandler "github.com/fahmyabida/labbaika/pkg/http/handler"
	customMiddleware "github.com/fahmyabida/labbaika/pkg/http/middleware"
)

var menuSvc = &cobra.Command{
	Use:   "menu-svc",
	Short: "This is rest server for labbaika-menu",
	Run:   RunMenuService,
}

func init() {
	rootCmd.AddCommand(menuSvc)
}

func RunMenuService(cmd *cobra.Command, args []string) {
	e := echo.New()

	e.Use(
		customMiddleware.TraceMiddleware(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Skipper: middleware.DefaultSkipper,
			Format: `{"time":"${time_rfc3339_nano}","trace_id":"${header:x-trace-id}","remote_ip":"${remote_ip}",` +
				`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
				`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
				`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
			CustomTimeFormat: "2006-01-02 15:04:05.00000",
		}),
		customMiddleware.ErrorMiddleware(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:  []string{"*"},
			AllowHeaders:  []string{"*"},
			ExposeHeaders: []string{"Content-Disposition"},
		}),
	)

	v1 := e.Group("/api/v1")

	httpHandler.InitMenuHandler(v1, MenuUsecase)

	e.Logger.Fatal(e.Start(":8080"))
}

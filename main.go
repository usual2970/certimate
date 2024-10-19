package main

import (
	"log"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"certimate/internal/domains"
	"certimate/internal/utils/app"
	_ "certimate/migrations"
	"certimate/ui"

	_ "time/tzdata"
)

func main() {
	app := app.GetApp()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	domains.AddEvent()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		domains.InitSchedule()

		e.Router.GET(
			"/*",
			echo.StaticDirectoryHandler(ui.DistDirFS, false),
			middleware.Gzip(),
		)

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

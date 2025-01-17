package main

import (
	"flag"
	"log"
	"os"
	"strings"

	_ "time/tzdata"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/routes"
	"github.com/usual2970/certimate/internal/scheduler"
	"github.com/usual2970/certimate/internal/workflow"
	"github.com/usual2970/certimate/ui"

	_ "github.com/usual2970/certimate/migrations"
)

func main() {
	app := app.GetApp()

	var flagHttp string
	var flagDir string
	flag.StringVar(&flagHttp, "http", "127.0.0.1:8090", "HTTP server address")
	flag.StringVar(&flagDir, "dir", "/pb_data/database", "Pocketbase data directory")
	_ = flag.CommandLine.Parse(os.Args[2:]) // skip the first two arguments: "main.go serve"

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: strings.HasPrefix(os.Args[0], os.TempDir()),
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		scheduler.Register()

		workflow.Register()

		routes.Register(e.Router)
		e.Router.GET(
			"/*",
			echo.StaticDirectoryHandler(ui.DistDirFS, false),
			middleware.Gzip(),
		)

		return nil
	})

	app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		routes.Unregister()

		log.Println("Exit!")

		return nil
	})

	log.Printf("Visit the website: http://%s", flagHttp)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "github.com/usual2970/certimate/migrations"

	"github.com/usual2970/certimate/internal/domains"
	"github.com/usual2970/certimate/internal/routes"
	"github.com/usual2970/certimate/internal/utils/app"
	"github.com/usual2970/certimate/ui"

	_ "time/tzdata"
)

func main() {
	app := app.GetApp()

	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	// 获取启动命令中的http参数
	var httpFlag string
	flag.StringVar(&httpFlag, "http", "127.0.0.1:8090", "HTTP server address")
	// "serve"影响解析
	_ = flag.CommandLine.Parse(os.Args[2:])

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	domains.AddEvent()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		domains.InitSchedule()

		routes.Register(e.Router)

		e.Router.GET(
			"/*",
			echo.StaticDirectoryHandler(ui.DistDirFS, false),
			middleware.Gzip(),
		)

		return nil
	})

	defer log.Println("Exit!")
	log.Printf("Visit the website: http://%s", httpFlag)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

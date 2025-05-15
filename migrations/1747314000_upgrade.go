package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// migrate data
		{
			accesses, err := app.FindAllRecords("access")
			if err != nil {
				return err
			}

			for _, access := range accesses {
				changed := false

				if access.GetString("provider") == "goedge" {
					config := make(map[string]any)
					if err := access.UnmarshalJSONField("config", &config); err != nil {
						return err
					}

					config["apiRole"] = "user"
					access.Set("config", config)
					changed = true
				}

				if changed {
					err = app.Save(access)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	}, func(app core.App) error {
		return nil
	})
}

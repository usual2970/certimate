package migrations

import (
	"slices"

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

			providersToUpdate := []string{
				"1panel",
				"baotapanel",
				"baotawaf",
				"cdnfly",
				"flexcdn",
				"goedge",
				"lecdn",
				"powerdns",
				"proxmoxve",
				"ratpanel",
				"safeline",
			}
			for _, access := range accesses {
				changed := false

				if slices.Contains(providersToUpdate, access.GetString("provider")) {
					config := make(map[string]any)
					if err := access.UnmarshalJSONField("config", &config); err != nil {
						return err
					}

					config["serverUrl"] = config["apiUrl"]
					delete(config, "apiUrl")
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

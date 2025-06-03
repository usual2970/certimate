package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		tracer := NewTracer("(v0.3)1748959200")
		tracer.Printf("go ...")

		// migrate data
		{
			collection, err := app.FindCollectionByNameOrId("4yzbv8urny5ja1e")
			if err != nil {
				return err
			}

			records, err := app.FindAllRecords(collection)
			if err != nil {
				return err
			}

			for _, record := range records {
				changed := false

				if record.GetString("provider") == "ssh" {
					config := make(map[string]any)
					if err := record.UnmarshalJSONField("config", &config); err != nil {
						return err
					}

					if config["authMethod"] == nil || config["authMethod"] == "" {
						if config["key"] != nil && config["key"] != "" {
							config["authMethod"] = "key"
						} else if config["password"] != nil && config["password"] != "" {
							config["authMethod"] = "password"
						} else {
							config["authMethod"] = "none"
						}
						record.Set("config", config)
						changed = true
					}
				}

				if changed {
					if err := app.Save(record); err != nil {
						return err
					}

					tracer.Printf("record #%s in collection '%s' updated", record.Id, collection.Name)
				}
			}
		}

		tracer.Printf("done")
		return nil
	}, func(app core.App) error {
		return nil
	})
}

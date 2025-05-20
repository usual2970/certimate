package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// update collection `certificate`
		{
			collection, err := app.FindCollectionByNameOrId("4szxr9x43tpj6np")
			if err != nil {
				return err
			}

			if err := collection.Fields.AddMarshaledJSONAt(6, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "text2910474005",
				"max": 0,
				"min": 0,
				"name": "issuerOrg",
				"pattern": "",
				"presentable": false,
				"primaryKey": false,
				"required": false,
				"system": false,
				"type": "text"
			}`)); err != nil {
				return err
			}

			if err := app.Save(collection); err != nil {
				return err
			}
		}

		// migrate data
		{
			accesses, err := app.FindAllRecords("access")
			if err != nil {
				return err
			}

			for _, access := range accesses {
				changed := false

				if access.GetString("provider") == "1panel" {
					config := make(map[string]any)
					if err := access.UnmarshalJSONField("config", &config); err != nil {
						return err
					}

					config["apiVersion"] = "v1"
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

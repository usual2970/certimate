package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// update collection `access`
		{
			collection, err := app.FindCollectionByNameOrId("4yzbv8urny5ja1e")
			if err != nil {
				return err
			}

			if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "text2859962647",
				"max": 0,
				"min": 0,
				"name": "reserve",
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

				if access.GetString("provider") == "buypass" {
					access.Set("reserve", "ca")
					changed = true
				} else if access.GetString("provider") == "googletrustservices" {
					access.Set("reserve", "ca")
					changed = true
				} else if access.GetString("provider") == "sslcom" {
					access.Set("reserve", "ca")
					changed = true
				} else if access.GetString("provider") == "zerossl" {
					access.Set("reserve", "ca")
					changed = true
				}

				if access.GetString("provider") == "webhook" {
					config := make(map[string]any)
					if err := access.UnmarshalJSONField("config", &config); err != nil {
						return err
					}

					config["method"] = "POST"
					config["headers"] = "Content-Type: application/json"
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

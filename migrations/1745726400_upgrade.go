package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		tracer := NewTracer("(v0.3)1745726400")
		tracer.Printf("go ...")

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

			tracer.Printf("collection '%s' updated", collection.Name)
		}

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

				if record.GetString("provider") == "buypass" {
					record.Set("reserve", "ca")
					changed = true
				} else if record.GetString("provider") == "googletrustservices" {
					record.Set("reserve", "ca")
					changed = true
				} else if record.GetString("provider") == "sslcom" {
					record.Set("reserve", "ca")
					changed = true
				} else if record.GetString("provider") == "zerossl" {
					record.Set("reserve", "ca")
					changed = true
				}

				if record.GetString("provider") == "webhook" {
					config := make(map[string]any)
					if err := record.UnmarshalJSONField("config", &config); err != nil {
						return err
					}

					config["method"] = "POST"
					config["headers"] = "Content-Type: application/json"
					record.Set("config", config)
					changed = true
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

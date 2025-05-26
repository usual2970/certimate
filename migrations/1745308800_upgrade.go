package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		tracer := NewTracer("(v0.3)1745308800")
		tracer.Printf("go ...")

		// update collection `access`
		{
			collection, err := app.FindCollectionByNameOrId("4yzbv8urny5ja1e")
			if err != nil {
				return err
			}

			// add temp field `providerTmp`
			if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "text2024822322",
				"max": 0,
				"min": 0,
				"name": "providerTmp",
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

			// copy `provider` to `providerTmp`
			if _, err := app.DB().NewQuery("UPDATE access SET providerTmp = provider").Execute(); err != nil {
				return err
			}

			// remove old field `provider`
			collection.Fields.RemoveById("hwy7m03o")
			if err := json.Unmarshal([]byte(`{
				"indexes": [
					"CREATE INDEX `+"`"+`idx_wkoST0j`+"`"+` ON `+"`"+`access`+"`"+` (`+"`"+`name`+"`"+`)"
				]
			}`), &collection); err != nil {
				return err
			}
			if err := app.Save(collection); err != nil {
				return err
			}

			// rename field `providerTmp` to `provider`
			if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "text2024822322",
				"max": 0,
				"min": 0,
				"name": "provider",
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

			// rebuild indexes
			if err := json.Unmarshal([]byte(`{
				"indexes": [
					"CREATE INDEX `+"`"+`idx_wkoST0j`+"`"+` ON `+"`"+`access`+"`"+` (`+"`"+`name`+"`"+`)",
					"CREATE INDEX `+"`"+`idx_frh0JT1Aqx`+"`"+` ON `+"`"+`access`+"`"+` (`+"`"+`provider`+"`"+`)"
				]
			}`), &collection); err != nil {
				return err
			}
			if err := app.Save(collection); err != nil {
				return err
			}

			tracer.Printf("collection '%s' updated", collection.Name)
		}

		tracer.Printf("done")
		return nil
	}, func(app core.App) error {
		return nil
	})
}

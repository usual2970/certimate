package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// create collection `workflow_run`
		{
			collection, err := app.FindCollectionByNameOrId("qjp8lygssgwyqyz")
			if err != nil {
				return err
			}

			// update field
			if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "hvebkuxw",
				"max": 20000,
				"min": 0,
				"name": "error",
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

		// create collection `workflow_output`
		{
			collection, err := app.FindCollectionByNameOrId("bqnxb95f2cooowp")
			if err != nil {
				return err
			}

			// update field
			if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
				"hidden": false,
				"id": "he4cceqb",
				"maxSize": 5000000,
				"name": "outputs",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "json"
			}`)); err != nil {
				return err
			}

			if err := app.Save(collection); err != nil {
				return err
			}
		}

		// create collection `workflow_logs`
		{
			collection, err := app.FindCollectionByNameOrId("pbc_1682296116")
			if err != nil {
				return err
			}

			// update field
			if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "text3065852031",
				"max": 20000,
				"min": 0,
				"name": "message",
				"pattern": "",
				"presentable": false,
				"primaryKey": false,
				"required": false,
				"system": false,
				"type": "text"
			}`)); err != nil {
				return err
			}

			// update field
			if err := collection.Fields.AddMarshaledJSONAt(8, []byte(`{
				"hidden": false,
				"id": "json2918445923",
				"maxSize": 5000000,
				"name": "data",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "json"
			}`)); err != nil {
				return err
			}

			if err := app.Save(collection); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		return nil
	})
}

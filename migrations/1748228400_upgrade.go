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

			// add field
			if err := collection.Fields.AddMarshaledJSONAt(14, []byte(`{
				"hidden": false,
				"id": "bool810050391",
				"name": "acmeRenewed",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "bool"
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

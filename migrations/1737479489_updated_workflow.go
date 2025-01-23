package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("tovyif5ax6j62ur")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(10, []byte(`{
			"hidden": false,
			"id": "zivdxh23",
			"maxSelect": 1,
			"name": "lastRunStatus",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"pending",
				"running",
				"succeeded",
				"failed",
				"canceled"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("tovyif5ax6j62ur")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(10, []byte(`{
			"hidden": false,
			"id": "zivdxh23",
			"maxSelect": 1,
			"name": "lastRunStatus",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"pending",
				"running",
				"succeeded",
				"failed"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}

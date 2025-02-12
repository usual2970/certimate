package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("bqnxb95f2cooowp")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"cascadeDelete": true,
			"collectionId": "tovyif5ax6j62ur",
			"hidden": false,
			"id": "jka88auc",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "workflowId",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"cascadeDelete": true,
			"collectionId": "qjp8lygssgwyqyz",
			"hidden": false,
			"id": "relation821863227",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "runId",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("bqnxb95f2cooowp")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"cascadeDelete": false,
			"collectionId": "tovyif5ax6j62ur",
			"hidden": false,
			"id": "jka88auc",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "workflowId",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"cascadeDelete": false,
			"collectionId": "qjp8lygssgwyqyz",
			"hidden": false,
			"id": "relation821863227",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "runId",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}

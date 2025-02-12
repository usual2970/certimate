package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("bqnxb95f2cooowp")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE INDEX ` + "`" + `idx_BYoQPsz4my` + "`" + ` ON ` + "`" + `workflow_output` + "`" + ` (` + "`" + `workflowId` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_O9zxLETuxJ` + "`" + ` ON ` + "`" + `workflow_output` + "`" + ` (` + "`" + `runId` + "`" + `)"
			]
		}`), &collection); err != nil {
			return err
		}

		// add field
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
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("bqnxb95f2cooowp")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": []
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation821863227")

		return app.Save(collection)
	})
}

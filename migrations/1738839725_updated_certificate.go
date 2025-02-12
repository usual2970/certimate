package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("4szxr9x43tpj6np")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE INDEX ` + "`" + `idx_Jx8TXzDCmw` + "`" + ` ON ` + "`" + `certificate` + "`" + ` (` + "`" + `workflowId` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_kcKpgAZapk` + "`" + ` ON ` + "`" + `certificate` + "`" + ` (` + "`" + `workflowNodeId` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_2cRXqNDyyp` + "`" + ` ON ` + "`" + `certificate` + "`" + ` (` + "`" + `workflowRunId` + "`" + `)"
			]
		}`), &collection); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(15, []byte(`{
			"cascadeDelete": false,
			"collectionId": "qjp8lygssgwyqyz",
			"hidden": false,
			"id": "relation3917999135",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "workflowRunId",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("4szxr9x43tpj6np")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE INDEX ` + "`" + `idx_Jx8TXzDCmw` + "`" + ` ON ` + "`" + `certificate` + "`" + ` (` + "`" + `workflowId` + "`" + `)",
				"CREATE INDEX ` + "`" + `idx_kcKpgAZapk` + "`" + ` ON ` + "`" + `certificate` + "`" + ` (` + "`" + `workflowNodeId` + "`" + `)"
			]
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation3917999135")

		return app.Save(collection)
	})
}

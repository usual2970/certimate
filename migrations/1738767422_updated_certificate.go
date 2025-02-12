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
				"CREATE INDEX ` + "`" + `idx_kcKpgAZapk` + "`" + ` ON ` + "`" + `certificate` + "`" + ` (` + "`" + `workflowNodeId` + "`" + `)"
			]
		}`), &collection); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text2069360702",
			"max": 0,
			"min": 0,
			"name": "serialNumber",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(6, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text2910474005",
			"max": 0,
			"min": 0,
			"name": "issuer",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(8, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text4164403445",
			"max": 0,
			"min": 0,
			"name": "keyAlgorithm",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(11, []byte(`{
			"autogeneratePattern": "",
			"hidden": false,
			"id": "text2045248758",
			"max": 0,
			"min": 0,
			"name": "acmeAccountUrl",
			"pattern": "",
			"presentable": false,
			"primaryKey": false,
			"required": false,
			"system": false,
			"type": "text"
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
			"indexes": []
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("text2069360702")

		// remove field
		collection.Fields.RemoveById("text2910474005")

		// remove field
		collection.Fields.RemoveById("text4164403445")

		// remove field
		collection.Fields.RemoveById("text2045248758")

		return app.Save(collection)
	})
}

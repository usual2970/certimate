package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qjp8lygssgwyqyz")
		if err != nil {
			return err
		}

		collection.Name = "workflow_run"

		// add
		new_trigger := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "jlroa3fk",
			"name": "trigger",
			"type": "select",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"auto",
					"manual"
				]
			}
		}`), new_trigger); err != nil {
			return err
		}
		collection.Schema.AddField(new_trigger)

		// add
		new_startedAt := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "k9xvtf89",
			"name": "startedAt",
			"type": "date",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_startedAt); err != nil {
			return err
		}
		collection.Schema.AddField(new_startedAt)

		// add
		new_endedAt := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "3ikum7mk",
			"name": "endedAt",
			"type": "date",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_endedAt); err != nil {
			return err
		}
		collection.Schema.AddField(new_endedAt)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qjp8lygssgwyqyz")
		if err != nil {
			return err
		}

		collection.Name = "workflow_run_log"

		// remove
		collection.Schema.RemoveField("jlroa3fk")

		// remove
		collection.Schema.RemoveField("k9xvtf89")

		// remove
		collection.Schema.RemoveField("3ikum7mk")

		return dao.SaveCollection(collection)
	})
}

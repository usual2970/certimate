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

		// remove
		collection.Schema.RemoveField("cht6kqw9")

		// add
		new_status := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "qldmh0tw",
			"name": "status",
			"type": "select",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"pending",
					"running",
					"succeeded",
					"failed"
				]
			}
		}`), new_status); err != nil {
			return err
		}
		collection.Schema.AddField(new_status)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qjp8lygssgwyqyz")
		if err != nil {
			return err
		}

		// add
		del_succeeded := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "cht6kqw9",
			"name": "succeeded",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), del_succeeded); err != nil {
			return err
		}
		collection.Schema.AddField(del_succeeded)

		// remove
		collection.Schema.RemoveField("qldmh0tw")

		return dao.SaveCollection(collection)
	})
}

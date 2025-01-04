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

		collection, err := dao.FindCollectionByNameOrId("4szxr9x43tpj6np")
		if err != nil {
			return err
		}

		// add
		new_source := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "by9hetqi",
			"name": "source",
			"type": "select",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"workflow",
					"upload"
				]
			}
		}`), new_source); err != nil {
			return err
		}
		collection.Schema.AddField(new_source)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("4szxr9x43tpj6np")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("by9hetqi")

		return dao.SaveCollection(collection)
	})
}

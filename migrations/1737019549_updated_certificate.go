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
		new_deleted := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "klyf4nlq",
			"name": "deleted",
			"type": "date",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_deleted); err != nil {
			return err
		}
		collection.Schema.AddField(new_deleted)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("4szxr9x43tpj6np")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("klyf4nlq")

		return dao.SaveCollection(collection)
	})
}

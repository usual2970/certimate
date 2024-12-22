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

		collection, err := dao.FindCollectionByNameOrId("z3p974ainxjqlvs")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("1bspzuku")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z3p974ainxjqlvs")
		if err != nil {
			return err
		}

		// add
		del_group := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "1bspzuku",
			"name": "group",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "teolp9pl72dxlxq",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), del_group); err != nil {
			return err
		}
		collection.Schema.AddField(del_group)

		return dao.SaveCollection(collection)
	})
}

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
		collection.Schema.RemoveField("ghtlkn5j")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("z3p974ainxjqlvs")
		if err != nil {
			return err
		}

		// add
		del_lastDeployment := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "ghtlkn5j",
			"name": "lastDeployment",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "0a1o4e6sstp694f",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), del_lastDeployment); err != nil {
			return err
		}
		collection.Schema.AddField(del_lastDeployment)

		return dao.SaveCollection(collection)
	})
}

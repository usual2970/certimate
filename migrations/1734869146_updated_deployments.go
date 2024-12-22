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

		collection, err := dao.FindCollectionByNameOrId("0a1o4e6sstp694f")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("farvlzk7")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("0a1o4e6sstp694f")
		if err != nil {
			return err
		}

		// add
		del_domain := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "farvlzk7",
			"name": "domain",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "z3p974ainxjqlvs",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), del_domain); err != nil {
			return err
		}
		collection.Schema.AddField(del_domain)

		return dao.SaveCollection(collection)
	})
}

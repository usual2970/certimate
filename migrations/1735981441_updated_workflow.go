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

		collection, err := dao.FindCollectionByNameOrId("tovyif5ax6j62ur")
		if err != nil {
			return err
		}

		// add
		new_lastRunId := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "a23wkj9x",
			"name": "lastRunId",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "qjp8lygssgwyqyz",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), new_lastRunId); err != nil {
			return err
		}
		collection.Schema.AddField(new_lastRunId)

		// add
		new_lastRunStatus := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "zivdxh23",
			"name": "lastRunStatus",
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
		}`), new_lastRunStatus); err != nil {
			return err
		}
		collection.Schema.AddField(new_lastRunStatus)

		// add
		new_lastRunTime := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "u9bosu36",
			"name": "lastRunTime",
			"type": "date",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_lastRunTime); err != nil {
			return err
		}
		collection.Schema.AddField(new_lastRunTime)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("tovyif5ax6j62ur")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("a23wkj9x")

		// remove
		collection.Schema.RemoveField("zivdxh23")

		// remove
		collection.Schema.RemoveField("u9bosu36")

		return dao.SaveCollection(collection)
	})
}

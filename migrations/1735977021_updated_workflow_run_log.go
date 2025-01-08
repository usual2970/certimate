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

		// update
		edit_workflowId := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "m8xfsyyy",
			"name": "workflowId",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "tovyif5ax6j62ur",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), edit_workflowId); err != nil {
			return err
		}
		collection.Schema.AddField(edit_workflowId)

		// update
		edit_logs := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "2m9byaa9",
			"name": "logs",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_logs); err != nil {
			return err
		}
		collection.Schema.AddField(edit_logs)

		// update
		edit_succeeded := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "cht6kqw9",
			"name": "succeeded",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), edit_succeeded); err != nil {
			return err
		}
		collection.Schema.AddField(edit_succeeded)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("qjp8lygssgwyqyz")
		if err != nil {
			return err
		}

		// update
		edit_workflowId := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "m8xfsyyy",
			"name": "workflow",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "tovyif5ax6j62ur",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), edit_workflowId); err != nil {
			return err
		}
		collection.Schema.AddField(edit_workflowId)

		// update
		edit_logs := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "2m9byaa9",
			"name": "log",
			"type": "json",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSize": 2000000
			}
		}`), edit_logs); err != nil {
			return err
		}
		collection.Schema.AddField(edit_logs)

		// update
		edit_succeeded := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "cht6kqw9",
			"name": "succeed",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), edit_succeeded); err != nil {
			return err
		}
		collection.Schema.AddField(edit_succeeded)

		return dao.SaveCollection(collection)
	})
}

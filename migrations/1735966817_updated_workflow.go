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

		// update
		edit_trigger := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "vqoajwjq",
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
		}`), edit_trigger); err != nil {
			return err
		}
		collection.Schema.AddField(edit_trigger)

		// update
		edit_triggerCron := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "8ho247wh",
			"name": "triggerCron",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_triggerCron); err != nil {
			return err
		}
		collection.Schema.AddField(edit_triggerCron)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("tovyif5ax6j62ur")
		if err != nil {
			return err
		}

		// update
		edit_trigger := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "vqoajwjq",
			"name": "type",
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
		}`), edit_trigger); err != nil {
			return err
		}
		collection.Schema.AddField(edit_trigger)

		// update
		edit_triggerCron := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "8ho247wh",
			"name": "crontab",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_triggerCron); err != nil {
			return err
		}
		collection.Schema.AddField(edit_triggerCron)

		return dao.SaveCollection(collection)
	})
}

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
		new_effectAt := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "v40aqzpd",
			"name": "effectAt",
			"type": "date",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_effectAt); err != nil {
			return err
		}
		collection.Schema.AddField(new_effectAt)

		// update
		edit_subjectAltNames := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "fugxf58p",
			"name": "subjectAltNames",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_subjectAltNames); err != nil {
			return err
		}
		collection.Schema.AddField(edit_subjectAltNames)

		// update
		edit_acmeCertUrl := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "ayyjy5ve",
			"name": "acmeCertUrl",
			"type": "url",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"exceptDomains": null,
				"onlyDomains": null
			}
		}`), edit_acmeCertUrl); err != nil {
			return err
		}
		collection.Schema.AddField(edit_acmeCertUrl)

		// update
		edit_acmeCertStableUrl := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "3x5heo8e",
			"name": "acmeCertStableUrl",
			"type": "url",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"exceptDomains": null,
				"onlyDomains": null
			}
		}`), edit_acmeCertStableUrl); err != nil {
			return err
		}
		collection.Schema.AddField(edit_acmeCertStableUrl)

		// update
		edit_workflowId := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "uvqfamb1",
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
		edit_workflowNodeId := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "uqldzldw",
			"name": "workflowNodeId",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_workflowNodeId); err != nil {
			return err
		}
		collection.Schema.AddField(edit_workflowNodeId)

		// update
		edit_workflowOutputId := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "2ohlr0yd",
			"name": "workflowOutputId",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "bqnxb95f2cooowp",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), edit_workflowOutputId); err != nil {
			return err
		}
		collection.Schema.AddField(edit_workflowOutputId)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("4szxr9x43tpj6np")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("v40aqzpd")

		// update
		edit_subjectAltNames := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "fugxf58p",
			"name": "san",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_subjectAltNames); err != nil {
			return err
		}
		collection.Schema.AddField(edit_subjectAltNames)

		// update
		edit_acmeCertUrl := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "ayyjy5ve",
			"name": "certUrl",
			"type": "url",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"exceptDomains": null,
				"onlyDomains": null
			}
		}`), edit_acmeCertUrl); err != nil {
			return err
		}
		collection.Schema.AddField(edit_acmeCertUrl)

		// update
		edit_acmeCertStableUrl := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "3x5heo8e",
			"name": "certStableUrl",
			"type": "url",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"exceptDomains": null,
				"onlyDomains": null
			}
		}`), edit_acmeCertStableUrl); err != nil {
			return err
		}
		collection.Schema.AddField(edit_acmeCertStableUrl)

		// update
		edit_workflowId := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "uvqfamb1",
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
		edit_workflowNodeId := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "uqldzldw",
			"name": "nodeId",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_workflowNodeId); err != nil {
			return err
		}
		collection.Schema.AddField(edit_workflowNodeId)

		// update
		edit_workflowOutputId := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "2ohlr0yd",
			"name": "output",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "bqnxb95f2cooowp",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": 1,
				"displayFields": null
			}
		}`), edit_workflowOutputId); err != nil {
			return err
		}
		collection.Schema.AddField(edit_workflowOutputId)

		return dao.SaveCollection(collection)
	})
}

package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// update collection `certificate`
		{
			certimateCollection, err := app.FindCollectionByNameOrId("4szxr9x43tpj6np")
			if err != nil {
				return err
			}

			if err := certimateCollection.Fields.AddMarshaledJSONAt(4, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "plmambpz",
				"max": 100000,
				"min": 0,
				"name": "certificate",
				"pattern": "",
				"presentable": false,
				"primaryKey": false,
				"required": false,
				"system": false,
				"type": "text"
			}`)); err != nil {
				return err
			}

			if err := certimateCollection.Fields.AddMarshaledJSONAt(5, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "49qvwxcg",
				"max": 100000,
				"min": 0,
				"name": "privateKey",
				"pattern": "",
				"presentable": false,
				"primaryKey": false,
				"required": false,
				"system": false,
				"type": "text"
			}`)); err != nil {
				return err
			}

			if err := certimateCollection.Fields.AddMarshaledJSONAt(7, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "agt7n5bb",
				"max": 100000,
				"min": 0,
				"name": "issuerCertificate",
				"pattern": "",
				"presentable": false,
				"primaryKey": false,
				"required": false,
				"system": false,
				"type": "text"
			}`)); err != nil {
				return err
			}

			if err := app.Save(certimateCollection); err != nil {
				return err
			}
		}

		// update collection `workflow`
		{
			workflowCollection, err := app.FindCollectionByNameOrId("tovyif5ax6j62ur")
			if err != nil {
				return err
			}

			if err := workflowCollection.Fields.AddMarshaledJSONAt(6, []byte(`{
				"hidden": false,
				"id": "awlphkfe",
				"maxSize": 5000000,
				"name": "content",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "json"
			}`)); err != nil {
				return err
			}

			if err := workflowCollection.Fields.AddMarshaledJSONAt(7, []byte(`{
				"hidden": false,
				"id": "g9ohkk5o",
				"maxSize": 5000000,
				"name": "draft",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "json"
			}`)); err != nil {
				return err
			}

			if err := app.Save(workflowCollection); err != nil {
				return err
			}
		}

		// update collection `workflow_output`
		{
			workflowOutputCollection, err := app.FindCollectionByNameOrId("bqnxb95f2cooowp")
			if err != nil {
				return err
			}

			if err := workflowOutputCollection.Fields.AddMarshaledJSONAt(4, []byte(`{
				"hidden": false,
				"id": "c2rm9omj",
				"maxSize": 5000000,
				"name": "node",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "json"
			}`)); err != nil {
				return err
			}

			if err := app.Save(workflowOutputCollection); err != nil {
				return err
			}
		}

		return nil
	}, func(app core.App) error {
		return nil
	})
}

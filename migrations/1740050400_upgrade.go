package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		tracer := NewTracer("(v0.3)1740050400")
		tracer.Printf("go ...")

		// update collection `certificate`
		{
			collection, err := app.FindCollectionByNameOrId("4szxr9x43tpj6np")
			if err != nil {
				return err
			}

			if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
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

			if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
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

			if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
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

			if err := app.Save(collection); err != nil {
				return err
			}

			tracer.Printf("collection '%s' updated", collection.Name)
		}

		// update collection `workflow`
		{
			collection, err := app.FindCollectionByNameOrId("tovyif5ax6j62ur")
			if err != nil {
				return err
			}

			if err := collection.Fields.AddMarshaledJSONAt(6, []byte(`{
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

			if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
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

			if err := app.Save(collection); err != nil {
				return err
			}

			tracer.Printf("collection '%s' updated", collection.Name)
		}

		// update collection `workflow_output`
		{
			collection, err := app.FindCollectionByNameOrId("bqnxb95f2cooowp")
			if err != nil {
				return err
			}

			if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
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

			if err := app.Save(collection); err != nil {
				return err
			}

			tracer.Printf("collection '%s' updated", collection.Name)
		}

		tracer.Printf("done")
		return nil
	}, func(app core.App) error {
		return nil
	})
}

package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		tracer := NewTracer("(v0.3)1742644800")
		tracer.Printf("go ...")

		// update collection `workflow_run`
		{
			collection, err := app.FindCollectionByNameOrId("qjp8lygssgwyqyz")
			if err != nil {
				return err
			}

			// update field
			if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "hvebkuxw",
				"max": 20000,
				"min": 0,
				"name": "error",
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

		// update collection `workflow_output`
		{
			collection, err := app.FindCollectionByNameOrId("bqnxb95f2cooowp")
			if err != nil {
				return err
			}

			// update field
			if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
				"hidden": false,
				"id": "he4cceqb",
				"maxSize": 5000000,
				"name": "outputs",
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

		// update collection `workflow_logs`
		{
			collection, err := app.FindCollectionByNameOrId("pbc_1682296116")
			if err != nil {
				return err
			}

			// update field
			if err := collection.Fields.AddMarshaledJSONAt(7, []byte(`{
				"autogeneratePattern": "",
				"hidden": false,
				"id": "text3065852031",
				"max": 20000,
				"min": 0,
				"name": "message",
				"pattern": "",
				"presentable": false,
				"primaryKey": false,
				"required": false,
				"system": false,
				"type": "text"
			}`)); err != nil {
				return err
			}

			// update field
			if err := collection.Fields.AddMarshaledJSONAt(8, []byte(`{
				"hidden": false,
				"id": "json2918445923",
				"maxSize": 5000000,
				"name": "data",
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

		// update collection `access`
		{
			collection, err := app.FindCollectionByNameOrId("4yzbv8urny5ja1e")
			if err != nil {
				return err
			}

			// update field
			if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
				"hidden": false,
				"id": "hwy7m03o",
				"maxSelect": 1,
				"name": "provider",
				"presentable": false,
				"required": false,
				"system": false,
				"type": "select",
				"values": [
					"1panel",
					"acmehttpreq",
					"akamai",
					"aliyun",
					"aws",
					"azure",
					"baiducloud",
					"baishan",
					"baotapanel",
					"byteplus",
					"cachefly",
					"cdnfly",
					"cloudflare",
					"cloudns",
					"cmcccloud",
					"ctcccloud",
					"cucccloud",
					"desec",
					"dnsla",
					"dogecloud",
					"dynv6",
					"edgio",
					"fastly",
					"gname",
					"gcore",
					"godaddy",
					"goedge",
					"huaweicloud",
					"jdcloud",
					"k8s",
					"local",
					"namecheap",
					"namedotcom",
					"namesilo",
					"ns1",
					"porkbun",
					"powerdns",
					"qiniu",
					"qingcloud",
					"rainyun",
					"safeline",
					"ssh",
					"tencentcloud",
					"ucloud",
					"upyun",
					"vercel",
					"volcengine",
					"webhook",
					"westcn"
				]
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

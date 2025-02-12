package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
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
				"acmehttpreq",
				"aliyun",
				"aws",
				"azure",
				"baiducloud",
				"baotapanel",
				"byteplus",
				"cloudflare",
				"dogecloud",
				"godaddy",
				"huaweicloud",
				"k8s",
				"local",
				"namedotcom",
				"namesilo",
				"powerdns",
				"qiniu",
				"ssh",
				"tencentcloud",
				"ucloud",
				"volcengine",
				"webhook"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
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
				"acmehttpreq",
				"aliyun",
				"aws",
				"azure",
				"baiducloud",
				"byteplus",
				"cloudflare",
				"dogecloud",
				"godaddy",
				"huaweicloud",
				"k8s",
				"local",
				"namedotcom",
				"namesilo",
				"powerdns",
				"qiniu",
				"ssh",
				"tencentcloud",
				"ucloud",
				"volcengine",
				"webhook"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}

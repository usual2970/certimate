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
				"powerdns",
				"qiniu",
				"qingcloud",
				"rainyun",
				"safeline",
				"ssh",
				"tencentcloud",
				"ucloud",
				"upyun",
				"volcengine",
				"webhook",
				"westcn"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		return nil
	})
}

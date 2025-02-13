package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("4yzbv8urny5ja1e")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE INDEX `+"`"+`idx_wkoST0j`+"`"+` ON `+"`"+`access`+"`"+` (`+"`"+`name`+"`"+`)",
				"CREATE INDEX `+"`"+`idx_frh0JT1Aqx`+"`"+` ON `+"`"+`access`+"`"+` (`+"`"+`provider`+"`"+`)"
			]
		}`), &collection); err != nil {
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
				"dogecloud",
				"edgio",
				"fastly",
				"gname",
				"gcore",
				"godaddy",
				"goedge",
				"huaweicloud",
				"k8s",
				"local",
				"namedotcom",
				"namesilo",
				"ns1",
				"powerdns",
				"qiniu",
				"rainyun",
				"safeline",
				"ssh",
				"tencentcloud",
				"ucloud",
				"volcengine",
				"webhook",
				"westcn"
			]
		}`)); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("hsxcnlvd")

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("4yzbv8urny5ja1e")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE INDEX `+"`"+`idx_wkoST0j`+"`"+` ON `+"`"+`access`+"`"+` (`+"`"+`name`+"`"+`)"
			]
		}`), &collection); err != nil {
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

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"hidden": false,
			"id": "hsxcnlvd",
			"maxSelect": 1,
			"name": "usage",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"apply",
				"deploy",
				"all"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}

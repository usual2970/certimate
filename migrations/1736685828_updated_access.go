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

		collection, err := dao.FindCollectionByNameOrId("4yzbv8urny5ja1e")
		if err != nil {
			return err
		}

		// update
		edit_provider := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "hwy7m03o",
			"name": "provider",
			"type": "select",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
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
					"volcengine",
					"webhook"
				]
			}
		}`), edit_provider); err != nil {
			return err
		}
		collection.Schema.AddField(edit_provider)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("4yzbv8urny5ja1e")
		if err != nil {
			return err
		}

		// update
		edit_provider := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "hwy7m03o",
			"name": "provider",
			"type": "select",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"acmehttpreq",
					"aliyun",
					"aws",
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
					"volcengine",
					"webhook"
				]
			}
		}`), edit_provider); err != nil {
			return err
		}
		collection.Schema.AddField(edit_provider)

		return dao.SaveCollection(collection)
	})
}

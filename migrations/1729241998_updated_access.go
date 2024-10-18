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
		edit_configType := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "hwy7m03o",
			"name": "configType",
			"type": "select",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"aliyun",
					"tencent",
					"huaweicloud",
					"qiniu",
					"aws",
					"cloudflare",
					"namesilo",
					"godaddy",
					"local",
					"ssh",
					"webhook",
					"k8s"
				]
			}
		}`), edit_configType); err != nil {
			return err
		}
		collection.Schema.AddField(edit_configType)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("4yzbv8urny5ja1e")
		if err != nil {
			return err
		}

		// update
		edit_configType := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "hwy7m03o",
			"name": "configType",
			"type": "select",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 1,
				"values": [
					"aliyun",
					"tencent",
					"huaweicloud",
					"qiniu",
					"aws",
					"cloudflare",
					"namesilo",
					"godaddy",
					"local",
					"ssh",
					"webhook"
				]
			}
		}`), edit_configType); err != nil {
			return err
		}
		collection.Schema.AddField(edit_configType)

		return dao.SaveCollection(collection)
	})
}

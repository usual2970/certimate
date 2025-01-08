package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("teolp9pl72dxlxq")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	}, func(db dbx.Builder) error {
		jsonData := `{
			"id": "teolp9pl72dxlxq",
			"created": "2024-09-13 12:51:05.611Z",
			"updated": "2024-11-25 07:41:10.235Z",
			"name": "access_groups",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "7sajiv6i",
					"name": "name",
					"type": "text",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "xp8admif",
					"name": "access",
					"type": "relation",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"collectionId": "4yzbv8urny5ja1e",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": null,
						"displayFields": null
					}
				}
			],
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_RgRXp0R` + "`" + ` ON ` + "`" + `access_groups` + "`" + ` (` + "`" + `name` + "`" + `)"
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	})
}

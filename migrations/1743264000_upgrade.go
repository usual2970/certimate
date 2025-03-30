package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// update collection `settings`
		{
			collection, err := app.FindCollectionByNameOrId("dy6ccjb60spfy6p")
			if err != nil {
				return err
			}

			records, err := app.FindRecordsByFilter(collection, "name='sslProvider'", "-created", 1, 0)
			if err != nil {
				return err
			}

			if len(records) == 1 {
				record := records[0]

				content := make(map[string]any)
				if err := record.UnmarshalJSONField("content", &content); err != nil {
					return err
				}

				if provider, ok := content["provider"]; ok {
					if providerStr, ok := provider.(string); ok {
						if providerStr == "letsencrypt_staging" {
							content["provider"] = "letsencryptstaging"
						}
					}
				}

				if config, ok := content["config"]; ok {
					if configMap, ok := config.(map[string]any); ok {
						if _, ok := configMap["letsencrypt_staging"]; ok {
							configMap["letsencryptstaging"] = configMap["letsencrypt_staging"]
							delete(configMap, "letsencrypt_staging")
						}
						if _, ok := configMap["gts"]; ok {
							configMap["googletrustservices"] = configMap["gts"]
							delete(configMap, "gts")
						}
					}
				}

				record.Set("content", content)
				if err := app.Save(record); err != nil {
					return err
				}
			}
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
					"buypass",
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
					"googletrustservices",
					"huaweicloud",
					"jdcloud",
					"k8s",
					"letsencrypt",
					"letsencryptstaging",
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
					"sslcom",
					"tencentcloud",
					"ucloud",
					"upyun",
					"vercel",
					"volcengine",
					"webhook",
					"westcn",
					"zerossl"
				]
			}`)); err != nil {
				return err
			}

			if err := app.Save(collection); err != nil {
				return err
			}
		}

		// update collection `acme_accounts`
		{
			collection, err := app.FindCollectionByNameOrId("012d7abbod1hwvr")
			if err != nil {
				return err
			}

			records, err := app.FindRecordsByFilter(collection, "ca='letsencrypt_staging' || ca='gts'", "-created", 0, 0)
			if err != nil {
				return err
			}

			for _, record := range records {
				ca := record.GetString("ca")
				if ca == "letsencrypt_staging" {
					record.Set("ca", "letsencryptstaging")
				} else if ca == "gts" {
					record.Set("ca", "googletrustservices")
				} else {
					continue
				}

				if err := app.Save(record); err != nil {
					return err
				}
			}
		}

		return nil
	}, func(app core.App) error {
		return nil
	})
}

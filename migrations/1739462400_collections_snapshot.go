package migrations

import (
	x509 "crypto/x509"
	"log/slog"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	certutil "github.com/usual2970/certimate/internal/pkg/utils/cert"
)

func init() {
	m.Register(func(app core.App) error {
		slog.Info("[CERTIMATE] migration: ready ...")

		// backup collection records
		collectionRecords := make([]*core.Record, 0)
		collections, err := app.FindAllCollections(core.CollectionTypeBase)
		if err != nil {
			return err
		} else {
			for _, collection := range collections {
				switch collection.Name {
				case "acme_accounts", "access", "certificate", "workflow", "settings":
					{
						records, err := app.FindAllRecords(collection)
						if err != nil {
							return err
						}
						collectionRecords = append(collectionRecords, records...)

						slog.Info("[CERTIMATE] migration: collection '" + collection.Name + "' backed up")

						if collection.Name == "access" {
							collection.Fields.RemoveByName("usage")

							for i, field := range collection.Fields {
								if field.GetName() == "provider" {
									collection.Fields.AddMarshaledJSONAt(i+1, []byte(`{
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
                      "volcengine",
                      "webhook",
                      "westcn"
                    ]
                  }`))
								}
							}

							err := app.Save(collection)
							if err != nil {
								return err
							}
						}
					}

				case "domains", "deployments", "access_groups":
					{
						app.Delete(collection)

						slog.Info("[CERTIMATE] migration: collection '" + collection.Name + "' truncated")
					}
				}
			}
		}

		// migrate
		jsonData := `[
			{
				"createRule": null,
				"deleteRule": null,
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "geeur58v",
						"max": 0,
						"min": 0,
						"name": "name",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
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
							"volcengine",
							"webhook",
							"westcn"
						]
					},
					{
						"hidden": false,
						"id": "iql7jpwx",
						"maxSize": 2000000,
						"name": "config",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "json"
					},
					{
						"hidden": false,
						"id": "lr33hiwg",
						"max": "",
						"min": "",
						"name": "deleted",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "date"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": false,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": false,
						"type": "autodate"
					}
				],
				"id": "4yzbv8urny5ja1e",
				"indexes": [
					"CREATE INDEX ` + "`" + `idx_wkoST0j` + "`" + ` ON ` + "`" + `access` + "`" + ` (` + "`" + `name` + "`" + `)",
					"CREATE INDEX ` + "`" + `idx_frh0JT1Aqx` + "`" + ` ON ` + "`" + `access` + "`" + ` (` + "`" + `provider` + "`" + `)"
				],
				"listRule": null,
				"name": "access",
				"system": false,
				"type": "base",
				"updateRule": null,
				"viewRule": null
			},
			{
				"createRule": null,
				"deleteRule": null,
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "1tcmdsdf",
						"max": 0,
						"min": 0,
						"name": "name",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "f9wyhypi",
						"maxSize": 2000000,
						"name": "content",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "json"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": false,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": false,
						"type": "autodate"
					}
				],
				"id": "dy6ccjb60spfy6p",
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_RO7X9Vw` + "`" + ` ON ` + "`" + `settings` + "`" + ` (` + "`" + `name` + "`" + `)"
				],
				"listRule": null,
				"name": "settings",
				"system": false,
				"type": "base",
				"updateRule": null,
				"viewRule": null
			},
			{
				"createRule": null,
				"deleteRule": null,
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "fmjfn0yw",
						"max": 0,
						"min": 0,
						"name": "ca",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"exceptDomains": null,
						"hidden": false,
						"id": "qqwijqzt",
						"name": "email",
						"onlyDomains": null,
						"presentable": false,
						"required": false,
						"system": false,
						"type": "email"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "genxqtii",
						"max": 0,
						"min": 0,
						"name": "key",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "1aoia909",
						"maxSize": 2000000,
						"name": "resource",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "json"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": false,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": false,
						"type": "autodate"
					}
				],
				"id": "012d7abbod1hwvr",
				"indexes": [],
				"listRule": null,
				"name": "acme_accounts",
				"system": false,
				"type": "base",
				"updateRule": null,
				"viewRule": null
			},
			{
				"createRule": null,
				"deleteRule": null,
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "8yydhv1h",
						"max": 0,
						"min": 0,
						"name": "name",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "1buzebwz",
						"max": 0,
						"min": 0,
						"name": "description",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "vqoajwjq",
						"maxSelect": 1,
						"name": "trigger",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "select",
						"values": [
							"auto",
							"manual"
						]
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "8ho247wh",
						"max": 0,
						"min": 0,
						"name": "triggerCron",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "nq7kfdzi",
						"name": "enabled",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "bool"
					},
					{
						"hidden": false,
						"id": "awlphkfe",
						"maxSize": 2000000,
						"name": "content",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "json"
					},
					{
						"hidden": false,
						"id": "g9ohkk5o",
						"maxSize": 2000000,
						"name": "draft",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "json"
					},
					{
						"hidden": false,
						"id": "2rpfz9t3",
						"name": "hasDraft",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "bool"
					},
					{
						"cascadeDelete": false,
						"collectionId": "qjp8lygssgwyqyz",
						"hidden": false,
						"id": "a23wkj9x",
						"maxSelect": 1,
						"minSelect": 0,
						"name": "lastRunId",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "relation"
					},
					{
						"hidden": false,
						"id": "zivdxh23",
						"maxSelect": 1,
						"name": "lastRunStatus",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "select",
						"values": [
							"pending",
							"running",
							"succeeded",
							"failed",
							"canceled"
						]
					},
					{
						"hidden": false,
						"id": "u9bosu36",
						"max": "",
						"min": "",
						"name": "lastRunTime",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "date"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": false,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": false,
						"type": "autodate"
					}
				],
				"id": "tovyif5ax6j62ur",
				"indexes": [],
				"listRule": null,
				"name": "workflow",
				"system": false,
				"type": "base",
				"updateRule": null,
				"viewRule": null
			},
			{
				"createRule": null,
				"deleteRule": null,
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"cascadeDelete": true,
						"collectionId": "tovyif5ax6j62ur",
						"hidden": false,
						"id": "jka88auc",
						"maxSelect": 1,
						"minSelect": 0,
						"name": "workflowId",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "relation"
					},
					{
						"cascadeDelete": true,
						"collectionId": "qjp8lygssgwyqyz",
						"hidden": false,
						"id": "relation821863227",
						"maxSelect": 1,
						"minSelect": 0,
						"name": "runId",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "relation"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "z9fgvqkz",
						"max": 0,
						"min": 0,
						"name": "nodeId",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "c2rm9omj",
						"maxSize": 2000000,
						"name": "node",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "json"
					},
					{
						"hidden": false,
						"id": "he4cceqb",
						"maxSize": 2000000,
						"name": "outputs",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "json"
					},
					{
						"hidden": false,
						"id": "2yfxbxuf",
						"name": "succeeded",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "bool"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": false,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": false,
						"type": "autodate"
					}
				],
				"id": "bqnxb95f2cooowp",
				"indexes": [
					"CREATE INDEX ` + "`" + `idx_BYoQPsz4my` + "`" + ` ON ` + "`" + `workflow_output` + "`" + ` (` + "`" + `workflowId` + "`" + `)",
					"CREATE INDEX ` + "`" + `idx_O9zxLETuxJ` + "`" + ` ON ` + "`" + `workflow_output` + "`" + ` (` + "`" + `runId` + "`" + `)"
				],
				"listRule": null,
				"name": "workflow_output",
				"system": false,
				"type": "base",
				"updateRule": null,
				"viewRule": null
			},
			{
				"createRule": null,
				"deleteRule": null,
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "by9hetqi",
						"maxSelect": 1,
						"name": "source",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "select",
						"values": [
							"workflow",
							"upload"
						]
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "fugxf58p",
						"max": 0,
						"min": 0,
						"name": "subjectAltNames",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text2069360702",
						"max": 0,
						"min": 0,
						"name": "serialNumber",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "plmambpz",
						"max": 0,
						"min": 0,
						"name": "certificate",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "49qvwxcg",
						"max": 0,
						"min": 0,
						"name": "privateKey",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text2910474005",
						"max": 0,
						"min": 0,
						"name": "issuer",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "agt7n5bb",
						"max": 0,
						"min": 0,
						"name": "issuerCertificate",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text4164403445",
						"max": 0,
						"min": 0,
						"name": "keyAlgorithm",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "v40aqzpd",
						"max": "",
						"min": "",
						"name": "effectAt",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "date"
					},
					{
						"hidden": false,
						"id": "zgpdby2k",
						"max": "",
						"min": "",
						"name": "expireAt",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "date"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text2045248758",
						"max": 0,
						"min": 0,
						"name": "acmeAccountUrl",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"exceptDomains": null,
						"hidden": false,
						"id": "ayyjy5ve",
						"name": "acmeCertUrl",
						"onlyDomains": null,
						"presentable": false,
						"required": false,
						"system": false,
						"type": "url"
					},
					{
						"exceptDomains": null,
						"hidden": false,
						"id": "3x5heo8e",
						"name": "acmeCertStableUrl",
						"onlyDomains": null,
						"presentable": false,
						"required": false,
						"system": false,
						"type": "url"
					},
					{
						"cascadeDelete": false,
						"collectionId": "tovyif5ax6j62ur",
						"hidden": false,
						"id": "uvqfamb1",
						"maxSelect": 1,
						"minSelect": 0,
						"name": "workflowId",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "relation"
					},
					{
						"cascadeDelete": false,
						"collectionId": "qjp8lygssgwyqyz",
						"hidden": false,
						"id": "relation3917999135",
						"maxSelect": 1,
						"minSelect": 0,
						"name": "workflowRunId",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "relation"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "uqldzldw",
						"max": 0,
						"min": 0,
						"name": "workflowNodeId",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"cascadeDelete": false,
						"collectionId": "bqnxb95f2cooowp",
						"hidden": false,
						"id": "2ohlr0yd",
						"maxSelect": 1,
						"minSelect": 0,
						"name": "workflowOutputId",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "relation"
					},
					{
						"hidden": false,
						"id": "klyf4nlq",
						"max": "",
						"min": "",
						"name": "deleted",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "date"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": false,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": false,
						"type": "autodate"
					}
				],
				"id": "4szxr9x43tpj6np",
				"indexes": [
					"CREATE INDEX ` + "`" + `idx_Jx8TXzDCmw` + "`" + ` ON ` + "`" + `certificate` + "`" + ` (` + "`" + `workflowId` + "`" + `)",
					"CREATE INDEX ` + "`" + `idx_kcKpgAZapk` + "`" + ` ON ` + "`" + `certificate` + "`" + ` (` + "`" + `workflowNodeId` + "`" + `)",
					"CREATE INDEX ` + "`" + `idx_2cRXqNDyyp` + "`" + ` ON ` + "`" + `certificate` + "`" + ` (` + "`" + `workflowRunId` + "`" + `)"
				],
				"listRule": null,
				"name": "certificate",
				"system": false,
				"type": "base",
				"updateRule": null,
				"viewRule": null
			},
			{
				"createRule": null,
				"deleteRule": null,
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"cascadeDelete": true,
						"collectionId": "tovyif5ax6j62ur",
						"hidden": false,
						"id": "m8xfsyyy",
						"maxSelect": 1,
						"minSelect": 0,
						"name": "workflowId",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "relation"
					},
					{
						"hidden": false,
						"id": "qldmh0tw",
						"maxSelect": 1,
						"name": "status",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "select",
						"values": [
							"pending",
							"running",
							"succeeded",
							"failed",
							"canceled"
						]
					},
					{
						"hidden": false,
						"id": "jlroa3fk",
						"maxSelect": 1,
						"name": "trigger",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "select",
						"values": [
							"auto",
							"manual"
						]
					},
					{
						"hidden": false,
						"id": "k9xvtf89",
						"max": "",
						"min": "",
						"name": "startedAt",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "date"
					},
					{
						"hidden": false,
						"id": "3ikum7mk",
						"max": "",
						"min": "",
						"name": "endedAt",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "date"
					},
					{
						"hidden": false,
						"id": "2m9byaa9",
						"maxSize": 2000000,
						"name": "logs",
						"presentable": false,
						"required": false,
						"system": false,
						"type": "json"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "hvebkuxw",
						"max": 0,
						"min": 0,
						"name": "error",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": false,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": false,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": false,
						"type": "autodate"
					}
				],
				"id": "qjp8lygssgwyqyz",
				"indexes": [],
				"listRule": null,
				"name": "workflow_run",
				"system": false,
				"type": "base",
				"updateRule": null,
				"viewRule": null
			},
			{
				"authAlert": {
					"emailTemplate": {
						"body": "<p>Hello,</p>\n<p>We noticed a login to your {APP_NAME} account from a new location.</p>\n<p>If this was you, you may disregard this email.</p>\n<p><strong>If this wasn't you, you should immediately change your {APP_NAME} account password to revoke access from all other locations.</strong></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
						"subject": "Login from a new location"
					},
					"enabled": true
				},
				"authRule": "",
				"authToken": {
					"duration": 1209600
				},
				"confirmEmailChangeTemplate": {
					"body": "<p>Hello,</p>\n<p>Click on the button below to confirm your new email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Confirm new email</a>\n</p>\n<p><i>If you didn't ask to change your email address, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
					"subject": "Confirm your {APP_NAME} new email address"
				},
				"createRule": null,
				"deleteRule": null,
				"emailChangeToken": {
					"duration": 1800
				},
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"cost": 0,
						"hidden": true,
						"id": "password901924565",
						"max": 0,
						"min": 8,
						"name": "password",
						"pattern": "",
						"presentable": false,
						"required": true,
						"system": true,
						"type": "password"
					},
					{
						"autogeneratePattern": "[a-zA-Z0-9]{50}",
						"hidden": true,
						"id": "text2504183744",
						"max": 60,
						"min": 30,
						"name": "tokenKey",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"exceptDomains": null,
						"hidden": false,
						"id": "email3885137012",
						"name": "email",
						"onlyDomains": null,
						"presentable": false,
						"required": true,
						"system": true,
						"type": "email"
					},
					{
						"hidden": false,
						"id": "bool1547992806",
						"name": "emailVisibility",
						"presentable": false,
						"required": false,
						"system": true,
						"type": "bool"
					},
					{
						"hidden": false,
						"id": "bool256245529",
						"name": "verified",
						"presentable": false,
						"required": false,
						"system": true,
						"type": "bool"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": true,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": true,
						"type": "autodate"
					}
				],
				"fileToken": {
					"duration": 120
				},
				"id": "pbc_3142635823",
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_tokenKey_pbc_3142635823` + "`" + ` ON ` + "`" + `_superusers` + "`" + ` (` + "`" + `tokenKey` + "`" + `)",
					"CREATE UNIQUE INDEX ` + "`" + `idx_email_pbc_3142635823` + "`" + ` ON ` + "`" + `_superusers` + "`" + ` (` + "`" + `email` + "`" + `) WHERE ` + "`" + `email` + "`" + ` != ''"
				],
				"listRule": null,
				"manageRule": null,
				"mfa": {
					"duration": 1800,
					"enabled": false,
					"rule": ""
				},
				"name": "_superusers",
				"oauth2": {
					"enabled": false,
					"mappedFields": {
						"avatarURL": "",
						"id": "",
						"name": "",
						"username": ""
					}
				},
				"otp": {
					"duration": 180,
					"emailTemplate": {
						"body": "<p>Hello,</p>\n<p>Your one-time password is: <strong>{OTP}</strong></p>\n<p><i>If you didn't ask for the one-time password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
						"subject": "OTP for {APP_NAME}"
					},
					"enabled": false,
					"length": 8
				},
				"passwordAuth": {
					"enabled": true,
					"identityFields": [
						"email"
					]
				},
				"passwordResetToken": {
					"duration": 1800
				},
				"resetPasswordTemplate": {
					"body": "<p>Hello,</p>\n<p>Click on the button below to reset your password.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Reset password</a>\n</p>\n<p><i>If you didn't ask to reset your password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
					"subject": "Reset your {APP_NAME} password"
				},
				"system": true,
				"type": "auth",
				"updateRule": null,
				"verificationTemplate": {
					"body": "<p>Hello,</p>\n<p>Thank you for joining us at {APP_NAME}.</p>\n<p>Click on the button below to verify your email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Verify</a>\n</p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
					"subject": "Verify your {APP_NAME} email"
				},
				"verificationToken": {
					"duration": 259200
				},
				"viewRule": null
			},
			{
				"createRule": null,
				"deleteRule": "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId",
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text455797646",
						"max": 0,
						"min": 0,
						"name": "collectionRef",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text127846527",
						"max": 0,
						"min": 0,
						"name": "recordRef",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text2462348188",
						"max": 0,
						"min": 0,
						"name": "provider",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text1044722854",
						"max": 0,
						"min": 0,
						"name": "providerId",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": true,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": true,
						"type": "autodate"
					}
				],
				"id": "pbc_2281828961",
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_externalAuths_record_provider` + "`" + ` ON ` + "`" + `_externalAuths` + "`" + ` (collectionRef, recordRef, provider)",
					"CREATE UNIQUE INDEX ` + "`" + `idx_externalAuths_collection_provider` + "`" + ` ON ` + "`" + `_externalAuths` + "`" + ` (collectionRef, provider, providerId)"
				],
				"listRule": "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId",
				"name": "_externalAuths",
				"system": true,
				"type": "base",
				"updateRule": null,
				"viewRule": "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId"
			},
			{
				"createRule": null,
				"deleteRule": null,
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text455797646",
						"max": 0,
						"min": 0,
						"name": "collectionRef",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text127846527",
						"max": 0,
						"min": 0,
						"name": "recordRef",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text1582905952",
						"max": 0,
						"min": 0,
						"name": "method",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": true,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": true,
						"type": "autodate"
					}
				],
				"id": "pbc_2279338944",
				"indexes": [
					"CREATE INDEX ` + "`" + `idx_mfas_collectionRef_recordRef` + "`" + ` ON ` + "`" + `_mfas` + "`" + ` (collectionRef,recordRef)"
				],
				"listRule": "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId",
				"name": "_mfas",
				"system": true,
				"type": "base",
				"updateRule": null,
				"viewRule": "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId"
			},
			{
				"createRule": null,
				"deleteRule": null,
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text455797646",
						"max": 0,
						"min": 0,
						"name": "collectionRef",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text127846527",
						"max": 0,
						"min": 0,
						"name": "recordRef",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"cost": 8,
						"hidden": true,
						"id": "password901924565",
						"max": 0,
						"min": 0,
						"name": "password",
						"pattern": "",
						"presentable": false,
						"required": true,
						"system": true,
						"type": "password"
					},
					{
						"autogeneratePattern": "",
						"hidden": true,
						"id": "text3866985172",
						"max": 0,
						"min": 0,
						"name": "sentTo",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": false,
						"system": true,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": true,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": true,
						"type": "autodate"
					}
				],
				"id": "pbc_1638494021",
				"indexes": [
					"CREATE INDEX ` + "`" + `idx_otps_collectionRef_recordRef` + "`" + ` ON ` + "`" + `_otps` + "`" + ` (collectionRef, recordRef)"
				],
				"listRule": "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId",
				"name": "_otps",
				"system": true,
				"type": "base",
				"updateRule": null,
				"viewRule": "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId"
			},
			{
				"createRule": null,
				"deleteRule": "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId",
				"fields": [
					{
						"autogeneratePattern": "[a-z0-9]{15}",
						"hidden": false,
						"id": "text3208210256",
						"max": 15,
						"min": 15,
						"name": "id",
						"pattern": "^[a-z0-9]+$",
						"presentable": false,
						"primaryKey": true,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text455797646",
						"max": 0,
						"min": 0,
						"name": "collectionRef",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text127846527",
						"max": 0,
						"min": 0,
						"name": "recordRef",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"autogeneratePattern": "",
						"hidden": false,
						"id": "text4228609354",
						"max": 0,
						"min": 0,
						"name": "fingerprint",
						"pattern": "",
						"presentable": false,
						"primaryKey": false,
						"required": true,
						"system": true,
						"type": "text"
					},
					{
						"hidden": false,
						"id": "autodate2990389176",
						"name": "created",
						"onCreate": true,
						"onUpdate": false,
						"presentable": false,
						"system": true,
						"type": "autodate"
					},
					{
						"hidden": false,
						"id": "autodate3332085495",
						"name": "updated",
						"onCreate": true,
						"onUpdate": true,
						"presentable": false,
						"system": true,
						"type": "autodate"
					}
				],
				"id": "pbc_4275539003",
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_authOrigins_unique_pairs` + "`" + ` ON ` + "`" + `_authOrigins` + "`" + ` (collectionRef, recordRef, fingerprint)"
				],
				"listRule": "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId",
				"name": "_authOrigins",
				"system": true,
				"type": "base",
				"updateRule": null,
				"viewRule": "@request.auth.id != '' && recordRef = @request.auth.id && collectionRef = @request.auth.collectionId"
			}
		]`
		err = app.ImportCollectionsByMarshaledJSON([]byte(jsonData), false)
		if err != nil {
			return err
		}

		slog.Info("[CERTIMATE] migration: collections imported")

		// restore records
		for _, record := range collectionRecords {
			changed := false

			switch record.Collection().Name {
			case "access":
				{
					if record.GetString("provider") == "tencent" {
						record.Set("provider", "tencentcloud")
						changed = true
					} else if record.GetString("provider") == "pdns" {
						record.Set("provider", "powerdns")
						changed = true
					} else if record.GetString("provider") == "httpreq" {
						record.Set("provider", "acmehttpreq")
						changed = true
					}
				}

			case "certificate":
				{
					if record.GetString("issuer") == "" {
						cert, _ := certutil.ParseCertificateFromPEM(record.GetString("certificate"))
						if cert != nil {
							record.Set("issuer", strings.Join(cert.Issuer.Organization, ";"))
							changed = true
						}
					}
					if record.GetString("serialNumber") == "" {
						cert, _ := certutil.ParseCertificateFromPEM(record.GetString("certificate"))
						if cert != nil {
							record.Set("serialNumber", strings.ToUpper(cert.SerialNumber.Text(16)))
							changed = true
						}
					}
					if record.GetString("keyAlgorithm") == "" {
						cert, _ := certutil.ParseCertificateFromPEM(record.GetString("certificate"))
						if cert != nil {
							switch cert.SignatureAlgorithm {
							case x509.SHA256WithRSA, x509.SHA256WithRSAPSS:
								record.Set("keyAlgorithm", "RSA2048")
							case x509.SHA384WithRSA, x509.SHA384WithRSAPSS:
								record.Set("keyAlgorithm", "RSA3072")
							case x509.SHA512WithRSA, x509.SHA512WithRSAPSS:
								record.Set("keyAlgorithm", "RSA4096")
							case x509.ECDSAWithSHA256:
								record.Set("keyAlgorithm", "EC256")
							case x509.ECDSAWithSHA384:
								record.Set("keyAlgorithm", "EC384")
							case x509.ECDSAWithSHA512:
								record.Set("keyAlgorithm", "EC512")
							}
							changed = true
						}
					}
				}
			}

			if changed {
				err = app.Save(record)
				if err != nil {
					return err
				}

				slog.Info("[CERTIMATE] migration: collection '" + record.Collection().Name + "' record #" + record.Id + " updated")
			}
		}

		slog.Info("[CERTIMATE] migration: done")

		return nil
	}, func(app core.App) error {
		return nil
	})
}

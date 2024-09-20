package migrations

import (
	"github.com/pocketbase/dbx"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		// add up queries...
		db.NewQuery("update access set usage='all' where configType in ('aliyun', 'tencent','aws')").Execute()
		db.NewQuery("update access set usage='deploy' where configType in ('ssh', 'webhook','qiniu')").Execute()
		db.NewQuery("update access set usage='apply' where configType in ('cloudflare','namesilo','godaddy')").Execute()
		return nil
	}, func(db dbx.Builder) error {
		// add down queries...

		return nil
	})
}

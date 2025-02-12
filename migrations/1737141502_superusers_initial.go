package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		if err != nil {
			return err
		}

		record, _ := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, "admin@certimate.fun")
		if record == nil {
			record := core.NewRecord(superusers)
			record.Set("email", "admin@certimate.fun")
			record.Set("password", "1234567890")
			return app.Save(record)
		}

		return nil
	}, func(app core.App) error {
		return nil
	})
}

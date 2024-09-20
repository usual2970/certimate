package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		// add up queries...
		dao := daos.New(db)

		admin := &models.Admin{}
		admin.Email = "admin@devops.com"
		admin.SetPassword("1234567890")
		return dao.SaveAdmin(admin)
	}, func(db dbx.Builder) error {
		// add down queries...
		dao := daos.New(db)

		admin, _ := dao.FindAdminByEmail("admin@devops.com")
		if admin != nil {
			return dao.DeleteAdmin(admin)
		}
		return nil
	})
}

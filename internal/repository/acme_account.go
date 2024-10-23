package repository

import (
	"github.com/go-acme/lego/v4/registration"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/usual2970/certimate/internal/domain"
	"github.com/usual2970/certimate/internal/utils/app"
)

type AcmeAccountRepository struct{}

func NewAcmeAccountRepository() *AcmeAccountRepository {
	return &AcmeAccountRepository{}
}

func (r *AcmeAccountRepository) GetByCAAndEmail(ca, email string) (*domain.AcmeAccount, error) {
	resp, err := app.GetApp().Dao().FindFirstRecordByFilter("acme_accounts", "ca={:ca} && email={:email}", dbx.Params{"ca": ca, "email": email})
	if err != nil {
		return nil, err
	}

	resource := &registration.Resource{}
	if err := resp.UnmarshalJSONField("resource", resource); err != nil {
		return nil, err
	}

	return &domain.AcmeAccount{
		Id:       resp.GetString("id"),
		Ca:       resp.GetString("ca"),
		Email:    resp.GetString("email"),
		Key:      resp.GetString("key"),
		Resource: resource,
		Created:  resp.GetTime("created"),
		Updated:  resp.GetTime("updated"),
	}, nil
}

func (r *AcmeAccountRepository) Save(ca, email, key string, resource *registration.Resource) error {
	collection, err := app.GetApp().Dao().FindCollectionByNameOrId("acme_accounts")
	if err != nil {
		return err
	}

	record := models.NewRecord(collection)
	record.Set("ca", ca)
	record.Set("email", email)
	record.Set("key", key)
	record.Set("resource", resource)
	return app.GetApp().Dao().Save(record)
}

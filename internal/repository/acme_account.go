package repository

import (
	"fmt"

	"github.com/go-acme/lego/v4/registration"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/sync/singleflight"

	"github.com/usual2970/certimate/internal/app"
	"github.com/usual2970/certimate/internal/domain"
)

type AcmeAccountRepository struct{}

func NewAcmeAccountRepository() *AcmeAccountRepository {
	return &AcmeAccountRepository{}
}

var g singleflight.Group

func (r *AcmeAccountRepository) GetByCAAndEmail(ca, email string) (*domain.AcmeAccount, error) {
	resp, err, _ := g.Do(fmt.Sprintf("acme_account_%s_%s", ca, email), func() (interface{}, error) {
		resp, err := app.GetApp().FindFirstRecordByFilter(
			"acme_accounts",
			"ca={:ca} && email={:email}",
			dbx.Params{"ca": ca, "email": email},
		)
		if err != nil {
			return nil, err
		}
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, domain.ErrRecordNotFound
	}

	record, ok := resp.(*core.Record)
	if !ok {
		return nil, domain.ErrRecordNotFound
	}

	return r.castRecordToModel(record)
}

func (r *AcmeAccountRepository) Save(ca, email, key string, resource *registration.Resource) error {
	collection, err := app.GetApp().FindCollectionByNameOrId("acme_accounts")
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	record.Set("ca", ca)
	record.Set("email", email)
	record.Set("key", key)
	record.Set("resource", resource)
	return app.GetApp().Save(record)
}

func (r *AcmeAccountRepository) castRecordToModel(record *core.Record) (*domain.AcmeAccount, error) {
	if record == nil {
		return nil, fmt.Errorf("record is nil")
	}

	resource := &registration.Resource{}
	if err := record.UnmarshalJSONField("resource", resource); err != nil {
		return nil, err
	}

	acmeAccount := &domain.AcmeAccount{
		Meta: domain.Meta{
			Id:        record.Id,
			CreatedAt: record.GetDateTime("created").Time(),
			UpdatedAt: record.GetDateTime("updated").Time(),
		},
		CA:       record.GetString("ca"),
		Email:    record.GetString("email"),
		Key:      record.GetString("key"),
		Resource: resource,
	}
	return acmeAccount, nil
}

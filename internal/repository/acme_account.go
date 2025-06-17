package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-acme/lego/v4/registration"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/sync/singleflight"

	"github.com/certimate-go/certimate/internal/app"
	"github.com/certimate-go/certimate/internal/domain"
)

type AcmeAccountRepository struct{}

func NewAcmeAccountRepository() *AcmeAccountRepository {
	return &AcmeAccountRepository{}
}

var g singleflight.Group

func (r *AcmeAccountRepository) GetByCAAndEmail(ca, email string) (*domain.AcmeAccount, error) {
	resp, err, _ := g.Do(fmt.Sprintf("acme_account_%s_%s", ca, email), func() (interface{}, error) {
		resp, err := app.GetApp().FindFirstRecordByFilter(
			domain.CollectionNameAcmeAccount,
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

func (r *AcmeAccountRepository) Save(ctx context.Context, acmeAccount *domain.AcmeAccount) (*domain.AcmeAccount, error) {
	collection, err := app.GetApp().FindCollectionByNameOrId(domain.CollectionNameAcmeAccount)
	if err != nil {
		return acmeAccount, err
	}

	var record *core.Record
	if acmeAccount.Id == "" {
		record = core.NewRecord(collection)
	} else {
		record, err = app.GetApp().FindRecordById(collection, acmeAccount.Id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return acmeAccount, domain.ErrRecordNotFound
			}
			return acmeAccount, err
		}
	}

	record.Set("ca", acmeAccount.CA)
	record.Set("email", acmeAccount.Email)
	record.Set("key", acmeAccount.Key)
	record.Set("resource", acmeAccount.Resource)
	if err := app.GetApp().Save(record); err != nil {
		return acmeAccount, err
	}

	acmeAccount.Id = record.Id
	acmeAccount.CreatedAt = record.GetDateTime("created").Time()
	acmeAccount.UpdatedAt = record.GetDateTime("updated").Time()
	return acmeAccount, nil
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

package domain

import (
	"github.com/go-acme/lego/v4/registration"
)

type AcmeAccount struct {
	Meta
	CA       string                 `json:"ca" db:"ca"`
	Email    string                 `json:"email" db:"email"`
	Resource *registration.Resource `json:"resource" db:"resource"`
	Key      string                 `json:"key" db:"key"`
}

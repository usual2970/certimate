package domain

import (
	"time"

	"github.com/go-acme/lego/v4/registration"
)

type AcmeAccount struct {
	Id       string
	Ca       string
	Email    string
	Resource *registration.Resource
	Key      string
	Created  time.Time
	Updated  time.Time
}

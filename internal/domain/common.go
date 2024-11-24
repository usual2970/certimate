package domain

import "time"

type Meta struct {
	Id      string    `json:"id" db:"id"`
	Created time.Time `json:"created" db:"created"`
	Updated time.Time `json:"updated" db:"updated"`
}

package domain

import "time"

type Meta struct {
	Id        string    `json:"id" db:"id"`
	CreatedAt time.Time `json:"created" db:"created"`
	UpdatedAt time.Time `json:"updated" db:"updated"`
}

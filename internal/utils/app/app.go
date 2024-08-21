package app

import (
	"sync"

	"github.com/pocketbase/pocketbase"
)

var instance *pocketbase.PocketBase

var intanceOnce sync.Once

func GetApp() *pocketbase.PocketBase {
	intanceOnce.Do(func() {
		instance = pocketbase.New()
	})

	return instance
}

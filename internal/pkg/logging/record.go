package logging

import (
	"time"

	types "github.com/pocketbase/pocketbase/tools/types"
)

type Record struct {
	Time    time.Time
	Level   Level
	Message string
	Data    types.JSONMap[any]
}

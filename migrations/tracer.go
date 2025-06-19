package migrations

import (
	"fmt"
	"log/slog"
)

type Tracer struct {
	logger *slog.Logger
	flag   string
}

func NewTracer(flag string) *Tracer {
	return &Tracer{
		logger: slog.Default(),
		flag:   flag,
	}
}

func (l *Tracer) Printf(format string, args ...any) {
	l.logger.Info("[CERTIMATE] migration " + l.flag + ": " + fmt.Sprintf(format, args...))
}

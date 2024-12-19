package log

import (
	"log/slog"
	"os"
)

type Log = *slog.Logger

func New() Log {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
}

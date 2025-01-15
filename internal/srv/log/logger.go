package log

import (
	"log/slog"
	"os"
)

type Log = *slog.Logger

func New(version string, prodEnv bool) Log {
	lvl := new(slog.LevelVar)
	if prodEnv {
		lvl.Set(slog.LevelInfo)
	} else {
		lvl.Set(slog.LevelDebug)
	}
	Log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	slog.SetDefault(Log)
	return Log.With("version", version, "app", "memopass")
}

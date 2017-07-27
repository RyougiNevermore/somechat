package log

import (
	slog "log"
	"sync"
	"os"
)

var logger *slog.Logger
var once sync.Once

func Log() *slog.Logger {
	once.Do(func() {
		logger = slog.New(os.Stdout, "", 0)
	})
	return logger
}
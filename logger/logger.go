package logger

import (
	"github.com/pterm/pterm"
)

var Log *pterm.Logger

func SetLogLevel(level pterm.LogLevel) {
	Log = pterm.DefaultLogger.WithLevel(level)
}

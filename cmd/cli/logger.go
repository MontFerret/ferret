package cli

import (
	"github.com/natefinch/lumberjack"
)

func NewLogger() *lumberjack.Logger {
	l := &lumberjack.Logger{
		Filename: "./ferret.log",
		MaxSize:  100,
	}

	return l
}

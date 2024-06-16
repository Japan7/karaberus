package server

import (
	"log"
	"os"
)

var _logger *log.Logger = nil

func getLogger() *log.Logger {
	if _logger == nil {
		_logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return _logger
}

func Warn(msg string) {
	getLogger().Printf("Warning: %s\n", msg)
}

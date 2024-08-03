package server

import (
	"io"
	"log"
	"os"
)

var _logger *log.Logger = nil

func getLogger() *log.Logger {
	if _logger == nil {
		_logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return _logger
}

func Warn(msg string) {
	getLogger().Printf("Warning: %s\n", msg)
}

func Closer(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		Warn(err.Error())
	}
}

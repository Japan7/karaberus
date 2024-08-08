package main

import (
	"C"
	"github.com/Japan7/karaberus/server"
)

//export MakeCli
func MakeCli() {
	server.MakeCli()
}

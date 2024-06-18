package main

import (
	"log"
	"os"

	"github.com/Japan7/karaberus/server"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "create-token":
			token, err := server.CreateSystemToken()
			if err != nil {
				log.Fatalln(err)
			}
			println(token)
		}
	} else {
		server.RunKaraberus()
	}
}

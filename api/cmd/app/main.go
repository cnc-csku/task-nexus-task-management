package main

import (
	"log"

	"github.com/cnc-csku/task-nexus/internal/wire"
)

func main() {
	app := wire.InitializeApp()

	err := app.Start()
	if err != nil {
		log.Fatalln(err)
	}
}

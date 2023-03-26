package main

import (
	"github.com/jattento/alien-invasion-simulator/cmd"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

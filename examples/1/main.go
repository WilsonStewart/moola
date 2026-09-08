package main

import (
	"log"

	"github.com/WilsonStewart/moola/internal/processor"
)

func main() {
	p := processor.NewProcessor()

	err := p.ReadAndProcessFile("master.moola")
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/WilsonStewart/moola/internal/processor"
)

func main() {
	err := processor.ProcessMoolaFile("master.moola", true)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/WilsonStewart/moola/internal/core"
	"github.com/WilsonStewart/moola/internal/processor"
)

func main() {
	mi := core.NewMoolaInstance()
	p := processor.NewProcessor(mi)

	err := p.ProcessMoolaFile("master.moola", true)
	if err != nil {
		log.Fatal(err)
	}
}

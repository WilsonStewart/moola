package processor

import (
	"time"

	"github.com/WilsonStewart/moola/internal/core"
)

type Processor struct {
	mi           *core.MoolaInstance
	currentYear  int
	currentMonth time.Month
	currentDay   int
}

type MoolaFile [][]string

type MoolaFileLine struct {
	directive  string
	arguments  []string
	filename   string
	lineNumber int
}

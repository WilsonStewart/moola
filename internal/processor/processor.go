package processor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/WilsonStewart/moola/internal/core"
)

type Processor struct {
	database     *core.MoolaDatabase
	currentYear  int
	currentMonth time.Month
	currentDay   int
}

type MoolaFile [][]string

func NewProcessor() *Processor {
	return &Processor{
		database: core.NewMoolaDatabase(),
	}
}

func readMoolaFile(path string) (MoolaFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result MoolaFile
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())

		if len(words) > 0 {
			result = append(result, words)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Processor) ProcessMoolaFile(path string, isMasterFile bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var mf MoolaFile
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())

		if len(words) > 0 {
			mf = append(mf, words)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, line := range mf {
		directive := strings.ToUpper(line[0])
		switch directive {
		case "OPENACCOUNT":
			if err := p.OPENACCOUNT(line); err != nil {
				return err
			}
		case "ASSERT":
			if err := p.ASSERT(line); err != nil {
				return err
			}
		case "TRANSACT":
			if err := p.TRANSACT(line); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%q is not a known directive!", directive)
		}
	}

	b, _ := json.MarshalIndent(p.database, "", "  ")
	fmt.Println(string(b))
	return nil
}

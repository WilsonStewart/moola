package processor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type SymbolsStore struct {
	AccountNames      []string
	AccountAliasNames []string
}

type Datafile struct {
	Accounts       map[string]*Account
	AccountAliases map[string]string
}

func NewDatafile() *Datafile {
	datafile := &Datafile{
		Accounts:       make(map[string]*Account),
		AccountAliases: make(map[string]string),
	}

	return datafile
}

type Processor struct {
	Datafile                   Datafile
	Symbols                    SymbolsStore
	currentDate                time.Time
	defaultEnvelopeAccountName string
}

func NewProcessor() *Processor {
	processor := &Processor{
		Datafile: *NewDatafile(),
	}

	processor.defaultEnvelopeAccountName = "to_be_assigned"

	defaultEnvelopeAccountAliasName := "tba"

	processor.OpenAccount(OpenAccountNode{
		accountKind:           "envelope",
		requestedAccountName:  "to_be_assigned",
		requestedAccountAlias: &defaultEnvelopeAccountAliasName,
	})

	return processor
}

func VerifyExpectedDirectiveAndFieldRange(
	rawFields []string,
	expectedDirective string,
	expectedMinimumFieldCount int,
	expectedMaximumFieldCount int,
) error {
	if rawFields[0] != expectedDirective {
		return fmt.Errorf("invalid directive %q: expected %q", rawFields[0], expectedDirective)
	}

	count := len(rawFields)
	if count < expectedMinimumFieldCount || count > expectedMaximumFieldCount {
		return fmt.Errorf("invalid field count %d for %q: expected between %d and %d",
			count, expectedDirective, expectedMinimumFieldCount, expectedMaximumFieldCount)
	}

	return nil
}

func (p *Processor) ReadAndProcessFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines [][]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		if len(fields) == 0 {
			continue
		}

		lines = append(lines, fields)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, line := range lines {
		switch strings.ToLower(line[0]) {
		case "openaccount":
			node := OpenAccountNode{}
			if err := node.Unmarshal(line); err != nil {
				return err
			}

			if err := p.OpenAccount(node); err != nil {
				return err
			}
		}
	}

	datafile, err := json.MarshalIndent(p.Datafile, "", "  ")
	if err != nil {
		panic(err)
	}
	symbols, err := json.MarshalIndent(p.Symbols, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(datafile))
	fmt.Println(string(symbols))

	return nil
}

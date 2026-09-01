package processor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/WilsonStewart/moola/internal/core"
)

func NewProcessor(moolaInstance *core.MoolaInstance) *Processor {
	return &Processor{
		mi: moolaInstance,
	}
}

// func readMoolaFile(path string) (MoolaFile, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	var result MoolaFile
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		words := strings.Fields(scanner.Text())

// 		if len(words) > 0 {
// 			result = append(result, words)
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }

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

	for idx, line := range mf {
		if len(line) == 0 {
			continue
		}

		if len(line) == 1 {
			return fmt.Errorf("this line is too short. must have at least one argument: %s", strings.Join(line, ","))
		}

		mfl := MoolaFileLine{
			directive:  strings.ToUpper(line[0]),
			arguments:  line[1:],
			lineNumber: idx,
			filename:   path,
		}
		switch mfl.directive {
		case "OPENACCOUNT":
			if err := p.OPENACCOUNT(mfl); err != nil {
				return err
			}
		case "ASSERT":
			if err := p.ASSERT(mfl); err != nil {
				return err
			}
		case "TRANSACT":
			if err := p.TRANSACT(mfl); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%q is not a known directive!", mfl.directive)
		}
	}

	b, _ := json.MarshalIndent(p.mi, "", "  ")
	fmt.Println(string(b))

	lipgloss.Println(p.mi.RenderCategorizedTables())
	return nil
}

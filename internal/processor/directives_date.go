package processor

import (
	"errors"
	"strconv"
)

func (p *Processor) SETYEAR(line []string) error {
	if len(line) < 2 {
		return errors.New("missing year argument")
	}

	yearInt, err := strconv.Atoi(line[1])
	if err != nil {
		return err
	}

	p.currentYear = yearInt

	return nil
}

func (p *Processor) SETMONTH(line []string) error {
	if len(line) < 2 {
		return errors.New("missing year argument")
	}

	yearInt, err := strconv.Atoi(line[1])
	if err != nil {
		return err
	}

	p.currentYear = yearInt

	return nil
}

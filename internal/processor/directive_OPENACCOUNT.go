package processor

import (
	"errors"
	"fmt"
)

func (p *Processor) OPENACCOUNT(line MoolaFileLine) error {
	if line.directive != "OPENACCOUNT" {
		return fmt.Errorf("tried to format a '%q' directive with the OPENACCOUNT method for some reason", line.directive)
	}

	if len(line.arguments) != 3 {
		return errors.New("3 arguments are required for the OPENACCOUNT directive: OPENACCOUNT name alias kind")
	}

	if err := p.mi.AddAccount(line.arguments[0], line.arguments[2]); err != nil {
		return fmt.Errorf("could not add account to database: %w", err)
	}

	return nil
}

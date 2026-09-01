package processor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (p *Processor) OPENACCOUNT(line []string) error {
	if strings.ToUpper(line[0]) != "OPENACCOUNT" {
		return fmt.Errorf("tried to format a '%q' directive with the OPENACCOUNT method for some reason", strings.ToUpper(line[0]))
	}

	if len(line) != 4 {
		return errors.New("4 arguments are required for the OPENACCOUNT directive: OPENACCOUNT name alias kind")
	}

	if err := p.database.AddAccount(line[1], line[3]); err != nil {
		return fmt.Errorf("could not add account to database: %w", err)
	}

	return nil
}

func (p *Processor) ASSERT(line []string) error {
	if strings.ToUpper(line[0]) != "ASSERT" {
		return fmt.Errorf("tried to format a '%q' directive with the ASSERT method for some reason", strings.ToUpper(line[0]))
	}

	if len(line) != 3 {
		return errors.New("3 arguments are required for the ASSERT directive: ASSERT accountName balance")
	}

	account, err := p.database.GetAccount(line[1])
	if err != nil {
		return fmt.Errorf("could not find account: %w", err)
	}

	newBalance, err := strconv.Atoi(line[2])
	if err != nil {
		return fmt.Errorf("could not convert %s to a balance int: %w", line[2], err)
	}

	account.Balance = newBalance

	return nil
}

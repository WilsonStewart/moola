package processor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (p *Processor) ASSERT(line MoolaFileLine) error {
	if line.directive != "ASSERT" {
		return fmt.Errorf("tried to format a '%q' directive with the ASSERT method for some reason", strings.ToUpper(line.directive))
	}

	if len(line.arguments) != 2 {
		return errors.New("2 arguments are required for the ASSERT directive: ASSERT accountName balance")
	}

	account, err := p.mi.GetAccount(line.arguments[0])
	if err != nil {
		return fmt.Errorf("could not find account: %w", err)
	}

	newBalance, err := strconv.Atoi(line.arguments[1])
	if err != nil {
		return fmt.Errorf("could not convert %s to a balance int: %w", line.arguments[1], err)
	}

	account.Balance = newBalance

	return nil
}

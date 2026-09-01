package processor

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func (p *Processor) TRANSACT(line []string) error {
	if strings.ToUpper(line[0]) != "TRANSACT" {
		return fmt.Errorf("tried to format a '%q' directive with the TRANSACT method for some reason", strings.ToUpper(line[0]))
	}

	if len(line) != 5 {
		return errors.New("5 arguments are required for the TRANSACT directive: TRANSACT payee amount realAccount envelope")
	}

	realAccount, err := p.database.GetAccount(line[3])
	if err != nil {
		return fmt.Errorf("'%q' real account not found", line[3])
	}
	validRealAccountKinds := []string{"cash", "creditcard"}
	if !slices.Contains(validRealAccountKinds, realAccount.Kind) {
		return fmt.Errorf("the provided real account is not of a valid kind, which are: %s", strings.Join(validRealAccountKinds, ","))
	}

	envelopeAccount, err := p.database.GetAccount(line[4])
	if err != nil {
		return fmt.Errorf("'%q' envelope not found", line[4])
	}
	if envelopeAccount.Kind != "envelope" {
		return fmt.Errorf("the provided envelope account was of the wrong kind")
	}

	balanceChange, err := strconv.Atoi(line[2])
	if err != nil {
		return fmt.Errorf("could not convert %q to int: %w", line[2], err)
	}

	realAccount.Balance += balanceChange
	envelopeAccount.Balance += balanceChange

	return nil
}

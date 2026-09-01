package processor

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func (p *Processor) TRANSACT(line MoolaFileLine) error {
	if line.directive != "TRANSACT" {
		return fmt.Errorf("tried to format a '%q' directive with the TRANSACT method for some reason", line.directive)
	}

	if len(line.arguments) != 4 {
		return errors.New("4 arguments are required for the TRANSACT directive: TRANSACT payee amount realAccount envelope")
	}

	realAccount, err := p.mi.GetAccount(line.arguments[2])
	if err != nil {
		return fmt.Errorf("'%q' real account not found", line.arguments[2])
	}
	validRealAccountKinds := []string{"cash", "creditcard"}
	if !slices.Contains(validRealAccountKinds, realAccount.Kind) {
		return fmt.Errorf("the provided real account is not of a valid kind, which are: %s", strings.Join(validRealAccountKinds, ","))
	}

	envelopeAccount, err := p.mi.GetAccount(line.arguments[3])
	if err != nil {
		return fmt.Errorf("'%q' envelope not found", line.arguments[3])
	}
	if envelopeAccount.Kind != "envelope" {
		return fmt.Errorf("the provided envelope account was of the wrong kind")
	}

	balanceChange, err := strconv.Atoi(line.arguments[1])
	if err != nil {
		return fmt.Errorf("could not convert %q to int: %w", line.arguments[1], err)
	}

	realAccount.Balance += balanceChange
	envelopeAccount.Balance += balanceChange

	return nil
}

package processor

import (
	"fmt"
	"strconv"
)

type TransactNode struct {
	BalanceDelta        float64
	Payee               string
	CashAccountName     string
	EnvelopeAccountName string
	Note                string
}

func (n *TransactNode) Unmarshal(rawFields []string) error {
	if err := VerifyExpectedDirectiveAndFieldRange(rawFields, "transact", 5, 6); err != nil {
		return err
	}

	parsedBalanceDelta, err := strconv.ParseFloat(rawFields[1], 64)
	if err != nil {
		return err
	}
	n.BalanceDelta = parsedBalanceDelta

	n.Payee = rawFields[2]
	n.CashAccountName = rawFields[3]
	n.EnvelopeAccountName = rawFields[4]

	if len(rawFields) == 6 {
		n.Note = rawFields[5]
	}

	return nil
}

func (p *Processor) Transact(node TransactNode) error {
	cashAccount, err := p.getAccount(node.CashAccountName)
	if err != nil {
		return err
	}

	envelopeAccount, err := p.getAccount(node.EnvelopeAccountName)
	if err != nil {
		return err
	}

	if cashAccount.Kind != "cash" {
		return fmt.Errorf("account-kind-mismatch: %s account provided, but cash required", cashAccount.Kind)
	}

	if envelopeAccount.Kind != "envelope" {
		return fmt.Errorf("account-kind-mismatch: %s account provided, but envelope required", envelopeAccount.Kind)
	}

	cashAccount.Balance += node.BalanceDelta
	envelopeAccount.Balance += node.BalanceDelta

	return nil
}

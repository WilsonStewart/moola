package processor

import (
	"fmt"
	"strconv"
)

type AssertNode struct {
	CashAccountName string
	Balance         float64
}

func (node *AssertNode) Unmarshal(rawFields []string) error {
	if err := VerifyExpectedDirectiveAndFieldRange(rawFields, "assert", 3, 3); err != nil {
		return err
	}

	node.CashAccountName = rawFields[1]

	parsedBalance, err := strconv.ParseFloat(rawFields[2], 64)
	if err != nil {
		return err
	}
	node.Balance = parsedBalance

	return nil
}

func (d *Datafile) Assert(node AssertNode) error {
	// Replace the alias name if it's there
	if realCashAccountName, ok := d.AccountAliases[node.CashAccountName]; ok {
		node.CashAccountName = realCashAccountName
	}

	account, ok := d.Accounts[node.CashAccountName]
	if !ok {
		return fmt.Errorf("account-not-found: %q not found", node.CashAccountName)
	}

	defaultEnvelopeAccount, ok := d.Accounts[d.defaultEnvelopeAccountName]
	if !ok {
		return fmt.Errorf("account-not-found: %q not found", d.defaultEnvelopeAccountName)
	}

	balanceDelta := node.Balance - account.Balance

	account.Balance += balanceDelta
	defaultEnvelopeAccount.Balance += balanceDelta

	return nil
}

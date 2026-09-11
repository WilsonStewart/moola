package processor

import (
	"fmt"
	"strconv"
)

type MoveNode struct {
	sourceEnvelopeAccountName string
	targetEnvelopeAccountName string
	balanceDelta              float64
}

func (n *MoveNode) Unmarshal(rawFields []string) error {
	if err := VerifyExpectedDirectiveAndFieldRange(rawFields, "move", 4, 4); err != nil {
		return err
	}

	n.sourceEnvelopeAccountName = rawFields[1]
	n.targetEnvelopeAccountName = rawFields[2]

	parsedBalanceDelta, err := strconv.ParseFloat(rawFields[3], 64)

	if parsedBalanceDelta < 0 {
		return fmt.Errorf("no-negative-allowed: negative numbers are not allowed in move")
	}

	if err != nil {
		return err
	}
	n.balanceDelta = parsedBalanceDelta

	return nil
}

func (d *Datafile) Move(node MoveNode) error {
	source, err := d.getAccount(node.sourceEnvelopeAccountName)
	if err != nil {
		return err
	}

	target, err := d.getAccount(node.targetEnvelopeAccountName)
	if err != nil {
		return err
	}

	if source.Kind != "envelope" {
		return fmt.Errorf("wrong-account-kind: got %q kind, expected envelope", source.Kind)
	}

	if target.Kind != "envelope" {
		return fmt.Errorf("wrong-account-kind: got %q kind, expected envelope", target.Kind)
	}

	source.Balance -= node.balanceDelta
	target.Balance += node.balanceDelta

	return nil
}

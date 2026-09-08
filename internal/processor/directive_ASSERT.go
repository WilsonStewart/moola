package processor

import "strconv"

type AssertNode struct {
	AccountName string
	Balance     float64
}

func (node *AssertNode) Unmarshal(rawFields []string) error {
	if err := VerifyExpectedDirectiveAndFieldRange(rawFields, "assert", 3, 3); err != nil {
		return err
	}

	node.AccountName = rawFields[1]

	parsedBalance, err := strconv.ParseFloat(rawFields[2], 64)
	if err != nil {
		return err
	}
	node.Balance = parsedBalance

	return nil
}

func (p *Processor) Assert(node AssertNode) error {

	return nil
}

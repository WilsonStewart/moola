package processor

import (
	"fmt"
	"slices"
)

type Account struct {
	Name    string
	Kind    string
	Balance float64
}

var validAccountKinds = []string{
	"cash",
	"envelope",
}

func (p *Processor) getAccount(name string) (*Account, error) {
	if slices.Contains(p.Symbols.AccountAliasNames, name) {
		name = p.Datafile.AccountAliases[name]
	}

	if !slices.Contains(p.Symbols.AccountNames, name) {
		return nil, fmt.Errorf("could not find account with name %q", name)
	}

	return p.Datafile.Accounts[name], nil
}

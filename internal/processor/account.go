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

func (d *Datafile) getAccount(name string) (*Account, error) {
	if slices.Contains(d.Symbols.AccountAliasNames, name) {
		name = d.AccountAliases[name]
	}

	if !slices.Contains(d.Symbols.AccountNames, name) {
		return nil, fmt.Errorf("could not find account with name %q", name)
	}

	return d.Accounts[name], nil
}

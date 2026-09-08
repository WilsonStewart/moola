package processor

import (
	"fmt"
	"slices"
	"strings"
)

type OpenAccountNode struct {
	accountKind           string
	requestedAccountName  string
	requestedAccountAlias *string
}

func (node *OpenAccountNode) Unmarshal(rawFields []string) error {
	if err := VerifyExpectedDirectiveAndFieldRange(rawFields, "openaccount", 3, 4); err != nil {
		return err
	}

	node.accountKind = rawFields[1]
	node.requestedAccountName = rawFields[2]
	if len(rawFields) == 4 {
		node.requestedAccountAlias = &rawFields[3]
	}

	return nil
}

func (p *Processor) OpenAccount(node OpenAccountNode) error {
	if !slices.Contains(validAccountKinds, node.accountKind) {
		return fmt.Errorf(
			"invalid type: %q is not a valid type out of %s",
			node.accountKind,
			strings.Join(validAccountKinds, ","),
		)
	}

	if slices.Contains(p.Symbols.AccountNames, node.requestedAccountName) {
		return fmt.Errorf(
			"duplicate symbol: cannot create account %q as an account with that name already exists!",
			node.requestedAccountName,
		)
	}

	if slices.Contains(p.Symbols.AccountAliasNames, node.requestedAccountName) {
		return fmt.Errorf(
			"duplicate symbol: cannot create account %q as an account alias with that name already exists!",
			node.requestedAccountName,
		)
	}

	p.Datafile.Accounts[node.requestedAccountName] = &Account{
		Name:    node.requestedAccountName,
		Kind:    node.accountKind,
		Balance: 0,
	}
	p.Symbols.AccountNames = append(p.Symbols.AccountNames, node.requestedAccountName)

	if node.requestedAccountAlias != nil {
		if slices.Contains(p.Symbols.AccountAliasNames, *node.requestedAccountAlias) {
			return fmt.Errorf(
				"duplicate symbol: cannot create account alias %q as an account alias with that name already exists!",
				*node.requestedAccountAlias,
			)
		}

		if slices.Contains(p.Symbols.AccountNames, *node.requestedAccountAlias) {
			return fmt.Errorf(
				"duplicate symbol: cannot create account alias %q as an account with that name already exists!",
				*node.requestedAccountAlias,
			)
		}
		p.Datafile.AccountAliases[*node.requestedAccountAlias] = node.requestedAccountName
		p.Symbols.AccountAliasNames = append(p.Symbols.AccountAliasNames, *node.requestedAccountAlias)
	}

	return nil
}

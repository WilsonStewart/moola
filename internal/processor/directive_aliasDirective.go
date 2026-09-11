package processor

import (
	"fmt"
	"slices"
)

type AliasDirectiveNode struct {
	DirectiveName string
	AliasName     string
}

func (n *AliasDirectiveNode) Unmarshal(rawFields []string) error {
	if err := VerifyExpectedDirectiveAndFieldRange(rawFields, "aliasdirective", 3, 3); err != nil {
		return err
	}

	n.DirectiveName = rawFields[1]
	n.AliasName = rawFields[2]

	return nil
}

func (d *Datafile) AliasDirective(node AliasDirectiveNode) error {
	if slices.Contains(d.Symbols.DirectiveAliasNames, node.AliasName) {
		return fmt.Errorf("duplicate alias: %q directive alias already exists", node.AliasName)
	}

	if slices.Contains(d.Symbols.DirectiveNames, node.AliasName) {
		return fmt.Errorf("invalid alias name: %q is already a directive, cannot make an alias with that name", node.AliasName)
	}

	if !slices.Contains(d.Symbols.DirectiveNames, node.DirectiveName) {
		return fmt.Errorf("invalid directive name: %q is not a known directive", node.DirectiveName)
	}

	d.DirectiveAliases[node.AliasName] = node.DirectiveName
	d.Symbols.DirectiveAliasNames = append(d.Symbols.DirectiveAliasNames, node.AliasName)

	return nil
}

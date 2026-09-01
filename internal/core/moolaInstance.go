package core

import (
	"fmt"
	"image/color"
	"slices"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

func NewMoolaInstance() *MoolaInstance {
	instance := &MoolaInstance{
		Accounts:       make(map[string]*Account),
		accountAliases: make(map[string]string),
	}

	return instance
}

func (mi *MoolaInstance) AddAccount(name string, kind string) error {
	mi.mu.Lock()
	defer mi.mu.Unlock()

	if _, exists := mi.Accounts[name]; exists {
		return fmt.Errorf("cannot create account '%q'  as an account of that name already exists", name)
	}

	if _, exists := mi.accountAliases[name]; exists {
		return fmt.Errorf("cannot create account '%q' as an alias of that name already exists", name)
	}

	validKinds := []string{"cash", "creditcard", "envelope"}
	if !slices.Contains(validKinds, kind) {
		return fmt.Errorf("the kind '%q' is not a valid choice. choose from: %s", kind, strings.Join(validKinds, ","))
	}

	mi.Accounts[name] = &Account{
		Name: name,
		Kind: kind,
	}

	return nil
}

func (mi *MoolaInstance) GetAccount(name string) (*Account, error) {
	mi.mu.RLock()
	defer mi.mu.RUnlock()

	if account, ok := mi.Accounts[name]; ok {
		return account, nil
	}
	if id, ok := mi.accountAliases[name]; ok {
		if account, ok := mi.Accounts[id]; ok {
			return account, nil
		}
	}
	return nil, fmt.Errorf("account %q not found", name)
}

func (mi *MoolaInstance) AccountExists(name string) bool {
	mi.mu.RLock()
	defer mi.mu.RUnlock()

	if _, ok := mi.Accounts[name]; ok {
		return true
	}
	if id, ok := mi.accountAliases[name]; ok {
		if _, ok := mi.Accounts[id]; ok {
			return true
		}
	}
	return false
}

func (m *MoolaInstance) RenderTable() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var rows [][]string
	for _, acc := range m.Accounts {
		// Assuming Balance is stored in cents (e.g. 1050 -> $10.50)
		formattedBalance := fmt.Sprintf("$%.2f", float64(acc.Balance)/100.0)

		rows = append(rows, []string{
			acc.Name,
			acc.Kind,
			formattedBalance,
		})
	}

	// Sort by account name for consistent display ordering
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	// Style definitions
	purple := lipgloss.Color("99")
	gray := lipgloss.Color("240")
	green := lipgloss.Color("42")

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(purple).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	balanceStyle := cellStyle.Copy().
		Foreground(green).
		Align(lipgloss.Right)

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(gray)).
		Headers("ACCOUNT NAME", "KIND", "BALANCE").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			if col == 2 {
				return balanceStyle
			}
			return cellStyle
		})

	return t.String()
}

func (m *MoolaInstance) RenderCategorizedTables() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Label style definition
	// makeLabel := func(title string, color color.Color) string {
	// 	return lipgloss.NewStyle().
	// 		Bold(true).
	// 		Foreground(color).
	// 		PaddingLeft(1).
	// 		Render(title)
	// }

	// Build individual sections
	// envelopeLabel := makeLabel("ENVELOPES", lipgloss.Color("99"))
	envelopeTable := renderKindTable(m.Accounts, "envelope", lipgloss.Color("99"))
	envelopeSection := lipgloss.JoinVertical(lipgloss.Left, envelopeTable)

	// cashLabel := makeLabel("CASH", lipgloss.Color("39"))
	cashTable := renderKindTable(m.Accounts, "cash", lipgloss.Color("39"))
	cashSection := lipgloss.JoinVertical(lipgloss.Left, cashTable)

	// Join both sections vertically with empty spacing between them
	return lipgloss.JoinVertical(
		lipgloss.Left,
		envelopeSection,
		"", // empty line spacer
		cashSection,
	)
}

func renderKindTable(accounts map[string]*Account, targetKind string, headerColor color.Color) string {
	var rows [][]string
	for _, acc := range accounts {
		if acc.Kind != targetKind {
			continue
		}
		formattedBalance := fmt.Sprintf("$%.2f", float64(acc.Balance)/100.0)
		rows = append(rows, []string{
			acc.Name,
			acc.Kind,
			formattedBalance,
		})
	}

	// Handle empty state gracefully
	if len(rows) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			PaddingLeft(1).
			Render("No accounts available")
	}

	// Sort rows alphabetically by account name
	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0]
	})

	gray := lipgloss.Color("240")
	green := lipgloss.Color("42")

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(headerColor).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	balanceStyle := cellStyle.Copy().
		Foreground(green).
		Align(lipgloss.Right)

	return table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(gray)).
		Headers("ACCOUNT NAME", "KIND", "BALANCE").
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			if col == 2 {
				return balanceStyle
			}
			return cellStyle
		}).
		String()
}

package core

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

type Account struct {
	Name                  string
	Kind                  string // cash, creditcard, envelope
	Balance               int
	ccBalanceNeedsCovered int
}

type MoolaDatabase struct {
	Accounts       map[string]*Account
	accountAliases map[string]string
	mu             sync.RWMutex
}

func NewMoolaDatabase() *MoolaDatabase {
	newDb := &MoolaDatabase{
		Accounts:       make(map[string]*Account),
		accountAliases: make(map[string]string),
	}

	return newDb
}

func (mdb *MoolaDatabase) AddAccount(name string, kind string) error {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()

	if _, exists := mdb.Accounts[name]; exists {
		return fmt.Errorf("cannot create account '%q'  as an account of that name already exists", name)
	}

	if _, exists := mdb.accountAliases[name]; exists {
		return fmt.Errorf("cannot create account '%q' as an alias of that name already exists", name)
	}

	validKinds := []string{"cash", "creditcard", "envelope"}
	if !slices.Contains(validKinds, kind) {
		return fmt.Errorf("the kind '%q' is not a valid choice. choose from: %s", kind, strings.Join(validKinds, ","))
	}

	mdb.Accounts[name] = &Account{
		Name: name,
		Kind: kind,
	}

	return nil
}

func (mdb *MoolaDatabase) GetAccount(name string) (*Account, error) {
	mdb.mu.RLock()
	defer mdb.mu.RUnlock()

	if account, ok := mdb.Accounts[name]; ok {
		return account, nil
	}
	if id, ok := mdb.accountAliases[name]; ok {
		if account, ok := mdb.Accounts[id]; ok {
			return account, nil
		}
	}
	return nil, fmt.Errorf("account %q not found", name)
}

func (mdb *MoolaDatabase) AccountExists(name string) bool {
	mdb.mu.RLock()
	defer mdb.mu.RUnlock()

	if _, ok := mdb.Accounts[name]; ok {
		return true
	}
	if id, ok := mdb.accountAliases[name]; ok {
		if _, ok := mdb.Accounts[id]; ok {
			return true
		}
	}
	return false
}

package core

import (
	"fmt"
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

func NewMoolaDatabase() (*MoolaDatabase, error) {
	var newDb MoolaDatabase

	newDb.Accounts = make(map[string]*Account)

	return ddnewDb, nil
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

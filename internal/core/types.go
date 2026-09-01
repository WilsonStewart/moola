package core

import "sync"

type Account struct {
	Name                  string
	Kind                  string // cash, creditcard, envelope
	Balance               int
	ccBalanceNeedsCovered int
}

type MoolaInstance struct {
	Accounts       map[string]*Account
	accountAliases map[string]string
	mu             sync.RWMutex
}

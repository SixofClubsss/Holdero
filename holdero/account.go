package holdero

import (
	dreams "github.com/dReam-dApps/dReams"
)

// Holdero account data
type accountData struct {
	Key    string       `json:"key"`
	Tables []string     `json:"tables"`
	Stats  Player_stats `json:"stats"`
}

// Get Holdero account data
func GetAccount() interface{} {
	return accountData{
		Key:    handKey,
		Tables: table.Favorites.SCIDs,
		Stats:  stats,
	}
}

// Set stored Holdero account data to variables
func SetAccount(ad interface{}) (err error) {
	var account accountData
	err = dreams.SetAccount(ad, &account)
	if err != nil {
		logger.Errorln("[SetAccount]", err)
		clearAccountData()
	} else {
		table.Favorites.SCIDs = account.Tables
		stats = account.Stats
		if account.Key != "" {
			handKey = account.Key
		}
	}

	return
}

// Clear existing Holdero account data
func clearAccountData() {
	table.Favorites.SCIDs = []string{}
	stats = Player_stats{}
}

// Save Holdero account data to datashards
func saveAccount() *dreams.AccountEncrypted {
	return dreams.AddAccountData(GetAccount(), "holdero")
}

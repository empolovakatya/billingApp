package repository

import (
	billing "billingApp"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

// get balance amount and freezed amount from db and find unfreeze amount to use
func getFreeBalance(balanceId uint64) (uint64, error) {

	db, err := NewPostgresDB()
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	tx, err := db.Begin()
	var balance billing.Balance
	var freezed billing.Freeze
	getBalanceQuery := fmt.Sprintf("SELECT balance_id, amount FROM %s WHERE balance_id = $1",
		balancesTable)
	row := tx.QueryRow(getBalanceQuery, balanceId)
	err = row.Scan(&balance.BalanceId, &balance.Amount)
	if err != nil {
		tx.Rollback()
		return 0, nil
	}
	getFreezedQuery := fmt.Sprintf("SELECT freeze_id, balance_id, freezed_amount FROM %s WHERE balance_id = $1",
		freezeTable)
	row = tx.QueryRow(getFreezedQuery, balanceId)
	err = row.Scan(&freezed.FreezeId, &freezed.BalanceId, &freezed.FreezedAmount)
	logrus.Print("freeze ", freezed)
	if err != nil {
		return balance.Amount, nil
	}
	freeBalance := balance.Amount - freezed.FreezedAmount
	return freeBalance, tx.Commit()
}

// ErrorResponse help to return error message on response
func ErrorResponse(msg string, err error) []byte {
	response := billing.Errors{ErrMessage: fmt.Sprintf(msg, err)}
	data := billing.ErrorResponse{Data: response}
	body, _ := json.Marshal(data)
	return body
}

package repository

import (
	billing "billingApp"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

const presicion = 1e6

//IntToFloat converts int to float representation
func IntToFloat(x float64) float64 {
	return x / presicion
}

//getFreeBalance gets balance amount and freezed amount from db and finds unfreeze amount to use
func getFreeBalance(balanceId uint64) (float64, error) {

	db, err := NewPostgresDB()
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var balance billing.Balance
	var freezed billing.Freeze
	var freezes []billing.Freeze

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
	rows, err := tx.Query(getFreezedQuery, balanceId)
	for rows.Next() {
		err = rows.Scan(&freezed.FreezeId, &freezed.BalanceId, &freezed.FreezedAmount)
		freezes = append(freezes, freezed)
	}
	if err != nil {
		return balance.Amount, nil
	}

	var freezedBalance float64
	for _, value := range freezes {
		freezedBalance += value.FreezedAmount
	}

	freeBalance := balance.Amount - freezedBalance
	return freeBalance, tx.Commit()
}

// ErrorResponse helps to return error message on response
func ErrorResponse(msg string, err error) []byte {
	response := billing.Errors{ErrMessage: fmt.Sprintf(msg, err)}
	data := billing.ErrorResponse{Data: response}
	body, _ := json.Marshal(data)
	return body
}

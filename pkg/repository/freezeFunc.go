package repository

import (
	billing "billingApp"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

//FreezeAmount freezes balance and update db
func FreezeAmount(input billing.Worker) ([]byte, error) {
	db, err := NewPostgresDB()
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		body := ErrorResponse("Failed on begin tx %s", err)
		return body, nil
	}

	var balance billing.Balance
	query := fmt.Sprintf("SELECT balance_id, amount FROM %s WHERE balance_id = $1",
		balancesTable)
	row := tx.QueryRow(query, input.BalanceId)
	err = row.Scan(&balance.BalanceId, &balance.Amount)
	if err != nil {
		tx.Rollback()
		body := ErrorResponse("Error on getting balance from db %s", err)
		return body, nil
	}
	freeBalance, err := getFreeBalance(input.BalanceId)
	if err != nil {
		tx.Rollback()
		body := ErrorResponse("Error on getting freebalance from db %s", err)
		return body, nil
	}
	if freeBalance >= input.FreezedAmount {
		var freezeBalance billing.Freeze
		query := fmt.Sprintf("INSERT INTO %s (balance_id, freezed_amount) values ($1, $2) RETURNING freeze_id", freezeTable)
		row := tx.QueryRow(query, input.BalanceId, input.FreezedAmount)
		err = row.Scan(&freezeBalance.FreezeId)
		if err != nil {
			tx.Rollback()
			body := ErrorResponse("Failed on insert freezed to db %s", err)
			return body, nil
		}
		response := billing.ArgsFreezes{FreezeId: freezeBalance.FreezeId, FreezedAmount: IntToFloat(input.FreezedAmount), Msg: "balance-freezed"}
		data := billing.ResponseFreezes{Data: response}
		body, _ := json.Marshal(data)
		return body, tx.Commit()
	}
	tx.Rollback()
	body := ErrorResponse("Not enough money", err)
	return body, nil
}

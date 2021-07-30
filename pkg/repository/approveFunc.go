package repository

import (
	billing "billingApp"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

//Approve approves freezed balance, send or unfreeze money
func Approve(input billing.Worker) ([]byte, error) {
	var freezedBalance billing.Balance
	var balance billing.Balance
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

	getFreezedAmountQuery := fmt.Sprintf("SELECT balance_id, freezed_amount FROM %s WHERE freeze_id = %d", freezeTable, input.FreezeId)
	row := tx.QueryRow(getFreezedAmountQuery)
	err = row.Scan(&freezedBalance.BalanceId, &freezedBalance.Amount)
	if err != nil {
		tx.Rollback()
		body := ErrorResponse("Error on getting freezed balance from db %s", err)
		return body, nil
	}
	getBalanceQuery := fmt.Sprintf("SELECT balance_id, amount FROM %s WHERE balance_id = $1", balancesTable)
	row = tx.QueryRow(getBalanceQuery, freezedBalance.BalanceId)
	err = row.Scan(&balance.BalanceId, &balance.Amount)
	if err != nil {
		tx.Rollback()
		body := ErrorResponse("Failed on getting free balance %s", err)
		return body, nil
	}
	if input.IsApproved == true {
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE freeze_id = $1", freezeTable)
		_, err = tx.Exec(deleteQuery, input.FreezeId)
		if err != nil {
			tx.Rollback()
			body := ErrorResponse("Failed on deleting from db %s", err)
			return body, nil
		}
		tx.Commit()
		result, err := Decrease(billing.Worker{BalanceId: freezedBalance.BalanceId, Amount: freezedBalance.Amount})
		if err != nil {
			tx.Rollback()
			body := ErrorResponse("Failed on decrease balance %s", err)
			return body, nil
		}
		return result, nil
	} else if input.IsApproved == false {
		deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE freeze_id = $1", freezeTable)
		_, err = tx.Exec(deleteQuery, input.FreezeId)
		if err != nil {
			tx.Rollback()
			body := ErrorResponse("Failed on deleting from db %s", err)
			return body, nil
		}
		response := billing.Args{BalanceId: balance.BalanceId, Amount: IntToFloat(balance.Amount), Msg: "balance-unfreezed"}
		data := billing.Response{Data: response}
		body, _ := json.Marshal(data)
		return body, tx.Commit()
	} else {
		tx.Rollback()
		body := ErrorResponse("Invalid command", err)
		return body, nil
	}
}

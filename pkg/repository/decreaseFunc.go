package repository

import (
	billing "billingApp"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

// Decrease decrease balance amount and update db
func Decrease(input billing.Worker) ([]byte, error) {
	db, err := NewPostgresDB()
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}
	tx, err := db.Begin()

	var balance billing.Balance
	balanceId := input.BalanceId
	sumToDecrease := input.Amount

	query := fmt.Sprintf("SELECT balance_id, amount FROM %s WHERE balance_id = $1",
		balancesTable)
	row := tx.QueryRow(query, balanceId)
	err = row.Scan(&balance.BalanceId, &balance.Amount)
	if err != nil {
		tx.Rollback()
		body := ErrorResponse("Error on getting balance from db %s", err)
		return body, nil
	}
	freeBalance, err := getFreeBalance(balanceId)
	if err != nil {
		tx.Rollback()
		body := ErrorResponse("Error on getting freebalance from db %s", err)
		return body, nil
	}
	if freeBalance >= sumToDecrease {
		balance.Amount -= sumToDecrease
		decreaseQuery := fmt.Sprintf("UPDATE %s SET amount = %d WHERE balance_id = %d",
			balancesTable, balance.Amount, balanceId)

		_, err = tx.Exec(decreaseQuery)
		if err != nil {
			tx.Rollback()
			body := ErrorResponse("Failed on decrease balance %s", err)
			return body, nil
		}
		response := billing.Args{BalanceId: balanceId, Amount: balance.Amount, Msg: "balance-changed"}
		data := billing.Response{Data: response}
		body, _ := json.Marshal(data)
		return body, tx.Commit()

	}
	tx.Rollback()
	body := ErrorResponse("Not enough money", err)
	return body, nil
}

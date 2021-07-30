package repository

import (
	billing "billingApp"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

// Increase increase balance amount and update db
func Increase(input billing.Worker) ([]byte, error) {
	db, err := NewPostgresDB()
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	var balance billing.Balance
	balanceId := input.BalanceId
	sumToIncrease := input.Amount

	updateQuery := fmt.Sprintf("UPDATE %s SET amount = amount+ %e WHERE balance_id = %d",
		balancesTable, sumToIncrease, balanceId)
	_, err = db.Exec(updateQuery)
	if err != nil {
		body := ErrorResponse("Failed on update db %s", err)
		return body, nil
	}

	selectQuery := fmt.Sprintf("SELECT balance_id, amount FROM %s WHERE balance_id = $1",
		balancesTable)
	err = db.Get(&balance, selectQuery, balanceId)
	if err != nil {
		body := ErrorResponse("Error on getting balance from db %s", err)
		return body, nil
	}

	response := billing.Args{BalanceId: balanceId, Amount: IntToFloat(balance.Amount), Msg: "balance-changed"}
	data := billing.Response{Data: response}
	body, _ := json.Marshal(data)
	return body, nil
}

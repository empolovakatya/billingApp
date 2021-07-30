package repository

import (
	billing "billingApp"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

//SendToOther sends money to other person
func SendToOther(input billing.Worker) ([]byte, error) {
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

	var senderBalance, receiverBalance billing.Balance
	senderBalanceId := input.BalanceId
	payment := input.Amount
	receiverBalanceId := input.Receiver

	senderQuery := fmt.Sprintf("SELECT balance_id, amount FROM %s WHERE balance_id = $1",
		balancesTable)
	row := tx.QueryRow(senderQuery, senderBalanceId)
	err = row.Scan(&senderBalance.BalanceId, &senderBalance.Amount)
	if err != nil || senderBalance.BalanceId == 0 {
		tx.Rollback()
		body := ErrorResponse("Failed on get sender", err)
		return body, nil
	}

	receiverQuery := fmt.Sprintf("SELECT balance_id, amount FROM %s WHERE balance_id = $1",
		balancesTable)
	row = tx.QueryRow(receiverQuery, receiverBalanceId)
	err = row.Scan(&receiverBalance.BalanceId, &receiverBalance.Amount)
	if err != nil || receiverBalance.BalanceId == 0 {
		tx.Rollback()
		body := ErrorResponse("Failed on find receiver", err)
		return body, nil
	}

	senderFreeBalance, err := getFreeBalance(senderBalanceId)
	if senderFreeBalance >= payment && receiverBalance.BalanceId > 0 && senderBalance.BalanceId > 0 {

		updateBalanceSenderQuery := fmt.Sprintf("UPDATE %s SET amount = amount - %e WHERE balance_id = %d",
			balancesTable, payment, senderBalanceId)
		updateBalanceReceiverQuery := fmt.Sprintf("UPDATE %s SET amount = amount + %e WHERE balance_id = %d",
			balancesTable, payment, receiverBalanceId)
		_, err = tx.Exec(updateBalanceSenderQuery)
		if err != nil {
			tx.Rollback()
			body := ErrorResponse("Failed on update sender balance %s", err)
			return body, nil
		}
		_, err = tx.Exec(updateBalanceReceiverQuery)
		if err != nil {
			tx.Rollback()
			body := ErrorResponse("Failed on update receiver balance %s", err)
			return body, nil
		}

		response := billing.MoneyTransfered{
			SenderId:        senderBalanceId,
			SenderBalance:   IntToFloat(senderBalance.Amount - payment),
			ReceiverId:      receiverBalanceId,
			ReceiverBalance: IntToFloat(receiverBalance.Amount + payment),
			Msg:             "money-transfered"}
		data := billing.TransferedResponse{Data: response}
		body, _ := json.Marshal(data)
		return body, tx.Commit()
	} else {
		tx.Rollback()
		body := ErrorResponse("Not enough money", err)
		return body, nil
	}
}

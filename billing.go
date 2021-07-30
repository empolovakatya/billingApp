package billing

// Balance uses to do smth with `balances` table
type Balance struct {
	BalanceId uint64  `json:"balance_id" db:"balance_id"`
	Amount    float64 `json:"amount" db:"amount" binding:"required"`
}

//Freeze uses to do smth with `freezes` table
type Freeze struct {
	FreezeId      uint64  `json:"freeze_id" db:"freeze_id"`
	BalanceId     uint64  `json:"balance_id" db:"balance_id" binding:"required"`
	FreezedAmount float64 `json:"freezed_amount" db:"freezed_amount" binding:"required"`
	IsApproved    bool    `json:"is_approved" db:"is_approved"`
}

//Worker uses to get requests from sender
type Worker struct {
	Method        string  `json:"method"`
	BalanceId     uint64  `json:"balance_id"`
	Amount        float64 `json:"amount"`
	Receiver      uint64  `json:"receiver"`
	FreezeId      uint64  `json:"freeze_id"`
	FreezedAmount float64 `json:"freezed_amount"`
	IsApproved    bool    `json:"is_approved"`
}

//Args uses to get response from events with `balances` table
type Args struct {
	BalanceId uint64  `json:"balance_id"`
	Amount    float64 `json:"amount"`
	Msg       string  `json:"msg"`
}

//ArgsFreezes uses to get response from events with `freezes` table
type ArgsFreezes struct {
	FreezeId      uint64  `json:"freeze_id"`
	FreezedAmount float64 `json:"freezed_amount"`
	Msg           string  `json:"msg"`
}

type Response struct {
	Data Args `json:"data"`
}

type ResponseFreezes struct {
	Data ArgsFreezes `json:"data"`
}

//WorkerSender uses to send request from sender
type WorkerSender struct {
	Method        string  `json:"method"`
	BalanceId     float64 `json:"balance_id"`
	Amount        float64 `json:"amount"`
	Receiver      uint64  `json:"receiver"`
	FreezeId      float64 `json:"freeze_id"`
	FreezedAmount float64 `json:"freezed_amount"`
	IsApproved    bool    `json:"is_approved"`
}

//Errors uses to return error messages on response
type Errors struct {
	ErrMessage string `json:"error"`
}

type ErrorResponse struct {
	Data Errors `json:"data"`
}

//MoneyTransfered uses to get response from method `send`
type MoneyTransfered struct {
	SenderId        uint64  `json:"sender_id"`
	SenderBalance   float64 `json:"sender_balance"`
	ReceiverId      uint64  `json:"receiver_id"`
	ReceiverBalance float64 `json:"receiver_balance"`
	Msg             string  `json:"msg"`
}

type TransferedResponse struct {
	Data MoneyTransfered `json:"data"`
}

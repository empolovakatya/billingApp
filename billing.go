package billing

type Balance struct {
	BalanceId uint64 `json:"balance_id" db:"balance_id"`
	Amount    uint64 `json:"amount" db:"amount" binding:"required"`
}

type Freeze struct {
	FreezeId      uint64 `json:"freeze_id" db:"freeze_id"`
	BalanceId     uint64 `json:"balance_id" db:"balance_id" binding:"required"`
	FreezedAmount uint64 `json:"freezed_amount" db:"freezed_amount" binding:"required"`
	IsApproved    bool   `json:"is_approved" db:"is_approved"`
}

type Worker struct {
	Method        string `json:"method"`
	BalanceId     uint64 `json:"balance_id"`
	Amount        uint64 `json:"amount"`
	Receiver      uint64 `json:"receiver"`
	FreezeId      uint64 `json:"freeze_id"`
	FreezedAmount uint64 `json:"freezed_amount"`
	IsApproved    bool   `json:"is_approved"`
}

type Args struct {
	BalanceId uint64 `json:"balance_id"`
	Amount    uint64 `json:"amount"`
	Msg       string `json:"msg"`
}

type ArgsFreezes struct {
	FreezeId      uint64 `json:"freeze_id"`
	FreezedAmount uint64 `json:"freezed_amount"`
	Msg           string `json:"msg"`
}

type Response struct {
	Data Args `json:"data"`
}

type ResponseFreezes struct {
	Data ArgsFreezes `json:"data"`
}

type WorkerSender struct {
	Method        string  `json:"method"`
	BalanceId     float64 `json:"balance_id"`
	Amount        float64 `json:"amount"`
	Receiver      float64 `json:"receiver"`
	FreezeId      float64 `json:"freeze_id"`
	FreezedAmount float64 `json:"freezed_amount"`
	IsApproved    bool    `json:"is_approved"`
}

type Errors struct {
	ErrMessage string `json:"error"`
}

type ErrorResponse struct {
	Data Errors `json:"data"`
}

type MoneyTransfered struct {
	SenderId        uint64 `json:"sender_id"`
	SenderBalance   uint64 `json:"sender_balance"`
	ReceiverId      uint64 `json:"receiver_id"`
	ReceiverBalance uint64 `json:"receiver_balance"`
	Msg             string `json:"msg"`
}

type TransferedResponse struct {
	Data MoneyTransfered `json:"data"`
}

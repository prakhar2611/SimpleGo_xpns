package Models

type PayloadToMongo struct {
	Data []ExpenseBO `json:"data"`
}

type ExpenseBO struct {
	TxnId          int64   `json:"txnId"`
	Date           string  `json:"date"`
	Amount         float32 `json:"amount"`
	UpiOfRecipient string  `json:"upiOfRecipient"`
	Category       string  `json:"category"`
	Mode           string  `json:"mode"`
	UserID         int32   `json:"userId"`
	Day            string  `json:"day"`
	Month          string  `json:"month"`
}

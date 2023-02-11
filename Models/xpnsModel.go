package Models

type PayloadToMongo struct {
	BaseResponse
	Data  []interface{} `json:"data"`
	Month string        `json:"month"`
}

type ExpenseBO struct {
	TxnId           string  `json:"txnId"`
	TransactionType string  `json:"transactionType"`
	Date            string  `json:"date"`
	DAmount         float32 `json:"dAmount"`
	CAmount         float32 `json:"cAmount"`
	ReceiverBody    string  `json:"receiverBody"`
	ReceiverPG      string  `json:"receiverPG"`
	ToAccount       string  `json:"toAccount"`
	Info            string  `json:"Info"`
	BalanceLeft     float32 `json:"balanceLeft"`
	Category        string  `json:"category"`
	Mode            string  `json:"mode"`
	UserID          int32   `json:"userId"`
}

type BaseResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

//how to use the object.function call in go lang
//based on the type if object you called
// func (a author) fullName() string {
//     return fmt.Sprintf("%s %s", a.firstName, a.lastName)
// }

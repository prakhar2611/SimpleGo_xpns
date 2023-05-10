package Models

import (
	"time"

	"gorm.io/gorm"
)

type PayloadToMongo struct {
	BaseResponse
	Data  []interface{} `json:"data"`
	Month string        `json:"month"`
}

type ExpenseBO struct {
	gorm.Model
	TxnId           string `gorm:"unique"`
	UserID          string
	Month           string
	Date            time.Time
	TransactionType string
	DAmount         float32
	CAmount         float32
	ReceiverBody    string
	ReceiverPG      string
	ToAccount       string
	Info            string
	BalanceLeft     float32
	Category        string
	Mode            string
}

type B64decodedResponse struct {
	gorm.Model
	AmountDebited string    `json:"amount_debited"`
	Date          string    `json:"date"`
	ETime         time.Time `json:"etime"`
	TransactionId string    `json:"msgId" gorm:"unique"`
	ToAccount     string    `json:"to_vpa"`
	Category      string    `json:"category"`
	Label         string    `json:"label"`
}

type GetEncodedDataReq struct {
	MsgEncodedData string `json:"msgEncodedData"`
	MsgId          string `json:"msgId"`
}

type BaseResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

type GmailExcelThreadSnapShot struct {
	gorm.Model
	ThreadId      int64
	Label         string
	LastHistoryId string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

//how to use the object.function call in go lang
//based on the type if object you called
// func (a author) fullName() string {
//     return fmt.Sprintf("%s %s", a.firstName, a.lastName)
// }

type XmlSheetExpense struct {
	Data []ExpenseBO `json:"data"`
}

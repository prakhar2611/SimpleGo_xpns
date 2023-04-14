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
	TxnId           string `gorm:"primarykey"`
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

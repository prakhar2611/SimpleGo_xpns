package Models

import (
	"time"

	"gorm.io/gorm"
)

//db Models

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

type PocketsMappingDbo struct {
	gorm.Model
	UserId string   `json:"userId"`
	Pocket string   `json:"pocket" gorm:"unique"` //making it unique
	Labels []string `json:"labels"`
}

type B64decodedResponse struct {
	gorm.Model
	AmountDebited float64   `json:"amount_debited"`
	Date          string    `json:"date"`
	ETime         time.Time `json:"etime"`
	TransactionId string    `json:"msgId" gorm:"unique"`
	ToAccount     string    `json:"to_vpa"`
	Pocket        string    `json:"pocket"`
	Label         string    `json:"label"`
	UserId        string    `json:"UserId"`
}

//REST models

type VpaMapping struct {
	Id          int    `json:"id"`
	Vpa         string `json:"vpa"`
	TotalAmount int64  `json:"totalAmount"`
	TotalTxn    int64  `json:"totalTxn"`
	Category    string `json:"category"`
	Label       string `json:"label"`
	UserId      string `json:"UserId"`
}

type GetEncodedDataReq struct {
	MsgEncodedData string `json:"msgEncodedData"`
	MsgId          string `json:"msgId"`
}

type BaseResponse struct {
	IsNewUser bool   `json:"isNewUser"`
	Status    bool   `json:"status"`
	Error     string `json:"error"`
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

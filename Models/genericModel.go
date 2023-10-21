package Models

import (
	"time"
)

type SignInCallbackResponse struct {
	BaseResponse
	Name       string `json:"name"`
	Token      string `json:"token"`
	RedirectTo string `json:"redirectTo"`
}

type SignInResponse struct {
	BaseResponse
	URl string `json:"url"`
}

type UserProfile struct {
	BaseResponse
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Userid  string `json:"userId"`
}

type UserXpnsData struct {
	AmountDebited string    `json:"amount_debited"`
	Date          string    `json:"date"`
	ETime         time.Time `json:"etime"`
	TransactionId string    `json:"msgId" gorm:"unique"`
	ToAccount     string    `json:"to_vpa"`
}

type UpdatecategoryPayload struct {
	Data map[string]categoryL `json:"data"`
}

type UpdateCategoryResponse struct {
	FailureMsgId []string `json:"failuerMsgId"`
}

type SyncUpResp struct {
	BaseResponse
	FailedTxns []string `json:"failedTxns"`
}

type categoryL struct {
	Category string `json:"category"`
	Label    string `json:"label"`
}

type UpdatePocketsPayload struct {
	Data []PocketsMappingDbo `json:"data"`
}

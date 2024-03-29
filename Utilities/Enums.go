package Utilities

const (
	TxnDate     = 0
	Description = 2
	RefNo       = 3
	DebitAmt    = 4
	CreditAmt   = 5
	Balance     = 6
)

const (
	Credit = "Credit"
	Debit  = "Debit"
)

func ElementExists(array []string, target string) bool {
	for _, item := range array {
		if item == target {
			return true
		}
	}
	return false
}

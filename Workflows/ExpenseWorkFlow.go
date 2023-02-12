package Workflows

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	Model "SimpleGo_xpns/Models"

	"SimpleGo_xpns.go/Utilities"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/spf13/viper"
)

func SBIparser(FileName string) Model.PayloadToMongo {
	var req Model.PayloadToMongo
	//k := Utilities.GetKaonf()

	f, err := excelize.OpenFile(FileName)
	if err != nil {
		log.Fatalln(err)
	}
	firstSheet := f.WorkBook.Sheets.Sheet[0].Name
	fmt.Printf("'%s' is first sheet of %d sheets.\n", firstSheet, f.SheetCount)
	rows := f.GetRows(firstSheet)
	if err != nil {
		log.Fatalln(err)
	}

	sheetPreOffset := viper.Get("SBISheetOffset")
	var i int
	if sheetPreOffset == "" {
		i = 20
	} else {
		d, e := sheetPreOffset.(float64)
		if e {
			i = int(d)
		}
	}

	r := rows[i : len(rows)-2]
	for _, row := range r {
		//for _, colCell := range row {
		// 	fmt.Print(colCell, "\t")
		// }
		//fmt.Println(row)

		//prepare a single transaction from the col
		var expense Model.ExpenseBO

		expense.Date = row[Utilities.TxnDate]

		if row[Utilities.DebitAmt] != "" || len(row[Utilities.DebitAmt]) == 0 {
			j, _ := strconv.ParseFloat(row[Utilities.DebitAmt], 32)
			expense.DAmount = float32(j)
			expense.TransactionType = Utilities.Debit
		}
		if row[Utilities.CreditAmt] != "" || len(row[Utilities.CreditAmt]) == 0 {
			j, _ := strconv.ParseFloat(row[Utilities.CreditAmt], 32)
			expense.CAmount = float32(j)
			expense.TransactionType = Utilities.Credit
		}

		//getting majordata from the description
		description := strings.Split(row[Utilities.Description], "/")
		if len(description) > 1 {
			expense.TxnId = strings.ToUpper(strings.TrimSpace(description[2]))
			expense.ReceiverBody = strings.ToUpper(strings.TrimSpace(description[3]))
			expense.ReceiverPG = strings.ToUpper(strings.TrimSpace(description[4]))
			expense.ToAccount = strings.ToUpper(strings.TrimSpace(description[5]))
			expense.Info = description[6]
		} else {
			expense.Info = row[Utilities.Description]
		}

		bal, _ := strconv.ParseFloat(row[Utilities.Balance], 32)
		expense.BalanceLeft = float32(bal)

		//adding transction to the list of daily transcation based on date of transaction
		var e interface{} = expense
		req.Data = append(req.Data, e)
		req.Month = strings.Split(expense.Date, "-")[1]
	}
	return req
}

func SendToMongo(payload Model.PayloadToMongo) bool {

	c, _ := Utilities.ConectToMongo()
	col := c.Database("Expense_new").Collection(payload.Month)

	if _, err := col.InsertMany(context.Background(), payload.Data); err == nil {
		return true
	} else {
		return false
	}
}

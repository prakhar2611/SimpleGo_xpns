package Workflows

import (
	Model "SimpleGo_xpns/Models"
	"SimpleGo_xpns/Utilities"
	Executer "SimpleGo_xpns/Utilities/APIExecuter"
	dbConnector "SimpleGo_xpns/Utilities/DbConnector"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
)

func SBIparser(FileName string, userid string, str io.Reader) []Model.ExpenseBO {
	var req []Model.ExpenseBO
	//k := Utilities.GetKaonf()

	files := strings.Split(FileName, ",")

	for _, file := range files {
		var err error
		var f *excelize.File
		filepath := "SBIFiles/" + file + ".xlsx"
		if file != "" {
			f, err = excelize.OpenFile(filepath)
		} else if str != nil {
			f, err = excelize.OpenReader(str)
		}

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

			//prepare a single transaction from thSe col
			var expense Model.ExpenseBO
			d := row[Utilities.TxnDate]
			myDate, err := time.Parse("2-Jan-06", d)
			if err != nil {
				panic(err)
			}

			expense.Date = myDate

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
			if len(description) == 7 {
				expense.TxnId = strings.ToUpper(strings.TrimSpace(description[2]))
				expense.ReceiverBody = strings.ToUpper(strings.TrimSpace(description[3]))
				expense.ReceiverPG = strings.ToUpper(strings.TrimSpace(description[4]))
				expense.ToAccount = strings.ToUpper(strings.TrimSpace(description[5]))
				expense.Info = description[6]
			} else {
				expense.TxnId = row[Utilities.Description]
			}

			bal, _ := strconv.ParseFloat(row[Utilities.Balance], 32)
			expense.BalanceLeft = float32(bal)
			expense.Month = strings.Split(d, "-")[1]
			expense.UserID = userid
			fmt.Println(expense)
			req = append(req, expense)
		}

	}
	return req
}

func GetDataForbase64(payload []Model.GetEncodedDataReq) []*Model.B64decodedResponse {
	var req Executer.APIRequest

	req.BaseURL = viper.Get("internalService").(string)
	req.Action = fmt.Sprintf("/api/v1/decode/")
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	req.Headers = headers
	resp := Executer.POST(req, payload)
	if resp.Status == 200 {
		var data []*Model.B64decodedResponse
		err := json.Unmarshal([]byte(resp.Response), &data)
		if err == nil {
			return data
		}
	}
	return nil
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

//sending data to sql db
func SendPayloadTodb(token string, payload []Model.ExpenseBO) bool {
	//addding user information here as extending the model
	//verifyIdToken(accessToken)
	//.
	if dbConnector.SendDataToPostgres(payload) {
		return true
	} else {
		return false
	}
}

//sending data to Excel sheets
func SendPayloadToExcel(token string, payload []Model.ExpenseBO) bool {
	var sheet Model.XmlSheetExpense
	sheet.Data = payload
	//calling google sheet api to create speardsheet and write excel on it

	return true
}

func getContentsFromBody(ebd string) (string, string, string) {
	amountDebitedRegex := regexp.MustCompile(`Rs.\d*`)
	toVPARegex := regexp.MustCompile(`VPA.+?on`)
	dateRegex := regexp.MustCompile(`\d{2}-\d{2}-\d{2}`)

	amountDebited := amountDebitedRegex.FindString(ebd)
	toVPA := toVPARegex.FindString(ebd)
	date := dateRegex.FindString(ebd)

	if date == "" {
		date = "None"
	}
	if toVPA != "" {
		toVPA = strings.TrimPrefix(toVPA, "VPA ")
		toVPA = strings.TrimSuffix(toVPA, " on")
	} else {
		toVPA = "None"
	}
	if amountDebited != "" {
		amountDebited = strings.TrimPrefix(amountDebited, "Rs.")
	} else {
		amountDebited = "None"
	}

	return date, toVPA, amountDebited
}

func getEmailBody(codedBody string) (string, string, string) {
	data := strings.ReplaceAll(codedBody, "-", "+")
	data = strings.ReplaceAll(data, "_", "/")

	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatal("Error decoding base64 data:", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(decodedData)))
	if err != nil {
		log.Fatal("Error parsing HTML:", err)
	}

	emailBody := doc.Find("td.td").First().Text()
	date, toVPA, amountDebited := getContentsFromBody(emailBody)
	return date, toVPA, amountDebited
}

func ExtractBodyFromEncodedData(codedRelevantJSON []Model.GetEncodedDataReq) []Model.B64decodedResponse {
	var decodedExtractedInfo []Model.B64decodedResponse

	for _, message := range codedRelevantJSON {
		data := message.MsgEncodedData

		date, toVPA, amountDebited := getEmailBody(data)
		var amount float64
		if amountDebited != "" {
			amount, _ = strconv.ParseFloat(amountDebited, 64)
		}
		d := Model.B64decodedResponse{
			TransactionId: message.MsgId,
			Date:          date,
			ToAccount:     toVPA,
			AmountDebited: amount,
		}

		decodedExtractedInfo = append(decodedExtractedInfo, d)
	}

	return decodedExtractedInfo
}

func UpdateVpaMapping(request *Model.UpdatecategoryPayload, userId string) bool {

	//two things to do
	// 1. update the vpa mapped in the Db
	// 2. update the vpa in the vpa mapping db - maitaning for fun

	//1
	//dbConnector.PushVPAMAppingToDb(request)
	success, _ := dbConnector.UpdateVPATxnLevel(request, userId)
	if success {
		return true

	}
	return false
}
